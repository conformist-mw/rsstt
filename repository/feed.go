package repository

import (
	"time"

	"github.com/conformist-mw/rsstt/models"
	"github.com/mmcdole/gofeed"
)

func CreateFeed(url string, feed *gofeed.Feed) models.Feed {
	feedModel := models.Feed{
		Url:         url,
		Title:       &feed.Title,
		Description: &feed.Description,
		Link:        feed.Link,
		FeedLink:    &feed.FeedLink,
		FeedType:    feed.FeedType,
		FeedVersion: feed.FeedVersion,
	}
	models.DB.Create(&feedModel)
	return feedModel
}

func GetFeedByUrl(url string) models.Feed {
	var feed models.Feed
	models.DB.Where("url = ?", url).First(&feed)
	return feed
}

func FetchAllFeeds() []models.Feed {
	var feeds []models.Feed
	models.DB.Find(&feeds)
	return feeds
}

func FetchOutdatedFeeds() []models.Feed {
	var feeds []models.Feed
	fifteenMinutesAgo := time.Now().Add(-15 * time.Minute)
	models.DB.Where("last_fetched_at < ? OR last_fetched_at IS NULL", fifteenMinutesAgo).Find(&feeds)
	return feeds
}

func SetLastFetchedAt(feed *models.Feed) {
	models.DB.Model(feed).Update("last_fetched_at", models.DB.NowFunc())
}

func UpdateFeed(feed *models.Feed) {
	models.DB.Save(feed)
}

func FetchItemLinksForFeed(feedID uint) []string {
	var items []models.Item
	var itemLinks []string
	models.DB.Where("feed_id = ?", feedID).Find(&items)
	for _, item := range items {
		itemLinks = append(itemLinks, item.Link)
	}
	return itemLinks
}

func AddItem(feedId uint, feedItem *gofeed.Item) {
	item := models.Item{
		FeedID:      feedId,
		Title:       &feedItem.Title,
		Link:        feedItem.Link,
		Description: &feedItem.Description,
		Content:     &feedItem.Content,
		PublishedAt: feedItem.PublishedParsed,
	}
	models.DB.Create(&item)
}
