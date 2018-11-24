package crawlers

// URLFetcher defines extraction of links from an html page
type URLFetcher interface {
	ExtractURLs(url string) ([]string, error)
}
