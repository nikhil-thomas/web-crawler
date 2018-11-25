package sitemap

import (
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

// Crawler interface defines the behavior of a Crawler
type Crawler interface {
	Crawl(url string) (map[string]Children, error)
}

// Children defines a list of children links in a html page
type Children []string

// SiteMapManager defines a sitemap generator
// SiteMapManager can work on a link and generate its sitemap
type SiteMapManager struct {
	rootDomain string
	Sitemap    map[string]Children
	urlQueue   []string
	crawler    Crawler
}

// NewSiteManager creates and returns a SiteMapManager
func NewSiteManager(url string, crawler Crawler) *SiteMapManager {
	return &SiteMapManager{
		rootDomain: url,
		Sitemap:    map[string]Children{},
		urlQueue:   []string{url},
		crawler:    crawler,
	}
}

// Crawl crawls a site starting from specified root url
// Crawl popolates the Sitemap map[string]Children
func (sm *SiteMapManager) Crawl() {
	var err error
	sm.Sitemap, err = sm.crawler.Crawl(sm.rootDomain)
	if err != nil {
		log.Error("sitemap : ", err)
	}
}

// PrintMap prints site map as a tree
func (sm *SiteMapManager) PrintMap() {
	sm.FPrintMap(os.Stdout)
}

// FPrintMap writes site map as a tree to io.Writer
func (sm *SiteMapManager) FPrintMap(w io.Writer) {
	fmt.Fprintf(w, "\n::::: Site Map: %s ::::\n\n", sm.rootDomain)

	sm.printTree(w, sm.rootDomain, 0)
}

func (sm *SiteMapManager) printTree(w io.Writer, url string, depth int) {

	fmt.Fprintf(w, "%*s%s\n", depth, "", url)
	for _, val := range sm.Sitemap[url] {
		sm.printTree(w, val, depth+2)
	}
}
