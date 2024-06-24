package service

import (
	"log"
	"net/http"

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
		{Command: "sayhi", Description: "Say hi"},
		{Command: "status", Description: "Show status"},
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
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if update.Message.Chat.ID != int64(adminChatID) {
			continue
		}
		if !update.Message.IsCommand() {
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Command() {
		case "help":
			msg.Text = "I understand /sayhi and /status."
		case "sayhi":
			msg.Text = "Hi :)"
		case "status":
			msg.Text = "I'm ok."
		default:
			msg.Text = "I don't know that command"
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
