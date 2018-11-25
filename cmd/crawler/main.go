package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nikhil-thomas/web-crawler/internal/crawlers/simple"
	"github.com/nikhil-thomas/web-crawler/internal/platform/http"
	"github.com/nikhil-thomas/web-crawler/internal/sitemap"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Printf("usage %s url", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	urls := os.Args[1:]

	fetcher := http.NewFetcher()
	crwlMng := simple.NewCrawlManager(fetcher)

	siteMap := sitemap.NewSiteManager(urls[0], crwlMng)

	siteMap.Crawl()

	siteMap.PrintMap()
}
