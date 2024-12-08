package main

import (
	"time"

	"github.com/conformist-mw/rsstt/models"
	"github.com/conformist-mw/rsstt/service"
)

func updateFeeds() {
	service.UpdateFeeds()
	time.Sleep(5 * time.Minute)
}

func sendToAllUsers(botUrl string) {
	service.SendToAllUsers(botUrl)
	time.Sleep(5 * time.Minute)
}

func handleUpdates(token string, url string, adminChatId int64) {
	bot, updates := service.GetBotAndUpdates(token, url)
	service.HandleUpdates(bot, updates, adminChatId)
}

func main() {
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	models.ConnectDb(config.DatabasePath)
	go updateFeeds()
	go sendToAllUsers(config.TelegramBotURL())
	go handleUpdates(config.TelegramBotToken, config.TelegramBotUrl, config.TelegramAdminChatID())
	select {}
}
