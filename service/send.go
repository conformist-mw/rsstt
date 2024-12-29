package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"rsstt/models"
	"rsstt/repository"
)

type TelegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func generateMessage(item *models.Item) string {
	message := fmt.Sprintf("<b><a href=\"%s\">%s</a></b>", item.Link, *item.Title)
	if item.Description != nil {
		if len(message)+len(*item.Description) <= 4096-2 {
			message += "\n\n" + *item.Description
		}
	}
	return message
}

func SendItemToUser(item *models.Item, tgChatId string, botUrl string) error {
	tgMessage := TelegramMessage{
		ChatID:    tgChatId,
		Text:      generateMessage(item),
		ParseMode: "HTML",
	}

	jsonData, err := json.Marshal(tgMessage)
	if err != nil {
		return err
	}

	resp, err := http.Post(botUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
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
