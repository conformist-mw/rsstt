package models

import "gorm.io/gorm"

// TODO: make indexes and constraints
type User struct {
	gorm.Model
	Username      *string        `gorm:"type:varchar(255)" json:"username"`
	FirstName     *string        `gorm:"type:varchar(255)" json:"first_name"`
	LastName      *string        `gorm:"type:varchar(255)" json:"last_name"`
	TgChatId      *string        `gorm:"type:varchar(255);unique" json:"tg_chat_id"`
	Subscriptions []Subscription `json:"subscriptions"`
}

type Subscription struct {
	gorm.Model
	UserID   uint `gorm:"uniqueIndex:idx_user_sub" json:"user_id"`
	FeedID   uint `gorm:"uniqueIndex:idx_user_sub" json:"feed_id"`
	IsActive bool `json:"is_active"`
}

type SeenItem struct {
	ID     uint `gorm:"primarykey"`
	ItemID uint `gorm:"uniqueIndex:idx_seen_item" json:"item_id"`
	UserID uint `gorm:"uniqueIndex:idx_seen_item" json:"user_id"`
}
