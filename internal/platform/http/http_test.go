package http_test

import (
	"encoding/json"
	nethttp "net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/nikhil-thomas/web-crawler/internal/crawlers"
	"github.com/nikhil-thomas/web-crawler/internal/platform/http"
)

func htmlPageHandler(w nethttp.ResponseWriter, r *nethttp.Request) {
	pageContent := `<!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <meta http-equiv="X-UA-Compatible" content="ie=edge">
        <title>Document</title>
    </head>
    <body>
        <a href="https://www.example.com"></a>
        <p>Fugiat ullamco occaecat eiusmod veniam pariatur ipsum esse do do
            excepteur ad dolore. Irure laborum exercitation sunt amet nulla
            reprehenderit. Cupidatat voluptate nostrud voluptate Lorem
            aliqua aliquip velit amet amet.
        </p>
        <a href="https://github.com"></a>
        <a href="https://golang.org"></a>
    </body>
    </html>`
	w.Write([]byte(pageContent))
}

func jsonHandler(w nethttp.ResponseWriter, r *nethttp.Request) {
	data := map[string]string{
		"key1": "val1",
		"key2": "val2",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func TestExtractURLs(t *testing.T) {
	r := nethttp.NewServeMux()
	r.Handle("/index.html", nethttp.HandlerFunc(htmlPageHandler))
	r.Handle("/data.json", nethttp.HandlerFunc(jsonHandler))
	server := httptest.NewServer(r)

	t.Run("it should returns urls from html pages", func(t *testing.T) {

		fetcher := http.NewFetcher()
		urls, err := fetcher.ExtractURLs(server.URL + "/index.html")

		if err != nil {
			t.Errorf("error unexpected, got %s", err)
		}

		expected := 3
		got := len(urls)
		if expected != got {
			t.Errorf("expected %d, got %d", expected, got)
		}

		expectedURLs := []string{
			"https://www.example.com",
			"https://github.com",
			"https://golang.org",
		}
		gotURLs := urls

		if !reflect.DeepEqual(expectedURLs, gotURLs) {
			t.Errorf("expected %v, got %v", expectedURLs, gotURLs)
		}
	})

	t.Run("it should returns urls from html pages", func(t *testing.T) {

		fetcher := http.NewFetcher()
		urls, err := fetcher.ExtractURLs(server.URL + "/data.json")

		if urls != nil {
			t.Errorf("expected nil, got %v", urls)
		}

		if err == nil {
			t.Errorf("error expected, but got %s", err)
		}

		expected := crawlers.ErrPageNotHTML

		if err != expected {
			t.Errorf("expected %s, but got %s", expected, err)
		}
	})
}
