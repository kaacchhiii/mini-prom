package main

import (
	"log"
	"mini-prometheus/api"
	"mini-prometheus/config"
	"mini-prometheus/scraper"
	"mini-prometheus/storage"
	"time"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	store := storage.NewMemoryStorage()

	// Start retention policy worker
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			if cfg.Storage.Retention > 0 {
				store.Prune(cfg.Storage.Retention)
			}
		}
	}()

	for _, target := range cfg.Targets {
		go scraper.StartScraper(target.URL, target.Interval, store)
	}

	port := cfg.Server.Port
	if port == 0 {
		port = 9090 // Default if not specified
	}
	api.StartServer(store, port)
}
