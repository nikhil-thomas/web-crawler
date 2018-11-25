package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nikhil-thomas/web-crawler/internal/crawlers/concurrent"
	"github.com/nikhil-thomas/web-crawler/internal/platform/http"
	"github.com/nikhil-thomas/web-crawler/internal/sitemap"
	"github.com/spf13/viper"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Printf("usage %s url", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	viper.Set("WORKER_COUNT", 10)
	viper.Set("PAGE_LIMIT", 1000)
	viper.Set("CRAWLER_TIMEOUT", "5s")
	viper.Set("CRAWLER_QUEUE_LENGTH", 500)
	urls := os.Args[1:]

	fetcher := http.NewFetcher()
	crwlMng := concurrent.NewCrawlManager(fetcher)

	siteMap := sitemap.NewSiteManager(urls[0], 100, crwlMng)

	siteMap.Crawl()

	siteMap.PrintMap()
}
