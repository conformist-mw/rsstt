package service

import (
	"sort"

	"github.com/conformist-mw/rsstt/models"
	"github.com/conformist-mw/rsstt/repository"
	"github.com/mmcdole/gofeed"
)

func CreateFeed(url string) {
	existingFeed := repository.GetFeedByUrl(url)
	if existingFeed.ID != 0 {
		return
	}
	fp := gofeed.NewParser()
	f, _ := fp.ParseURL(url)
	repository.CreateFeed(url, f)
}

func FetchFeed(feed *models.Feed) {
	existingItemLinks := repository.FetchItemLinksForFeed(feed.ID)
	existingItemsMap := make(map[string]struct{})
	for _, link := range existingItemLinks {
		existingItemsMap[link] = struct{}{}
	}
	fp := gofeed.NewParser()
	f, _ := fp.ParseURL(feed.Url)

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
	repository.SetLastFetchedAt(feed)
}

func UpdateFeeds() {
	feeds := repository.FetchAllFeeds()
	for _, feed := range feeds {
		FetchFeed(&feed)
	}
}
