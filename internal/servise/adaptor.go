package servise

import (
	repoModel "github.com/iTukaev/news_service/internal/repo/models"
	rssModel "github.com/iTukaev/news_service/pkg/rss/models"
)

func rssToRepoNews(n *rssModel.News, id uint32) *repoModel.News {
	return &repoModel.News{
		ID:          id,
		Title:       n.Title,
		Link:        n.Link,
		Description: n.Description,
		PubDate:     n.PubDate,
		Article:     n.Article,
	}
}
