package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Feed struct {
	gorm.Model
	Url           string       `gorm:"type:varchar(255);uniqueIndex" json:"url"`
	Title         *string      `gorm:"type:varchar(255)" json:"title"`
	Description   *string      `json:"description"`
	Link          string       `json:"link"`
	FeedLink      *string      `json:"feed_link"`
	FeedType      string       `json:"feed_type"`
	FeedVersion   string       `json:"feed_version"`
	LastFetchedAt sql.NullTime `json:"last_fetched_at"`
	IsActive      bool         `gorm:"default:true" json:"is_active"`
	Items         []Item       `json:"items"`
}

type Item struct {
	gorm.Model
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Content     *string    `json:"content"`
	Link        string     `json:"link"`
	FeedID      uint       `json:"feed_id"`
	PublishedAt *time.Time `json:"published_at"`
}
