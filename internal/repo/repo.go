//go:generate mockgen -source=repo.go -destination=./mock/repo_mock.go -package=mock

package repo

import (
	"context"

	"github.com/iTukaev/news_service/internal/repo/models"
)

type Interface interface {
	NewsExists(ctx context.Context, id uint32) error
	NewsInsert(ctx context.Context, news *models.News) error
	NewsGet(ctx context.Context, search string) (*models.News, error)
	NewsList(ctx context.Context, params models.ListParams) ([]models.News, error)
	Close()
}
