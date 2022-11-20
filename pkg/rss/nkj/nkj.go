package nkj

import (
	"context"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	rssModel "github.com/iTukaev/news_service/pkg/rss/models"
)

const (
	baseTicker = time.Hour
	fullNews   = 30 * 24 * time.Hour
)

type Client struct {
	parser  *gofeed.Parser
	logger  *zap.SugaredLogger
	baseURL string
}

func NewClient(logger *zap.SugaredLogger, url string) *Client {
	return &Client{
		parser:  gofeed.NewParser(),
		logger:  logger,
		baseURL: url,
	}
}

func (c *Client) News(ctx context.Context, news chan rssModel.News) {
	ticker := time.NewTicker(baseTicker)

	c.setNewsToChannel(ctx, fullNews, news)
loop:
	for {
		select {
		case <-ticker.C:
			c.setNewsToChannel(ctx, baseTicker, news)
		case <-ctx.Done():
			break loop
		}
	}
}

func (c *Client) setNewsToChannel(ctx context.Context, lastGet time.Duration, news chan rssModel.News) {
	rssNews, err := c.getNews(ctx, lastGet)
	if err != nil {
		// todo: alert may be implemented
		c.logger.Errorln("News getting failed from", c.baseURL, err)
		return
	}
	for _, n := range rssNews {
		news <- n
	}
}

func (c *Client) getNews(ctx context.Context, lastGet time.Duration) ([]rssModel.News, error) {
	now := time.Now()
	feed, err := c.parser.ParseURLWithContext(c.baseURL, ctx)
	if err != nil {
		return nil, errors.Wrap(err, "parse URL")
	}
	news := make([]rssModel.News, 0)
	for _, item := range feed.Items {
		if item.PublishedParsed != nil {
			if item.PublishedParsed.After(now.Add(-lastGet)) {
				news = append(news, rssModel.News{
					Title:       item.Title,
					Link:        item.Link,
					Description: item.Description,
					PubDate:     *item.PublishedParsed,
					Article:     item.Description,
				})
			}
		} else {
			news = append(news, rssModel.News{
				Title:       item.Title,
				Link:        item.Link,
				Description: item.Description,
				PubDate:     now,
				Article:     item.Description,
			})
		}
	}

	return news, nil
}
