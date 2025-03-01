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

func generateMessageAndDescription(item *models.Item) (string, string) {
	title := escapeMarkdown(*item.Title)
	message := fmt.Sprintf("[%s](%s)", title, item.Link)
	
	var dateStr string
	if item.PublishedAt != nil {
		dateStr = fmt.Sprintf("Published: %s", item.PublishedAt.Format("2006-01-02"))
	} else {
		dateStr = fmt.Sprintf("Created: %s", item.CreatedAt.Format("2006-01-02"))
	}
	message += "\n" + escapeMarkdown(dateStr)
	description := ""
	
	if item.Description != nil {
		remainingLength := 4096 - len(message) - 2
		if remainingLength > 0 {
			description = regexp.MustCompile("<[^>]*>").ReplaceAllString(*item.Description, "")
			description = html.UnescapeString(description)
			description = escapeMarkdown(description)

			if len(description) > remainingLength {
				description = description[:remainingLength-3] + "..."
			}
		}
	}
	return message, description
}

func SendItemToUser(item *models.Item, tgChatId string, botUrl string) error {
	message, description := generateMessageAndDescription(item)
	
	fullMessage := message + "\n\n" + description
	err := sendTelegramMessage(item.ID, tgChatId, fullMessage, botUrl)
	
	if err != nil {
		logger.Log.Infof("Trying to send simplified message. Item: %d", item.ID)
		err = sendTelegramMessage(item.ID, tgChatId, message, botUrl)
	}
	
	return err
}

func sendTelegramMessage(itemID uint, chatID string, text string, botUrl string) error {
	tgMessage := TelegramMessage{
		ChatID:    chatID,
		Text:      text,
		ParseMode: "MarkdownV2",
	}
	logger.Log.Debugf("Sending telegram message. Item: %d, Message: %s", itemID, tgMessage.Text)
	jsonData, err := json.Marshal(tgMessage)
	if err != nil {
		logger.Log.Errorf("Error marshalling telegram message. Item: %d, Error: %s", itemID, err)
		return err
	}

	resp, err := http.Post(botUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Log.Errorf("Error sending telegram message. Item: %d, Error: %s", itemID, err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			logger.Log.Errorf("Failed to decode error response. Status code: %d, Item: %d, Error: %s", resp.StatusCode, itemID, err)
		} else {
			logger.Log.Errorf("Failed to send message. Status code: %d, Item: %d, Response: %v", resp.StatusCode, itemID, result)
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
