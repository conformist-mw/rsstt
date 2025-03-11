package repository

import (
	"rsstt/models"
	"time"
)

func CreateUser(tg_chat_id string) {
	user := models.User{
		TgChatId: &tg_chat_id,
	}
	models.DB.Create(&user)
}

func GetUserByTgChatId(tg_chat_id string) models.User {
	var user models.User
	models.DB.Where("tg_chat_id = ?", tg_chat_id).First(&user)
	return user
}

func CreateSubscription(user_id uint, feed_id uint) {
	subscription := models.Subscription{
		UserID:   user_id,
		FeedID:   feed_id,
		IsActive: true,
	}
	models.DB.Create(&subscription)
}

func CreateOrActivateSubscription(user_id uint, feed_id uint) {
	subscription := models.Subscription{
		UserID: user_id,
		FeedID: feed_id,
	}
	models.DB.FirstOrCreate(&subscription, subscription)
	subscription.IsActive = true
	models.DB.Save(&subscription)
}

func Unsubscribe(user_id uint, feed_id uint) {
	var subscription models.Subscription
	models.DB.Where("user_id = ? AND feed_id = ?", user_id, feed_id).First(&subscription)
	subscription.IsActive = false
	models.DB.Save(&subscription)
}

func GetSubscriptions(user_id uint) []models.Subscription {
	var subscriptions []models.Subscription
	models.DB.Where("user_id = ?", user_id).Preload("Feed").Find(&subscriptions)
	return subscriptions
}

func GetUsers() []models.User {
	var users []models.User
	models.DB.Find(&users)
	return users
}

func GetUnseenItems(user_id uint) []models.Item {
	var items []models.Item
	last24Hours := time.Now().Add(-24 * time.Hour)

	// get new items from last 24 hours
	models.DB.Where("subscriptions.user_id = ? AND seen_items.item_id IS NULL AND items.created_at >= ?", user_id, last24Hours).
		Joins("JOIN feeds ON items.feed_id = feeds.id").
		Joins("JOIN subscriptions ON feeds.id = subscriptions.feed_id AND subscriptions.user_id = ? AND subscriptions.is_active = TRUE", user_id).
		Joins("LEFT OUTER JOIN seen_items ON items.id = seen_items.item_id AND seen_items.user_id = ?", user_id).
		Find(&items)

	// add old unseen items every 3 hours between 8am and 9pm
	now := time.Now()
	hour := now.Hour()
	if hour >= 8 && hour <= 21 && hour%3 == 0 && now.Minute() < 5 {
		var oldItems []models.Item
		models.DB.Where("subscriptions.user_id = ? AND seen_items.item_id IS NULL AND items.created_at >= subscriptions.created_at", user_id).
			Joins("JOIN feeds ON items.feed_id = feeds.id").
			Joins("JOIN subscriptions ON feeds.id = subscriptions.feed_id AND subscriptions.user_id = ? AND subscriptions.is_active = TRUE", user_id).
			Joins("LEFT OUTER JOIN seen_items ON items.id = seen_items.item_id AND seen_items.user_id = ?", user_id).
			Order("items.created_at DESC").
			Limit(2).
			Find(&oldItems)

		items = append(items, oldItems...)
	}

	return items
}

func MarkItemsAsSeen(user_id uint, item_ids []uint) {
	for _, item_id := range item_ids {
		seenItem := models.SeenItem{
			ItemID: item_id,
			UserID: user_id,
		}
		models.DB.Create(&seenItem)
	}
}
