package simple

import (
	"fmt"
	"strings"

	"github.com/nikhil-thomas/web-crawler/internal/crawlers"
	"github.com/nikhil-thomas/web-crawler/internal/sitemap"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// CrawlManager implements sitemap.Crawler interface
type CrawlManager struct {
	fetcher crawlers.URLFetcher
}

// NewCrawlManager creates and returns a CrawlManager
func NewCrawlManager(fetcher crawlers.URLFetcher) *CrawlManager {
	return &CrawlManager{fetcher: fetcher}
}

func filterDomains(links []string, rootDomain string) []string {
	var filteredLinks []string
	for _, link := range links {
		if strings.HasPrefix(link, rootDomain) {
			filteredLinks = append(filteredLinks, link)
			log.Info("add  : ", link)
		} else {
			log.Info("skip : ", link)
		}
	}
	return filteredLinks
}

// Crawl crawls a webpage and cretes sitemap
func (cm *CrawlManager) Crawl(rootURL string) (map[string]sitemap.Children, error) {
	stmp := map[string]sitemap.Children{}
	i := 0
	urls := []string{rootURL}
	linksPerPage := viper.GetInt("LINKS_PER_PAGE")
	pageLimit := viper.GetInt("PAGE_LIMIT")
	for len(urls) > 0 {
		url := urls[0]
		links, err := cm.fetcher.ExtractURLs(url)
		if err != nil {
			if err != nil {
				if err != crawlers.ErrPageNotHTML {
					return nil, fmt.Errorf("crawl manager: %s", err)
				}
				log.Error("crawler : ", err)
			}
		}

		children := filterDomains(links, rootURL)

		k := 0
		for _, link := range children {
			if _, ok := stmp[link]; !ok {
				stmp[url] = append(stmp[url], link)
				stmp[link] = sitemap.Children{}
				urls = append(urls, link)
				k++
			}
			if linksPerPage > 0 && k >= linksPerPage {
				break
			}
		}

		urls = urls[1:]
		i++
		log.Info("links processed : ", i, " : links in queue : ", len(urls))
		if pageLimit != 0 && i >= pageLimit {
			break
		}
	}
	return stmp, nil
}
