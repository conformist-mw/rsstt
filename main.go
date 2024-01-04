package main

import (
	"fmt"

	"github.com/mmcdole/gofeed"
)

func main() {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://www.opennet.me/opennews/opennews_all_utf.rss")
	fmt.Println(feed.Title)
	for _, item := range feed.Items {
		fmt.Println(item.Title)
		fmt.Println(item.Description)
		fmt.Println(item.Link)

	}
}
