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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {

	url := parseFlags()

	log.Info("root   : ", url)

	fetcher := http.NewFetcher()

	var crwlMng sitemap.Crawler

	crwlMng = concurrent.NewCrawlManager(fetcher)

	if viper.GetBool("DISABLE_CONCURRENCY") {
		crwlMng = simple.NewCrawlManager(fetcher)
	}

	siteMap := sitemap.NewSiteManager(url, crwlMng)
	siteMap.Crawl()
	siteMap.PrintMap()
}

func parseFlags() string {
	disableConcurrency := flag.Bool(
		"con-off",
		false,
		"set false to turn off concurrency")

	pageLimit := flag.Int(
		"p",
		250,
		"maximum number of pages to be crawled (set 0 for no limit)")

	linksPerPage := flag.Int(
		"l",
		100,
		"maximum number of links to be extracted per page to be crawled (set 0 for no limit)")

	crawlerTimeout := flag.String(
		"t",
		"5s",
		"timeout to stop concurrent crawler when no new links are available [eg: 1s,1ns,1ms,1µs]")

	logLevel := flag.Int(
		"log",
		4,
		"log level [0-panic, 1-fatal, 2-error, 3-warn, 4-info, 5-debug")

	queueLength := flag.Int(
		"q",
		500,
		"length of crawler process queue")

	numWorkers := flag.Int(
		"w",
		10,
		"number of workers(goroutines) in concurrent crawling")

	trimRoot := flag.Bool(
		"trim",
		false,
		"trim root domain name from sitemap")
	flag.Parse()

	viper.Set("DISABLE_CONCURRENCY", *disableConcurrency)
	viper.Set("PAGE_LIMIT", *pageLimit)
	viper.Set("LINKS_PER_PAGE", *linksPerPage)
	viper.Set("CRAWLER_TIMEOUT", *crawlerTimeout)
	viper.Set("WORKER_COUNT", *numWorkers)
	viper.Set("CRAWLER_QUEUE_LENGTH", *queueLength)
	viper.Set("TRIM_ROOT", *trimRoot)

	log.SetLevel(log.Level(*logLevel))

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
	return url.String()
}
