package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/conformist-mw/rsstt/logger"
	"github.com/conformist-mw/rsstt/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func setCommands(bot *tgbotapi.BotAPI) {
	commands := []tgbotapi.BotCommand{
		{Command: "help", Description: "Show help"},
		{Command: "feeds", Description: "List all feeds"},
		{Command: "subs", Description: "List all subscriptions"},
		{Command: "subscribe", Description: "Subscribe to feed (pass feed url)"},
		{Command: "unsubscribe", Description: "Unsubscribe from feed (pass feed url)"},
	}
	if _, err := bot.Request(tgbotapi.SetMyCommandsConfig{
		Commands: commands,
	}); err != nil {
		log.Panic(err)
	}
}

func CreateBot(token string) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	setCommands(bot)
	return bot
}

func SetupWebhook(bot *tgbotapi.BotAPI, url string) {
	wh, _ := tgbotapi.NewWebhook(url)
	if _, err := bot.Request(wh); err != nil {
		log.Fatal(err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		logger.Log.Error("Error getting webhook info", err)
	}

	if info.LastErrorDate != 0 {
		logger.Log.Error("Telegram callback failed: ", info.LastErrorMessage)
	}
}

func GetUpdates(bot *tgbotapi.BotAPI, url string) tgbotapi.UpdatesChannel {
	logger.Log.Debugf("Setup updates channel. URL: %s", url)
	if url != "" {
		SetupWebhook(bot, url)
		updates := bot.ListenForWebhook("/")
		logger.Log.Debug("Start HTTP server")
		go func() {
			err := http.ListenAndServe(":8080", nil)
			if err != nil {
				logger.Log.Error("Failed to start HTTP server: ", err)
			}
		}()
		return updates
	}

	logger.Log.Debug("Start polling updates")
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return bot.GetUpdatesChan(u)
}

func ProcessUpdates(ctx context.Context, bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel, adminChatID int64) {
	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("Stopping updates processing")
			return
		case update := <-updates:
			if update.Message == nil {
				continue
			}

			tgUserId := update.Message.From.ID
			logger.Log.Debugf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if tgUserId != int64(adminChatID) {
				continue
			}

			if !update.Message.IsCommand() {
				continue
			}

			handleCommand(bot, update, tgUserId)
		}
	}
}

func handleCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update, tgUserId int64) {
	logger.Log.Debugf("Handling command: %s", update.Message.Command())
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.DisableWebPagePreview = true

	switch update.Message.Command() {
	case "help":
		msg.Text = "I can help you with RSS feeds. /subscribe feed_url, /unsubscribe feed_url, /feeds, /subs"
	case "feeds":
		feeds := repository.FetchAllFeeds()
		for _, feed := range feeds {
			msg.Text += strconv.Itoa(int(feed.ID)) + " " + feed.Link + " " + feed.Url + "\n"
		}
	case "subs":
		user := repository.GetUserByTgChatId(strconv.Itoa(int(tgUserId)))
		subs := repository.GetSubscriptions(user.ID)
		for _, sub := range subs {
			isActive := "no"
			if sub.IsActive {
				isActive = "yes"
			}
			msg.Text += strconv.Itoa(int(sub.ID)) + " " + sub.Feed.Link + ". Active: " + isActive + "\n"
		}
	case "subscribe":
		feedUrl := update.Message.CommandArguments()
		if len(feedUrl) < 1 {
			msg.Text = "Please provide feed URL"
			break
		}
		msg.Text = fmt.Sprintf("Subscribing to %s", feedUrl)
		feed, error := CreateFeed(feedUrl)
		if error != nil {
			msg.Text = fmt.Sprintf("Failed to subscribe to %s, invalid URL", feedUrl)
			break
		}
		user := repository.GetUserByTgChatId(strconv.Itoa(int(tgUserId)))
		repository.CreateOrActivateSubscription(user.ID, feed.ID)
		msg.Text = fmt.Sprintf("Subscribed to %s", feedUrl)
	case "unsubscribe":
		feedUrl := update.Message.CommandArguments()
		if len(feedUrl) < 1 {
			msg.Text = "Please provide feed URL"
			break
		}
		msg.Text = fmt.Sprintf("Unsubscribing from %s", feedUrl)
		feed := repository.GetFeedByUrl(feedUrl)
		if feed.ID == 0 {
			msg.Text = fmt.Sprintf("Failed to unsubscribe from %s, not subscribed", feedUrl)
			break
		}
		user := repository.GetUserByTgChatId(strconv.Itoa(int(tgUserId)))
		repository.Unsubscribe(user.ID, feed.ID)
		msg.Text = fmt.Sprintf("Unsubscribed from %s", feedUrl)
	default:
		msg.Text = "I don't know that command"
	}

	if _, err := bot.Send(msg); err != nil {
		logger.Log.Error("Error sending message", err)
	}
}
