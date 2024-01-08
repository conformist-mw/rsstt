package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/conformist-mw/rsstt/models"
	"github.com/conformist-mw/rsstt/service"
)

type Config struct {
	TelegramBotToken string `json:"telegram_bot_token"`
}

func (c *Config) TelegramBotURL() string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.TelegramBotToken)
}

func LoadConfig() (Config, error) {
	var config Config

	file, err := os.Open("config.json")
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func main() {
	models.ConnectDb()
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	service.UpdateFeeds()
	botUrl := config.TelegramBotURL()
	service.SendToAllUsers(botUrl)
}
