package simple

import (
	"fmt"
	"strings"

	"github.com/nikhil-thomas/web-crawler/internal/crawlers"
)

// CrawlManager implements sitemap.Crawler interface
type CrawlManager struct {
	fetcher crawlers.URLFetcher
}

// NewCrawlManager creates and returns a CrawlManager
func NewCrawlManager(fetcher crawlers.URLFetcher) *CrawlManager {
	return &CrawlManager{fetcher: fetcher}
}

// Crawl collects links form a page and returns a list of links
// Crawl skips links which are not from root domain
func (cm *CrawlManager) Crawl(url, rootDomain string) ([]string, error) {
	var children []string
	links, err := cm.fetcher.ExtractURLs(url)
	if err != nil {
		return nil, fmt.Errorf("crawl manager: %s", err)
	}

	for _, link := range links {
		if strings.HasPrefix(link, rootDomain) {
			children = append(children, link)
			fmt.Println("add:", link)
		} else {
			fmt.Println("skip:", link)
		}
	}

	return children, nil
}
