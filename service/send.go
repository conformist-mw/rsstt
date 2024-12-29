package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"regexp"
	"strings"

	"rsstt/logger"
	"rsstt/models"
	"rsstt/repository"
)

type TelegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func escapeMarkdown(text string) string {
	markdownSpecial := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	escaped := text
	for _, char := range markdownSpecial {
		escaped = strings.ReplaceAll(escaped, char, "\\"+char)
	}
	return escaped
}

func generateMessage(item *models.Item) string {
	title := escapeMarkdown(*item.Title)
	message := fmt.Sprintf("[%s](%s)", title, item.Link)

	var dateStr string
	if item.PublishedAt != nil {
		dateStr = fmt.Sprintf("Дата публикации: %s", item.PublishedAt.Format("2006-01-02 15:04:05"))
	} else {
		dateStr = fmt.Sprintf("Дата создания: %s", item.CreatedAt.Format("2006-01-02 15:04:05"))
	}
	message += "\n" + escapeMarkdown(dateStr)

	if item.Description != nil {
		remainingLength := 4096 - len(message) - 2
		if remainingLength > 0 {
			description := regexp.MustCompile("<[^>]*>").ReplaceAllString(*item.Description, "")
			description = html.UnescapeString(description)
			description = escapeMarkdown(description)

			if len(description) > remainingLength {
				description = description[:remainingLength-3] + "..."
			}

			message += "\n\n" + description
		}
	}
	return message
}

func SendItemToUser(item *models.Item, tgChatId string, botUrl string) error {
	tgMessage := TelegramMessage{
		ChatID:    tgChatId,
		Text:      generateMessage(item),
		ParseMode: "MarkdownV2",
	}
	logger.Log.Debugf("Sending telegram message. Item: %d, Message: %s", item.ID, tgMessage.Text)
	jsonData, err := json.Marshal(tgMessage)
	if err != nil {
		logger.Log.Errorf("Error marshalling telegram message. Item: %d, Error: %s", item.ID, err)
		return err
	}

	resp, err := http.Post(botUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Log.Errorf("Error sending telegram message. Item: %d, Error: %s", item.ID, err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			logger.Log.Errorf("Failed to decode error response. Status code: %d, Item: %d, Error: %s", resp.StatusCode, item.ID, err)
		} else {
			logger.Log.Errorf("Failed to send message. Status code: %d, Item: %d, Response: %v", resp.StatusCode, item.ID, result)
		}
		return fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
	}

	return nil
}

func SendToAllUsers(botUrl string) {
	for _, user := range repository.GetUsers() {
		items := repository.GetUnseenItems(user.ID)
		if len(items) > 0 {
			seenItems := []uint{}
			for _, item := range items {
				err := SendItemToUser(&item, *user.TgChatId, botUrl)
				if err == nil {
					seenItems = append(seenItems, item.ID)
				}
			}
			if len(seenItems) > 0 {
				repository.MarkItemsAsSeen(user.ID, seenItems)
			}
		}
	}
}
