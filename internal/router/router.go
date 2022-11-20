package router

import (
	"context"
	"net/http"

	gorilla "github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/iTukaev/news_service/internal/client"
)

type Core struct {
	router *gorilla.Router
	app    app
	logger *zap.SugaredLogger
}

func New(app app, logger *zap.SugaredLogger) *Core {
	return &Core{
		router: gorilla.NewRouter(),
		app:    app,
		logger: logger,
	}
}

func (c *Core) Mux() *gorilla.Router {
	return c.router
}

func (c *Core) RegisterServices() {
	c.router.Path("/news").
		Queries("search", "").
		HandlerFunc(c.getNews()).
		Methods(http.MethodGet)
	c.router.Path("/news-list").
		Queries("size", "", "page", "", "order", "").
		HandlerFunc(c.getNewsList()).
		Methods(http.MethodGet)
}

type app interface {
	NewsGet(ctx context.Context, search string) (*client.News, error)
	NewsList(ctx context.Context, limit, offset uint64, order bool) ([]client.News, error)
}
