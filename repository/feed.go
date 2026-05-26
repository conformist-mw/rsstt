package repository

import (
	"rsstt/models"

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
		IsActive:    true,
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

func FetchActiveFeeds() []models.Feed {
	var feeds []models.Feed
	models.DB.Where("is_active = ?", true).Find(&feeds)
	return feeds
}

func ToggleFeedActive(id uint) {
	var feed models.Feed
	if err := models.DB.First(&feed, id).Error; err != nil {
		return
	}
	models.DB.Model(&feed).Update("is_active", !feed.IsActive)
}

func DeleteFeed(id uint) {
	models.DB.Delete(&models.Feed{}, id)
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
