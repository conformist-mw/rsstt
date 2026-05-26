package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DatabasePath     string `json:"DATABASE_PATH"`
	TelegramBotToken string `json:"TELEGRAM_BOT_TOKEN"`
	TelegramAdminId  string `json:"TELEGRAM_ADMIN_ID"`
	LogLevel         string `json:"LOG_LEVEL"`
	AdminPort        int    `json:"ADMIN_PORT"`

	FeedUpdateInterval  int `json:"FEED_UPDATE_INTERVAL"`
	SendToUsersInterval int `json:"SEND_TO_USERS_INTERVAL"`
	HTTPTimeout         int `json:"HTTP_TIMEOUT"`
	MaxRetries          int `json:"MAX_RETRIES"`
	MaxConcurrentFeeds  int `json:"MAX_CONCURRENT_FEEDS"`
	MaxMessagesPerUser  int `json:"MAX_MESSAGES_PER_USER"`
}

func (c *Config) TelegramBotURL() string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.TelegramBotToken)
}

func LoadConfig() (Config, error) {
	var config Config

	// defaults
	config.DatabasePath = "data/feed.db"
	config.AdminPort = 8080
	config.FeedUpdateInterval = 5
	config.SendToUsersInterval = 5
	config.HTTPTimeout = 30
	config.MaxRetries = 3
	config.MaxConcurrentFeeds = 5
	config.MaxMessagesPerUser = 10

	// try config.json first
	if file, err := os.Open("config.json"); err == nil {
		defer file.Close()
		if err := json.NewDecoder(file).Decode(&config); err != nil {
			return config, err
		}
	}

	// env vars always override (partial overrides are fine)
	stringVars := map[string]*string{
		"DATABASE_PATH":      &config.DatabasePath,
		"TELEGRAM_BOT_TOKEN": &config.TelegramBotToken,
		"TELEGRAM_ADMIN_ID":  &config.TelegramAdminId,
		"LOG_LEVEL":          &config.LogLevel,
	}
	for envVar, field := range stringVars {
		if value, exists := os.LookupEnv(envVar); exists {
			*field = value
		}
	}

	intVars := map[string]*int{
		"ADMIN_PORT":             &config.AdminPort,
		"FEED_UPDATE_INTERVAL":   &config.FeedUpdateInterval,
		"SEND_TO_USERS_INTERVAL": &config.SendToUsersInterval,
		"HTTP_TIMEOUT":           &config.HTTPTimeout,
		"MAX_RETRIES":            &config.MaxRetries,
		"MAX_CONCURRENT_FEEDS":   &config.MaxConcurrentFeeds,
		"MAX_MESSAGES_PER_USER":  &config.MaxMessagesPerUser,
	}
	for envVar, field := range intVars {
		if value, exists := os.LookupEnv(envVar); exists {
			if parsed, err := strconv.Atoi(value); err == nil {
				*field = parsed
			}
		}
	}

	return config, nil
}
