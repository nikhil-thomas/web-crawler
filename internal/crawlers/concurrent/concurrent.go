package concurrent

import (
	"strings"
	"time"

	"github.com/nikhil-thomas/web-crawler/internal/crawlers"
	"github.com/nikhil-thomas/web-crawler/internal/sitemap"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// CrawlManager implements sitemap.Crawler interface
type CrawlManager struct {
	fetcher crawlers.URLFetcher
}

// Page defines a HTML page and links inside the page
type Page struct {
	url      string
	children []string
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
			log.Info("add    : ", link)
		} else {
			log.Info("skip   : ", link)
		}
	}
	return filteredLinks
}

// Crawl crawls a webpage and cretes sitemap
func (cm *CrawlManager) Crawl(rootURL string) (map[string]sitemap.Children, error) {
	stmp := map[string]sitemap.Children{}

	done := make(chan bool)

	queueLength := viper.GetInt("CRAWLER_QUEUE_LENGTH")
	urlSUpplyChan := make(chan string, queueLength)

	PageChan := enqueue(done, urlSUpplyChan)

	linksChan := make(chan Page)

	launchWorkers(done, PageChan, linksChan, 10, cm.fetcher, rootURL)

	sitemapChan := makeSiteMap(done, linksChan, urlSUpplyChan, stmp)
	urlSUpplyChan <- rootURL
	stmpOut := <-sitemapChan
	return stmpOut, nil
}

func enqueue(done chan bool, inChan chan string) chan Page {
	outChan := make(chan Page)
	go func() {
	forLoop:
		for {
			select {
			case <-done:
				close(outChan)
				break forLoop
			case link := <-inChan:
				page := Page{
					url:      link,
					children: nil,
				}
				outChan <- page
			}
		}
	}()
	return outChan
}

func launchWorkers(done chan bool, inChan chan Page, outChan chan Page, workerCount int, fetcher crawlers.URLFetcher, rootURL string) {
	// numWorkers := viper.GetInt("WORKER_COUNT")
	// if numWorkers == 0 {
	numWorkers := 10
	// }
	for i := 0; i < numWorkers; i++ {
		extractWorker(done, inChan, outChan, i+1, fetcher, rootURL)
	}
}

// func extractWorker(inChan chan Page, outChan chan Page, fetcher crawlers.URLFetcher, rootURL string) chan Page {
func extractWorker(done chan bool, inChan, outChan chan Page, id int, fetcher crawlers.URLFetcher, rootURL string) {
	go func() {
	forLoop:
		for {
			select {
			case <-done:
				break forLoop
			case page := <-inChan:
				log.Debug("worker : ", id, " : url : ", page.url)
				links, err := fetcher.ExtractURLs(page.url)
				if err != nil {
					if err != nil {
						log.Error("crawl  : ", err, page.url)
					}
				}

				page.children = filterDomains(links, rootURL)

				outChan <- page
			}
		}
		log.Debug("worker : exited : ", id)
	}()
}

func makeSiteMap(done chan bool, inChan chan Page, supplyChan chan string, stmp map[string]sitemap.Children) chan map[string]sitemap.Children {
	outSiteMap := make(chan map[string]sitemap.Children)
	linksPerPage := viper.GetInt("LINKS_PER_PAGE")
	pageLimit := viper.GetInt("PAGE_LIMIT")

	go func() {
		i := 0
	forLoop:
		for {
			select {
			case <-done:
				break forLoop
			case page := <-inChan:
				k := 0
				for _, link := range page.children {
					if _, ok := stmp[link]; !ok {
						stmp[page.url] = append(stmp[page.url], link)
						stmp[link] = sitemap.Children{}
						go func(l string) {
							supplyChan <- l
						}(link)
						//urls = append(urls, link)
						k++
					}
					if linksPerPage > 0 && k >= linksPerPage {
						break
					}
				}
				log.Info("links  : ", i, " : queue : ", len(supplyChan))
				if len(supplyChan) == 0 {
					go endOperationTimeout(done, supplyChan)
				}
				i++
				if pageLimit != 0 && i >= pageLimit {
					close(done)
					break forLoop
				}
			}
		}
		outSiteMap <- stmp
	}()
	return outSiteMap
}

func endOperationTimeout(done chan bool, checkChan chan string) {
	timeout := viper.GetDuration("CRAWLER_TIMEOUT")
	log.Info("queue  : empty : start crawiling stop timeout : ", timeout)
	time.Sleep(timeout)
	if len(checkChan) == 0 {
		log.Info("queue  : empty : stop crawiling")
		close(done)
	}
}
