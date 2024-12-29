package service

import (
	"sort"

	"rsstt/logger"
	"rsstt/models"
	"rsstt/repository"

	"github.com/mmcdole/gofeed"
)

func CreateFeed(url string) (models.Feed, error) {
	logger.Log.Debug("Creating feed", url)
	existingFeed := repository.GetFeedByUrl(url)
	if existingFeed.ID != 0 {
		return existingFeed, nil
	}
	fp := gofeed.NewParser()
	f, error := fp.ParseURL(url)
	if error != nil {
		return models.Feed{}, error
	}
	newFeed := repository.CreateFeed(url, f)
	return newFeed, nil
}

func FetchFeed(feed *models.Feed) {
	logger.Log.Debug("Fetching feed", feed.Url)
	existingItemLinks := repository.FetchItemLinksForFeed(feed.ID)
	existingItemsMap := make(map[string]struct{})
	for _, link := range existingItemLinks {
		existingItemsMap[link] = struct{}{}
	}
	fp := gofeed.NewParser()
	f, err := fp.ParseURL(feed.Url)
	if err != nil {
		logger.Log.Error("Error parsing feed", feed.Url, err)
		return
	}

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
	logger.Log.Debug("Feed fetched", feed.Url)
}

func UpdateFeeds() error {
	logger.Log.Debug("Starting feeds update")
	feeds := repository.FetchAllFeeds()
	for _, feed := range feeds {
		logger.Log.Debug("Updating feed", feed.Url)
		FetchFeed(&feed)
	}
	logger.Log.Debug("Feeds update completed")
	return nil
}
