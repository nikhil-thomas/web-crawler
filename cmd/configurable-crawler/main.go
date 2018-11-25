package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/nikhil-thomas/web-crawler/internal/crawlers/concurrent"
	"github.com/nikhil-thomas/web-crawler/internal/crawlers/simple"
	"github.com/nikhil-thomas/web-crawler/internal/platform/http"
	"github.com/nikhil-thomas/web-crawler/internal/sitemap"
	"github.com/spf13/viper"
)

func main() {

	disableConcurrency := flag.Bool("con-off", false, "set false to turn off concurrency")
	pageLimit := flag.Int("p", 50, "maximum number of pages to be crawled (set 0 for no limit)")
	linksPerPage := flag.Int("l", 5, "maximum number of links to be extracted per page to be crawled (set 0 for no limit)")
	crawlerTimeout := flag.Int("t", 10, "timeout to stop crawler when no new links are available (only for concurrent crawler)")

	flag.Parse()

	viper.Set("DISABLE_CONCURRENCY", *disableConcurrency)
	viper.Set("PAGE_LIMIT", *pageLimit)
	viper.Set("LINKS_PER_PAGE", *linksPerPage)
	viper.Set("CRAWLER_TIMEOUT", *crawlerTimeout)

	args := flag.Args()
	if len(args) < 1 {
		fmt.Printf("\nusage %s <options> url\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(1)
	}

	url, err := url.Parse(args[0])
	if err != nil {
		fmt.Printf("url parse error: %s\n", err)
		os.Exit(1)
	}

	fetcher := http.NewFetcher()

	var crwlMng sitemap.Crawler
	crwlMng = concurrent.NewCrawlManager(fetcher)

	if viper.GetBool("DISABLE_CONCURRENCY") {
		crwlMng = simple.NewCrawlManager(fetcher)
	}

	siteMap := sitemap.NewSiteManager(url.String(), 100, crwlMng)

	siteMap.Crawl()

	siteMap.PrintMap()
}
