package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDb() {
	db, err := gorm.Open(sqlite.Open("feed.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Feed{}, &Item{}, &User{}, &Subscription{}, &SeenItem{})
	DB = db
}
