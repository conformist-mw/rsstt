package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"rsstt/logger"
	"rsstt/models"
	"rsstt/repository"
)

type TelegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

var telegramRateLimiter = time.NewTicker(time.Second / 30)

func escapeMarkdown(text string) string {
	special := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	escaped := text
	for _, char := range special {
		escaped = strings.ReplaceAll(escaped, char, "\\"+char)
	}
	return escaped
}

func generateMessage(item *models.Item) (string, string) {
	title := ""
	if item.Title != nil {
		title = escapeMarkdown(*item.Title)
	}
	message := fmt.Sprintf("[%s](%s)", title, item.Link)

	if item.PublishedAt != nil {
		message += "\n" + escapeMarkdown("Published: "+item.PublishedAt.Format("2006-01-02"))
	} else {
		message += "\n" + escapeMarkdown("Created: "+item.CreatedAt.Format("2006-01-02"))
	}

	description := ""
	if item.Description != nil {
		remaining := 4096 - len(message) - 2
		if remaining > 0 {
			desc := regexp.MustCompile("<[^>]*>").ReplaceAllString(*item.Description, "")
			desc = html.UnescapeString(desc)
			desc = escapeMarkdown(desc)
			description = truncateUTF8(desc, remaining)
		}
	}
	return message, description
}

func truncateUTF8(s string, max int) string {
	if len(s) <= max {
		return s
	}
	cut := max - len("\\.\\.\\.")
	if cut <= 0 {
		return ""
	}
	// step back to a rune boundary
	for cut > 0 && !utf8.RuneStart(s[cut]) {
		cut--
	}
	return s[:cut] + "\\.\\.\\."
}

func SendItemToUser(item *models.Item, chatID string, botURL string) error {
	message, description := generateMessage(item)

	err := sendTelegramMessage(item.ID, chatID, message+"\n\n"+description, botURL)
	if err != nil {
		logger.Log.Infof("Retrying with simplified message for item %d", item.ID)
		err = sendTelegramMessage(item.ID, chatID, message, botURL)
	}
	return err
}

func sendTelegramMessage(itemID uint, chatID string, text string, botURL string) error {
	<-telegramRateLimiter.C

	tgMessage := TelegramMessage{
		ChatID:    chatID,
		Text:      text,
		ParseMode: "MarkdownV2",
	}

	maxRetries := 3
	if serviceConfig != nil {
		maxRetries = serviceConfig.MaxRetries
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := doSendTelegramMessage(itemID, tgMessage, botURL)
		if err == nil {
			return nil
		}
		if attempt < maxRetries {
			logger.Log.Warnf("Attempt %d/%d failed for item %d: %v", attempt, maxRetries, itemID, err)
			time.Sleep(time.Duration(attempt*2) * time.Second)
		} else {
			return fmt.Errorf("item %d: all %d attempts failed: %w", itemID, maxRetries, err)
		}
	}
	return nil
}

func doSendTelegramMessage(itemID uint, tgMessage TelegramMessage, botURL string) error {
	jsonData, err := json.Marshal(tgMessage)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	timeout := 30 * time.Second
	if serviceConfig != nil {
		timeout = time.Duration(serviceConfig.HTTPTimeout) * time.Second
	}

	client := &http.Client{Timeout: timeout}
	resp, err := client.Post(botURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		return fmt.Errorf("status %d: %v", resp.StatusCode, result)
	}
	return nil
}

func SendToUser(chatID string, botURL string) {
	limit := 10
	if serviceConfig != nil {
		limit = serviceConfig.MaxMessagesPerUser
	}

	items := repository.GetUnseenItems(limit)
	if len(items) == 0 {
		return
	}

	var seen []uint
	for _, item := range items {
		if err := SendItemToUser(&item, chatID, botURL); err == nil {
			seen = append(seen, item.ID)
		}
	}
	if len(seen) > 0 {
		repository.MarkItemsAsSeen(seen)
	}
}
