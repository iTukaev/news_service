package rss

import (
	"context"

	"github.com/iTukaev/news_service/pkg/rss/models"
)

type NewsGetter interface {
	News(ctx context.Context) (chan models.News, error)
}
