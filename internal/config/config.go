package config

import (
	pgModel "github.com/iTukaev/news_service/internal/repo/models"
)

type Interface interface {
	LogLevel() string
	ServiceURL() string
	ClientHTTP() string
	PGConfig() pgModel.Config
}
