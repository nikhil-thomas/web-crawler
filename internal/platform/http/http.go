package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/nikhil-thomas/web-crawler/internal/crawlers"
	"golang.org/x/net/html"
)

// Fetcher implements crawlers.URLFetcher interface
type Fetcher struct{}

// NewFetcher creates and returns a Fetcher
func NewFetcher() *Fetcher {
	return &Fetcher{}
}

// ExtractURLs returns all the links from a page
// only links from anchor tags(<a href="url"></a>) are returned
func (f *Fetcher) ExtractURLs(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http fetcher: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http fetcher: http.Get status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	if !isHTML(resp) {
		return nil, crawlers.ErrPageNotHTML
	}

	rootNode, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("http fetcher: %s", err)
	}

	var links []string
	rawLinks := walkDOM(rootNode, parseHTMLAnchorTag)

	for _, link := range rawLinks {
		absoluteLink, err := resp.Request.URL.Parse(link)
		if err != nil {
			continue
		}
		links = append(links, absoluteLink.String())
	}

	return links, nil
}

func walkDOM(n *html.Node, fn func(n *html.Node) (string, bool)) []string {
	var links []string
	if fn != nil {
		link, ok := fn(n)
		if ok {
			links = append(links, link)
		}
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		childLinks := walkDOM(child, fn)
		links = append(links, childLinks...)
	}
	return links
}

func parseHTMLAnchorTag(node *html.Node) (string, bool) {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				return attr.Val, true
			}
		}
	}
	return "", false
}

func isHTML(resp *http.Response) bool {
	ct := resp.Header.Get("Content-Type")
	if ct != "text/html" && !strings.HasPrefix(ct, "text/html;") {
		return false
	}
	return true
}
