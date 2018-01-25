package http

import (
	"net/http"
	"strings"
	"log"
)

func CreateRequest(url string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, strings.Trim(url, "\n"), nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml,application/rss+xml;q=0.9,*/*;q=0.8;q=0.7")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:57.0) Gecko/20100101 Firefox/57.0")

	return req
}

func SendRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	return client.Do(req)
}
