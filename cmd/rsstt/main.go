package main

import (
	"flag"
	"fmt"

	"github.com/conformist-mw/rsstt/models"
	"github.com/conformist-mw/rsstt/service"
)

func main() {
	url := flag.String("url", "", "The URL of the feed to create")
	flag.Parse()

	if *url == "" {
		fmt.Println("Please provide a URL using the -url flag")
		return
	}

	models.ConnectDb()
	service.CreateFeed(*url)
}
