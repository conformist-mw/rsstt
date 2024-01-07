package main

import (
	"github.com/conformist-mw/rsstt/models"
	"github.com/conformist-mw/rsstt/repository"
	"github.com/conformist-mw/rsstt/service"
)

func main() {
	models.ConnectDb()
	service.CreateFeed("https://www.opennet.me/opennews/opennews_all_utf.rss")
	feeds := repository.FetchOutdatedFeeds()
	for _, feed := range feeds {
		service.FetchFeed(&feed)
	}
}
