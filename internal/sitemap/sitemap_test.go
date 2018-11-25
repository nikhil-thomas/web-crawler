package sitemap_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/nikhil-thomas/web-crawler/internal/sitemap"
)

type stubCrawler struct{}

func (sc *stubCrawler) Crawl(url string) (map[string]sitemap.Children, error) {
	stmp := map[string]sitemap.Children{
		"https://example.com": sitemap.Children{
			"https://example.com/about.html",
			"https://example.com/contact.html",
		},
		"https://example.com/about.html": sitemap.Children{
			"https://example.com/about/rev1.html",
			"https://example.com/about/rev2.html",
		},
		"https://example.com/contact.html": sitemap.Children{
			"https://example.com/contact/rev1.html",
			"https://example.com/contact/rev2.html",
		},
	}
	return stmp, nil
}

func TestSiteMapManager(t *testing.T) {
	crawler := &stubCrawler{}
	t.Run("it shoudl create an SiteMapManager", func(t *testing.T) {
		stmpMng := sitemap.NewSiteManager(
			"https://example.com",
			crawler,
		)

		if stmpMng == nil {
			t.Error("expeced non-nil SiteMapManger instance")
		}
	})

	t.Run("it should return sitemap", func(t *testing.T) {
		stmpMng := sitemap.NewSiteManager(
			"https://example.com",
			crawler,
		)
		stmpMng.Crawl()

		stmp := stmpMng.Sitemap

		gotBytes, _ := json.Marshal(stmp)
		got := string(gotBytes)

		expected := `{"https://example.com":["https://example.com/about.html","https://example.com/contact.html"],"https://example.com/about.html":["https://example.com/about/rev1.html","https://example.com/about/rev2.html"],"https://example.com/contact.html":["https://example.com/contact/rev1.html","https://example.com/contact/rev2.html"]}`

		if expected != got {
			t.Errorf("expected %s, got %s", expected, got)
		}
	})

	t.Run("it should print sitemap", func(t *testing.T) {
		stmpMng := sitemap.NewSiteManager(
			"https://example.com",
			crawler,
		)
		stmpMng.Crawl()

		got := &bytes.Buffer{}

		stmpMng.FPrintMap(got)

		expected := `
::::: Site Map: https://example.com ::::

https://example.com
  https://example.com/about.html
    https://example.com/about/rev1.html
    https://example.com/about/rev2.html
  https://example.com/contact.html
    https://example.com/contact/rev1.html
    https://example.com/contact/rev2.html
`

		if expected != got.String() {
			t.Errorf("expected %s, got %s", expected, got)
		}
	})
}
