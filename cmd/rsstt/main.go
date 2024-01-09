package main

import (
	"flag"
	"fmt"

	"github.com/conformist-mw/rsstt/models"
	"github.com/conformist-mw/rsstt/repository"
	"github.com/conformist-mw/rsstt/service"
)

func main() {
	url := flag.String("url", "", "The URL of the feed to create")
	tgUserId := flag.String("tg_user_id", "", "The Telegram user ID to send the feed to")
	listFeeds := flag.Bool("list_feeds", false, "List all feeds")
	feedId := flag.Int("feed_id", 0, "The ID of the feed to update")
	subscribe := flag.Bool("subscribe", false, "Subscribe to a feed")
	flag.Parse()

	models.ConnectDb()

	if url != nil && *url != "" {
		service.CreateFeed(*url)
	}

	if tgUserId != nil && *tgUserId != "" {
		if subscribe != nil && *subscribe {
			user := repository.GetUserByTgChatId(*tgUserId)
			if feedId != nil && *feedId != 0 {
				repository.CreateSubscription(user.ID, uint(*feedId))
			}
		} else {
			repository.CreateUser(*tgUserId)
		}
	}

	if listFeeds != nil && *listFeeds {
		feeds := repository.FetchAllFeeds()
		for _, feed := range feeds {
			fmt.Printf("id: %d url: %s\n", feed.ID, feed.Url)
		}
	}
}
