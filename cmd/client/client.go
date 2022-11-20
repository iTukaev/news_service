package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/iTukaev/news_service/internal/client"
	yamlPkg "github.com/iTukaev/news_service/internal/config/yaml"
	postgresPkg "github.com/iTukaev/news_service/internal/repo/postgres"
	customRouter "github.com/iTukaev/news_service/internal/router"
	loggerPkg "github.com/iTukaev/news_service/pkg/logger"
)

func main() {
	var cfg config
	var err error
	cfg, err = yamlPkg.New()
	if err != nil {
		log.Fatalln("Config init error:", err)
	}

	logger, err := loggerPkg.New(loggerPkg.LogLevel(cfg.LogLevel()))
	if err != nil {
		log.Fatalln("Logger init error:", err)
	}
	logger.Infoln("Start main")

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		logger.Infoln("Shutting down...")
		cancel()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		if err = start(ctx, cfg, logger); err != nil {
			logger.Errorln("client start", err)
		}
		c <- os.Interrupt
	}()

	<-c
}

func start(ctx context.Context, cfg config, logger *zap.SugaredLogger) (retErr error) {
	pool, err := postgresPkg.NewPostgres(ctx, cfg.PGConfig(), logger)
	if err != nil {
		return errors.Wrap(err, "postgres")
	}
	repo := postgresPkg.New(pool, logger)

	app := client.NewApp(repo, logger)
	router := customRouter.New(app, logger)
	router.RegisterServices()

	srv := http.Server{
		Addr:    cfg.ClientHTTP(),
		Handler: router.Mux(),
	}

	stopCh := make(chan struct{}, 0)
	go func() {
		if err = srv.ListenAndServe(); err != nil {
			retErr = errors.Wrap(err, "listen")
		}
		close(stopCh)
	}()

	select {
	case <-ctx.Done():
	case <-stopCh:
	}
	return retErr
}
