package repository

import "github.com/conformist-mw/rsstt/models"

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
	models.DB.Where("subscriptions.user_id = ? AND seen_items.item_id IS NULL", user_id).
		Joins("JOIN feeds ON items.feed_id = feeds.id").
		Joins("JOIN subscriptions ON feeds.id = subscriptions.feed_id").
		Joins("LEFT OUTER JOIN seen_items ON items.id = seen_items.item_id AND seen_items.user_id = ?", user_id).
		Find(&items)
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
