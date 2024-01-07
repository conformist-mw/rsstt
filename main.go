package main

import (
	"sort"

	"github.com/conformist-mw/rsstt/models"
	"github.com/conformist-mw/rsstt/repository"
	"github.com/mmcdole/gofeed"
)

func main() {
	models.ConnectDb()
	feeds := repository.FetchOutdatedFeeds()
	fp := gofeed.NewParser()
	for _, feed := range feeds {
		existingItemLinks := repository.FetchItemLinksForFeed(feed.ID)
		existingItemsMap := make(map[string]struct{})
		for _, link := range existingItemLinks {
			existingItemsMap[link] = struct{}{}
		}
		f, _ := fp.ParseURL(feed.Link)

		sort.Slice(f.Items, func(i, j int) bool {
			if f.Items[i].PublishedParsed != nil && f.Items[j].PublishedParsed != nil {
				return f.Items[i].PublishedParsed.Before(*f.Items[j].PublishedParsed)
			}
			return false
		})
		for _, item := range f.Items {
			if _, exists := existingItemsMap[item.Link]; !exists {
				repository.AddItem(feed.ID, item)
			}
		}
		repository.SetLastFetchedAt(&feed)
	}
}
