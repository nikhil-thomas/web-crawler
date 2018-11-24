package sitemap

import (
	"fmt"
)

// Crawler interface defines the behavior of a Crawler
type Crawler interface {
	Crawl(url string, root string) ([]string, error)
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
		pageLimit:    100,
		linksPerPage: 2,
	}
}

// Crawl crawls a site starting from specified root url
// Crawl popolates the Sitemap map[string]Children
func (sm *SiteMapManager) Crawl() {
	i := 0
	for len(sm.urlQueue) > 0 {
		url := sm.urlQueue[0]
		links, _ := sm.crawler.Crawl(url, sm.rootDomain)

		k := 0
		for _, link := range links {
			if _, ok := sm.Sitemap[link]; !ok {
				sm.Sitemap[url] = append(sm.Sitemap[url], link)
				sm.Sitemap[link] = Children{}
				sm.urlQueue = append(sm.urlQueue, link)
				k++
			}
			if sm.linksPerPage > 0 && k >= sm.linksPerPage {
				break
			}
		}

		sm.urlQueue = sm.urlQueue[1:]
		i++
		if sm.pageLimit != 0 && i >= sm.pageLimit {
			break
		}
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
