package main

import (
	pgModel "github.com/iTukaev/news_service/internal/repo/models"
)

type config interface {
	LogLevel() string
	ClientHTTP() string
	PGConfig() pgModel.Config
}
