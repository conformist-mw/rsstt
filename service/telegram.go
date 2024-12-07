package service

import (
	"log"
	"fmt"
	"net/http"
	"strconv"

	"github.com/conformist-mw/rsstt/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func getBot(token string) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	setCommands(bot)
	return bot
}

func getPollingBot(token string) (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel) {
	bot := getBot(token)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	return bot, updates
}

func getWebhookBot(token string, url string) (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel) {
	bot := getBot(token)
	wh, _ := tgbotapi.NewWebhook(url)
	_, err := bot.Request(wh)
	if err != nil {
		log.Fatal(err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()
	return bot, updates
}

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

func GetBotAndUpdates(token string, url string) (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel) {
	if url == "" {
		return getPollingBot(token)
	}
	return getWebhookBot(token, url)
}

func HandleUpdates(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel, adminChatID int64) {
	for update := range updates {
		if update.Message == nil {
			continue
		}
		tgUserId := update.Message.From.ID
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if tgUserId != int64(adminChatID) {
			continue
		}
		if !update.Message.IsCommand() {
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.DisableWebPagePreview = true

		switch update.Message.Command() {
		case "help":
			msg.Text = "I can help you with rss feeds. /subscribe feed_url, /unsubscribe feed_url, /feeds, /subs"
			fmt.Println(update.Message.CommandArguments())
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
				msg.Text = "Please provide feed url"
				break
			}
			msg.Text = fmt.Sprintf("Subscribing to %s", feedUrl)
			feed, error := CreateFeed(feedUrl)
			if error != nil {
				msg.Text = fmt.Sprintf("Failed to subscribe to %s, invalid url", feedUrl)
				break
			}
			user := repository.GetUserByTgChatId(strconv.Itoa(int(tgUserId)))
			repository.CreateOrActivateSubscription(user.ID, feed.ID)
			msg.Text = fmt.Sprintf("Subscribed to %s", feedUrl)
		case "unsubscribe":
			feedUrl := update.Message.CommandArguments()
			if len(feedUrl) < 1 {
				msg.Text = "Please provide feed url"
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
			log.Panic(err)
		}
	}
}
