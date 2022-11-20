package servise

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	errorsPkg "github.com/iTukaev/news_service/internal/customerrors"
	repoModel "github.com/iTukaev/news_service/internal/repo/models"
	"github.com/iTukaev/news_service/pkg/helper"
	rssModel "github.com/iTukaev/news_service/pkg/rss/models"
)

type App struct {
	repo   repo
	news   newsGetter
	logger *zap.SugaredLogger
}

func NewApp(repo repo, logger *zap.SugaredLogger, news newsGetter) *App {
	return &App{
		repo:   repo,
		news:   news,
		logger: logger,
	}
}

type newsGetter interface {
	News(ctx context.Context, news chan rssModel.News)
}

type repo interface {
	NewsExists(ctx context.Context, id uint32) error
	NewsInsert(ctx context.Context, news *repoModel.News) error
	Close()
}

func (a *App) HandleNews(ctx context.Context) error {
	news := make(chan rssModel.News)

	go a.news.News(ctx, news)

	for n := range news {
		select {
		case <-ctx.Done():
			close(news)
			return nil
		default:
			a.handler(ctx, &n, news)
		}
	}
	return nil
}

func (a *App) handler(ctx context.Context, n *rssModel.News, news chan rssModel.News) {
	hash, err := helper.StringToHash(fmt.Sprintf("%s%s", n.Title, n.PubDate.Format(time.RFC3339)))
	if err != nil {
		a.logger.Errorln("News ID generate", err)
		news <- *n
		return
	}

	if err = a.repo.NewsExists(ctx, hash); err != nil {
		if errors.Is(err, errorsPkg.ErrNewsAlreadyExists) {
			return
		}
	}

	if err = a.repo.NewsInsert(ctx, rssToRepoNews(n, hash)); err != nil {
		if !errors.Is(err, errorsPkg.ErrNewsAlreadyExists) {
			timeout := time.Duration(0)
			if errors.Is(err, errorsPkg.ErrTimeout) {
				timeout = 5 * time.Second
			}
			time.Sleep(timeout)
			news <- *n
		}
	}
}
