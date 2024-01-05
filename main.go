package main

import (
	"fmt"

	"github.com/conformist-mw/rsstt/models"
	"github.com/mmcdole/gofeed"
)

func main() {
	models.ConnectDb()
	var feeds []models.Feed
	models.DB.Find(&feeds)
	fp := gofeed.NewParser()
	for _, feed := range feeds {
		f, _ := fp.ParseURL("https://www.opennet.me/opennews/opennews_all_utf.rss")
		for _, item := range f.Items {
			fmt.Println(item.Title)
			fmt.Println(item.Link)
		}
	}
}
