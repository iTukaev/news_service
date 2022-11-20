package main

import (
	pgModel "github.com/iTukaev/news_service/internal/repo/models"
)

type config interface {
	LogLevel() string
	ServiceURL() string
	PGConfig() pgModel.Config
}
