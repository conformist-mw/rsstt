package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/conformist-mw/rsstt/models"
	"github.com/conformist-mw/rsstt/repository"
)

type TelegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func generateMessage(item *models.Item) (string, error) {
	markup := `
	<b><a href="{{ .Link }}">{{ .Title }}</a></b>

	{{ .Description }}
	`
	tmpl, err := template.New("").Parse(markup)
	if err != nil {
		return "", err
	}
	var rendered bytes.Buffer
	err = tmpl.Execute(&rendered, item)
	if err != nil {
		return "", err
	}

	return rendered.String(), nil
}

func SendItemToUser(item *models.Item, tgChatId string, botUrl string) error {
	msg, err := generateMessage(item)
	if err != nil {
		return err
	}

	tgMessage := TelegramMessage{
		ChatID:    tgChatId,
		Text:      msg,
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
