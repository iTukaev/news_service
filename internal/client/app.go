package client

import (
	"context"

	"go.uber.org/zap"

	repoModel "github.com/iTukaev/news_service/internal/repo/models"
)

const (
	maxLimit = 10
)

type App struct {
	repo   repo
	logger *zap.SugaredLogger
}

func NewApp(repo repo, logger *zap.SugaredLogger) *App {
	return &App{
		repo:   repo,
		logger: logger,
	}
}

type repo interface {
	NewsGet(ctx context.Context, search string) (*repoModel.News, error)
	NewsList(ctx context.Context, params repoModel.ListParams) ([]repoModel.News, error)
	Close()
}

func (a *App) NewsGet(ctx context.Context, search string) (*News, error) {
	news, err := a.repo.NewsGet(ctx, search)
	if err != nil {
		return nil, err
	}

	return repoToAppNews(news), nil
}

func (a *App) NewsList(ctx context.Context, limit, offset uint64, order bool) ([]News, error) {
	if limit < 1 || limit > maxLimit {
		limit = maxLimit
	}

	news, err := a.repo.NewsList(ctx, repoModel.ListParams{
		Limit:  limit,
		Offset: offset,
		Order:  order,
	})
	if err != nil {
		return nil, err
	}

	return repoToAppNewsList(news), nil
}
