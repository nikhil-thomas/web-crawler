package crawlers

import "errors"

// ErrPageNotHTML is returned when the fetched page is not HTML
var ErrPageNotHTML = errors.New("page is not html")

// URLFetcher defines extraction of links from an html page
type URLFetcher interface {
	ExtractURLs(url string) ([]string, error)
}
