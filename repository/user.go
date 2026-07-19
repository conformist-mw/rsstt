package repository

import (
	"time"

	"rsstt/models"

	"gorm.io/gorm/clause"
)

func GetUnseenItems(limit int) []models.Item {
	var items []models.Item
	last24Hours := time.Now().Add(-24 * time.Hour)

	models.DB.Where("seen_items.item_id IS NULL AND items.created_at >= ?", last24Hours).
		Joins("JOIN feeds ON items.feed_id = feeds.id AND feeds.is_active = TRUE AND feeds.deleted_at IS NULL").
		Joins("LEFT OUTER JOIN seen_items ON items.id = seen_items.item_id").
		Order("items.created_at DESC").
		Limit(limit).
		Find(&items)

	// add up to 2 old unseen items every 3 hours between 8am and 9pm
	now := time.Now()
	hour := now.Hour()
	if hour >= 8 && hour <= 21 && hour%3 == 0 && now.Minute() < 5 {
		var oldItems []models.Item
		models.DB.Where("seen_items.item_id IS NULL AND items.created_at < ?", last24Hours).
			Joins("JOIN feeds ON items.feed_id = feeds.id AND feeds.is_active = TRUE AND feeds.deleted_at IS NULL").
			Joins("LEFT OUTER JOIN seen_items ON items.id = seen_items.item_id").
			Order("items.created_at DESC").
			Limit(2).
			Find(&oldItems)
		items = append(items, oldItems...)
	}

	return items
}

func MarkItemsAsSeen(itemIDs []uint) {
	if len(itemIDs) == 0 {
		return
	}

	seenItems := make([]models.SeenItem, 0, len(itemIDs))
	for _, itemID := range itemIDs {
		seenItems = append(seenItems, models.SeenItem{ItemID: itemID})
	}
	models.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&seenItems)
}
