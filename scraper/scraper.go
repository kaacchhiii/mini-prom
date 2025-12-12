package scraper

import (
	"bufio"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"mini-prometheus/storage"
)

func StartScraper(url string, interval time.Duration, store *storage.MemoryStorage) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		<-ticker.C
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error scraping %s: %v", url, err)
			continue
		}
		scanner := bufio.NewScanner(resp.Body)
		now := time.Now()
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "#") || line == "" {
				continue
			}
			parts := strings.Fields(line)
			if len(parts) != 2 {
				continue
			}
			value, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				continue
			}
			store.AddMetric(parts[0], now, value)
		}
		resp.Body.Close()
	}
}
