package models

import (
	"log"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDb(dbPath string) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  true,
		},
	)
	db, err := gorm.Open(sqlite.Open(dbPath+"?_journal_mode=WAL&_busy_timeout=5000"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect database")
	}
	DB = db

	if err := DB.AutoMigrate(&Feed{}, &Item{}, &SeenItem{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	createIndexes()
}

func createIndexes() {
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_feeds_active ON feeds(is_active, deleted_at)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_items_feed_created ON items(feed_id, created_at)")
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_items_created_at ON items(created_at)")
}
