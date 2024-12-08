package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DatabasePath	 string `json:"DATABASE_PATH"`
	TelegramBotToken string `json:"TELEGRAM_BOT_TOKEN"`
	TelegramBotUrl   string `json:"TELEGRAM_BOT_URL"`
	TelegramAdminId  string `json:"TELEGRAM_ADMIN_ID"`
}

func (c *Config) TelegramBotURL() string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.TelegramBotToken)
}

func (c *Config) TelegramBotWebhookURL() string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook", c.TelegramBotToken)
}

func (c *Config) TelegramAdminChatID() int64 {
	chatID, _ := strconv.Atoi(c.TelegramAdminId)
	return int64(chatID)
}

func LoadConfig() (Config, error) {
	var config Config

	envVars := map[string]*string{
		"TELEGRAM_BOT_TOKEN": &config.TelegramBotToken,
		"TELEGRAM_BOT_URL":   &config.TelegramBotUrl,
		"TELEGRAM_ADMIN_ID":  &config.TelegramAdminId,
	}

	var loadedFromEnv bool

	for envVar, configField := range envVars {
		if value, exists := os.LookupEnv(envVar); exists {
			*configField = value
			loadedFromEnv = true
		}
	}

	if loadedFromEnv {
		return config, nil
	}

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
