package concurrent

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/nikhil-thomas/web-crawler/internal/crawlers"
	"github.com/nikhil-thomas/web-crawler/internal/sitemap"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// CrawlManager implements sitemap.Crawler interface
type CrawlManager struct {
	fetcher    crawlers.URLFetcher
	done       chan bool
	cache      []string
	supplyChan chan string
	mu         sync.Mutex
}

// Page defines a HTML page and links inside the page
type Page struct {
	url      string
	children []string
}

// NewCrawlManager creates and returns a CrawlManager
func NewCrawlManager(fetcher crawlers.URLFetcher) *CrawlManager {
	queueLength := viper.GetInt("CRAWLER_QUEUE_LENGTH")
	if queueLength == 0 {
		queueLength = 2
	}
	return &CrawlManager{
		fetcher:    fetcher,
		done:       make(chan bool),
		cache:      []string{},
		supplyChan: make(chan string, queueLength),
	}
}

func filterDomains(links []string, rootDomain string) []string {
	var filteredLinks []string
	for _, link := range links {
		if link == "" {
			continue
		}
		if strings.HasPrefix(link, rootDomain) {
			filteredLinks = append(filteredLinks, link)
		} else {
			log.Info("skip   : ", link)
		}
	}
	return filteredLinks
}

// Crawl crawls a webpage and cretes sitemap
func (cm *CrawlManager) Crawl(rootURL string) (map[string]sitemap.Children, error) {
	stmp := map[string]sitemap.Children{}

	PageChan := cm.enqueue()

	linksChan := cm.launchWorkers(PageChan, rootURL)

	sitemapChan := cm.makeSiteMap(linksChan, stmp)

	// pass first input to pipeline
	cm.addToQueue(rootURL)

	// wait for final sitemmap map[string][]string
	stmpOut := <-sitemapChan

	return stmpOut, nil
}

func (cm *CrawlManager) enqueue() chan Page {
	outChan := make(chan Page)
	go func() {
	forLoop:
		for {
			select {
			case <-cm.done:
				break forLoop
			case link := <-cm.supplyChan:
				page := Page{
					url:      link,
					children: nil,
				}
				outChan <- page
			}
		}
		close(outChan)
	}()
	return outChan
}

func (cm *CrawlManager) addToQueue(url string) {
	// add new links to input []string slice
	// pass links to crawl pipeline only if input channel is not full
	cm.mu.Lock()
	cm.cache = append(cm.cache, url)
	cm.mu.Unlock()

	queueLength := viper.GetInt("CRAWLER_QUEUE_LENGTH")
	if queueLength == 0 {
		queueLength = 2
	}
	if len(cm.supplyChan) < queueLength {
		cm.supplyChan <- cm.cache[0]
		cm.mu.Lock()
		cm.cache = cm.cache[1:]
		cm.mu.Unlock()
	}
}

func (cm *CrawlManager) launchWorkers(inChan chan Page, rootURL string) chan Page {

	numWorkers := viper.GetInt("WORKER_COUNT")
	if numWorkers == 0 {
		numWorkers = 10
	}
	outChanList := []chan Page{}

	// Fan Out
	for i := 0; i < numWorkers; i++ {
		outChan := cm.extractWorker(inChan, i+1, rootURL)
		outChanList = append(outChanList, outChan)
	}

	// Fan In
	outChan := cm.merge(outChanList...)
	return outChan
}

// func extractWorker(inChan chan Page, outChan chan Page, fetcher crawlers.URLFetcher, rootURL string) chan Page {
func (cm *CrawlManager) extractWorker(inChan chan Page, id int, rootURL string) chan Page {
	outChan := make(chan Page)
	go func() {
	forLoop:
		for {
			select {
			case <-cm.done:
				break forLoop
			case page := <-inChan:
				log.Debug("worker : ", id, " : url : ", page.url)
				links, err := cm.fetcher.ExtractURLs(page.url)
				if err != nil {
					log.Error("crawl : ", err, page.url)
				}
				page.children = filterDomains(links, rootURL)

				outChan <- page
			}
		}
		log.Debug("worker : exited : ", id)
		close(outChan)
	}()
	return outChan
}

func (cm *CrawlManager) makeSiteMap(inChan chan Page, stmp map[string]sitemap.Children) chan map[string]sitemap.Children {
	outSiteMapChan := make(chan map[string]sitemap.Children)
	linksPerPage := viper.GetInt("LINKS_PER_PAGE")
	pageLimit := viper.GetInt("PAGE_LIMIT")

	// queueEmptyTimeoutEvent trigger crawl stop
	// when queue is empty for more than specified time duration
	queueEmptyTimeoutEvent := make(chan bool)

	// queueCloseCtx, queueCloseCancel are used
	// to make sitemap generation safe to stop
	// without causing duplicate channel closes
	var queueCloseCtx context.Context
	queueCloseCancel := func() {}

	go func() {
		i := 0
	forLoop:
		for {
			select {
			case <-cm.done:
				break forLoop
			case <-queueEmptyTimeoutEvent:
				break forLoop
			case page := <-inChan:
				k := 0
				for _, link := range page.children {
					// save link only if it is new
					if _, ok := stmp[link]; !ok {
						// append link to parents children slice
						stmp[page.url] = append(stmp[page.url], link)
						log.Info("add    : ", link)
						// record link in sitemap for further crawling
						stmp[link] = sitemap.Children{}

						// push link to input queue
						cm.addToQueue(link)

						// if there is an active timeout
						// because of empty queue cance it
						// as queue is not empty anymore
						queueCloseCancel()

						k++
						// process only specified number of links perpage
						// if the linksPerPage == 0 process all links from the page
						if linksPerPage > 0 && k >= linksPerPage {
							break
						}
					}
				}
				i++
				// print number of pages processed and
				// number of links currently in the input queue
				log.Info("links  : ", i, " : queue : ", len(cm.supplyChan))

				// if specified number pages are processed stop crawlling
				// if pageLimit param is 0, then thre is no limit
				// crawl till the queue is empty
				if pageLimit != 0 && i >= pageLimit {
					log.Info("crawl  : page limit (", pageLimit, ") reached : stop crawiling")
					break forLoop
				}
				// if input queue is empty
				// stop crawling if after specified timeout
				if len(cm.supplyChan) == 0 {
					queueCloseCancel()
					queueCloseCtx, queueCloseCancel = context.WithCancel(context.Background())
					go endOperationTimeout(queueCloseCtx, queueEmptyTimeoutEvent)
				}
			}
		}
		// issue done signal for all pipeline stages
		cm.StopCrawl()
		outSiteMapChan <- stmp
	}()
	return outSiteMapChan
}

func endOperationTimeout(ctx context.Context, queueEmptyChan chan bool) {
	timeout := viper.GetDuration("CRAWLER_TIMEOUT")
	if timeout == time.Duration(0) {
		timeout = 1 * time.Second
	}
	log.Info("queue  : empty : start crawiling stop timeout : ", timeout)
	select {
	case <-time.After(timeout):
		log.Info("queue  : empty : stop crawiling")
		close(queueEmptyChan)

	// abort timeout if queue is not empty any more
	case <-ctx.Done():
		log.Info("queue  : empty : start crawiling stop timeout cancelled")
	}
}

// StopCrawl stops crawling
func (cm *CrawlManager) StopCrawl() {
	close(cm.done)
}

func (cm *CrawlManager) merge(inChans ...chan Page) chan Page {
	outChan := make(chan Page)

	var wg sync.WaitGroup
	wg.Add(len(inChans))

	for _, inChan := range inChans {
		go func(inChan chan Page) {
			defer wg.Done()
			for page := range inChan {
				outChan <- page
			}
		}(inChan)
	}

	go func() {
		wg.Wait()
		close(outChan)
	}()

	return outChan
}
