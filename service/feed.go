package service

import (
	"context"
	"net/http"
	"sort"
	"sync"
	"time"

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

	timeout := 30 * time.Second
	if serviceConfig != nil {
		timeout = time.Duration(serviceConfig.HTTPTimeout) * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	fp := gofeed.NewParser()
	fp.Client = &http.Client{Timeout: timeout}
	f, error := fp.ParseURLWithContext(url, ctx)
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

	timeout := 30 * time.Second
	if serviceConfig != nil {
		timeout = time.Duration(serviceConfig.HTTPTimeout) * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	fp := gofeed.NewParser()
	fp.Client = &http.Client{Timeout: timeout}
	f, err := fp.ParseURLWithContext(feed.Url, ctx)
	if err != nil {
		logger.Log.Error("Error parsing feed", feed.Url, err)
		return
	}

	// Sort items by published date (oldest first)
	sort.Slice(f.Items, func(i, j int) bool {
		if f.Items[i].PublishedParsed != nil && f.Items[j].PublishedParsed != nil {
			return f.Items[i].PublishedParsed.Before(*f.Items[j].PublishedParsed)
		}
		// If one has published date and other doesn't, prioritize the one with date
		if f.Items[i].PublishedParsed != nil {
			return true
		}
		if f.Items[j].PublishedParsed != nil {
			return false
		}
		// If neither has published date, keep original order
		return false
	})

	newItemsCount := 0
	for _, item := range f.Items {
		if _, exists := existingItemsMap[item.Link]; !exists {
			repository.AddItem(feed.ID, item)
			newItemsCount++
		}
	}

	repository.SetLastFetchedAt(feed)
	logger.Log.Debugf("Feed fetched %s, added %d new items", feed.Url, newItemsCount)
}

func UpdateFeeds() {
	logger.Log.Debug("Starting feeds update")
	feeds := repository.FetchActiveFeeds()

	maxConcurrent := 5
	if serviceConfig != nil {
		maxConcurrent = serviceConfig.MaxConcurrentFeeds
	}

	semaphore := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup

	for _, feed := range feeds {
		wg.Add(1)
		go func(f models.Feed) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			defer func() {
				if r := recover(); r != nil {
					logger.Log.Errorf("Panic while fetching feed %s: %v", f.Url, r)
				}
			}()

			logger.Log.Debug("Updating feed", f.Url)
			FetchFeed(&f)
		}(feed)
	}

	wg.Wait()
	logger.Log.Debug("Feeds update completed")
}
