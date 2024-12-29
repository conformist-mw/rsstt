package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"rsstt/logger"
	"rsstt/models"
	"rsstt/service"
)

func main() {
	config, err := LoadConfig()
	logger.Init(config.LogLevel)
	logger.Log.Info("Starting application")
	if err != nil {
		logger.Log.Error("Error loading config", err)
		panic(err)
	}
	models.ConnectDb(config.DatabasePath)
	bot := service.CreateBot(config.TelegramBotToken)
	updates := service.GetUpdates(bot, config.TelegramBotUrl)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Log.Debug("Start Processing updates")
		service.ProcessUpdates(ctx, bot, updates, config.TelegramAdminChatID())
		logger.Log.Debug("End Processing updates")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		updateTicker := time.NewTicker(5 * time.Second)
		defer updateTicker.Stop()

		for {
			select {
			case <-ctx.Done():
				logger.Log.Info("Stopping feeds update")
				return
			case <-updateTicker.C:
				logger.Log.Debug("Start Updating feeds")
				service.UpdateFeeds()
				logger.Log.Debug("End Updating feeds")
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		sendTicker := time.NewTicker(5 * time.Second)
		defer sendTicker.Stop()

		for {
			select {
			case <-ctx.Done():
				logger.Log.Info("Stopping sending to users")
				return
			case <-sendTicker.C:
				logger.Log.Debug("Start Sending to users")
				service.SendToAllUsers(config.TelegramBotURL())
				logger.Log.Debug("End Sending to users")
			}
		}
	}()

	sig := <-sigChan
	logger.Log.Infof("Received signal: %v", sig)

	cancel()

	logger.Log.Info("Waiting for all operations to complete...")
	wg.Wait()

	logger.Log.Info("Application shutdown complete")
}
