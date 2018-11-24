package sitemap

import (
	"fmt"
	"log"
)

// Crawler interface defines the behavior of a Crawler
type Crawler interface {
	Crawl(url string, pageLimit, linksPerPage int) (map[string]Children, error)
}

// Children defines a list of children links in a html page
type Children []string

// SiteMapManager defines a sitemap generator
// SiteMapManager can work on a link and generate its sitemap
type SiteMapManager struct {
	rootDomain   string
	Sitemap      map[string]Children
	urlQueue     []string
	crawler      Crawler
	pageLimit    int
	linksPerPage int
}

// NewSiteManager creates and returns a SiteMapManager
func NewSiteManager(url string, depth int, crawler Crawler) *SiteMapManager {
	return &SiteMapManager{
		rootDomain:   url,
		Sitemap:      map[string]Children{},
		urlQueue:     []string{url},
		crawler:      crawler,
		pageLimit:    50,
		linksPerPage: 2,
	}
}

// Crawl crawls a site starting from specified root url
// Crawl popolates the Sitemap map[string]Children
func (sm *SiteMapManager) Crawl() {
	var err error
	sm.Sitemap, err = sm.crawler.Crawl(sm.rootDomain, sm.pageLimit, sm.linksPerPage)
	if err != nil {
		log.Printf("sitemap: error: %s", err)
	}
}

// PrintMap prints site map as a tree
func (sm *SiteMapManager) PrintMap() {
	fmt.Printf("\n::::: Site Map: %s ::::\n", sm.rootDomain)
	sm.printTree(sm.rootDomain, 0)
}

func (sm *SiteMapManager) printTree(url string, depth int) {

	fmt.Printf("%*s%s\n", depth, "", url)
	for _, val := range sm.Sitemap[url] {
		sm.printTree(val, depth+2)
	}
}
