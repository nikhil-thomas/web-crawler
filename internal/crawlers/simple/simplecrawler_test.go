package simple_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/nikhil-thomas/web-crawler/internal/crawlers/simple"
)

type stubURLFetcher struct {
	urls  map[string][]string
	index int
}

func (suf *stubURLFetcher) ExtractURLs(url string) ([]string, error) {
	links := suf.urls[url]
	fmt.Println(links)
	return links, nil
}

func TestCrawlManager(t *testing.T) {
	urlFetcher := &stubURLFetcher{
		urls: map[string][]string{
			"https://example.com": []string{
				"https://example.com/about.html",
				"https://example.com/contact.html",
			},
			"https://example.com/about.html": []string{
				"https://example.com/about/rev1.html",
				"https://example.com/about/rev2.html",
			},
			"https://example.com/contact.html": []string{
				"https://example.com/contact/rev1.html",
				"https://example.com/contact/rev2.html",
			},
		},
	}

	t.Run("it shoudl create an concurrent crawler", func(t *testing.T) {
		crwl := simple.NewCrawlManager(urlFetcher)

		if crwl == nil {
			t.Error("expeced non-nil SiteMapManger instance")
		}
	})
	t.Run("it should generate a sitemap from urls", func(t *testing.T) {
		crwl := simple.NewCrawlManager(urlFetcher)
		stmp, err := crwl.Crawl("https://example.com")

		if err != nil {
			t.Errorf("expected no error, got %s", err)
		}

		gotBytes, _ := json.Marshal(stmp)
		got := string(gotBytes)

		expected := `{"https://example.com":["https://example.com/about.html","https://example.com/contact.html"],"https://example.com/about.html":["https://example.com/about/rev1.html","https://example.com/about/rev2.html"],"https://example.com/about/rev1.html":[],"https://example.com/about/rev2.html":[],"https://example.com/contact.html":["https://example.com/contact/rev1.html","https://example.com/contact/rev2.html"],"https://example.com/contact/rev1.html":[],"https://example.com/contact/rev2.html":[]}`

		if expected != got {
			t.Errorf("expected %s, got %s", expected, got)
		}

	})
}
