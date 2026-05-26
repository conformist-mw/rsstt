package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"rsstt/logger"
	"rsstt/models"
	"rsstt/service"
	"rsstt/web"
)

func main() {
	config, err := LoadConfig()
	logger.Init(config.LogLevel)
	logger.Log.Info("Starting application")
	if err != nil {
		logger.Log.Error("Error loading config", err)
		os.Exit(1)
	}

	models.ConnectDb(config.DatabasePath)

	service.SetServiceConfig(&service.ServiceConfig{
		HTTPTimeout:        config.HTTPTimeout,
		MaxRetries:         config.MaxRetries,
		MaxMessagesPerUser: config.MaxMessagesPerUser,
		MaxConcurrentFeeds: config.MaxConcurrentFeeds,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	// feed update loop
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				logger.Log.Errorf("Panic in feed update: %v", r)
			}
		}()
		ticker := time.NewTicker(time.Duration(config.FeedUpdateInterval) * time.Minute)
		defer ticker.Stop()
		logger.Log.Debug("Updating feeds (initial)")
		service.UpdateFeeds()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				logger.Log.Debug("Updating feeds")
				service.UpdateFeeds()
			}
		}
	}()

	// send loop
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				logger.Log.Errorf("Panic in send loop: %v", r)
			}
		}()
		ticker := time.NewTicker(time.Duration(config.SendToUsersInterval) * time.Minute)
		defer ticker.Stop()
		logger.Log.Debug("Sending to user (initial)")
		service.SendToUser(config.TelegramAdminId, config.TelegramBotURL())
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				logger.Log.Debug("Sending to user")
				service.SendToUser(config.TelegramAdminId, config.TelegramBotURL())
			}
		}
	}()

	// web admin
	wg.Add(1)
	go func() {
		defer wg.Done()
		mux := http.NewServeMux()
		web.RegisterHandlers(mux)
		addr := fmt.Sprintf(":%d", config.AdminPort)
		logger.Log.Infof("Web admin on http://localhost%s", addr)
		srv := &http.Server{Addr: addr, Handler: mux}
		go func() {
			<-ctx.Done()
			srv.Close()
		}()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Errorf("Web server error: %v", err)
		}
	}()

	sig := <-sigChan
	logger.Log.Infof("Received signal: %v, shutting down", sig)
	cancel()
	wg.Wait()
	logger.Log.Info("Shutdown complete")
}
