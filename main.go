package main

import (
	"github.com/conformist-mw/rsstt/models"
	"github.com/conformist-mw/rsstt/repository"
	"github.com/conformist-mw/rsstt/service"
)

func main() {
	models.ConnectDb()
	feeds := repository.FetchOutdatedFeeds()
	for _, feed := range feeds {
		service.FetchFeed(&feed)
	}
}
