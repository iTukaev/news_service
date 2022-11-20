package client

import (
	repoModel "github.com/iTukaev/news_service/internal/repo/models"
)

func repoToAppNews(n *repoModel.News) *News {
	return &News{
		Title:       n.Title,
		Link:        n.Link,
		Description: n.Description,
		PubDate:     n.PubDate,
		Article:     n.Article,
	}
}

func repoToAppNewsList(news []repoModel.News) []News {
	appNews := make([]News, 0, len(news))
	for _, n := range news {
		appNews = append(appNews, News{
			Title:       n.Title,
			Link:        n.Link,
			Description: n.Description,
			PubDate:     n.PubDate,
			Article:     n.Article,
		})
	}
	return appNews
}
