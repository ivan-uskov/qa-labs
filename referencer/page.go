package main

import (
	"strings"
	"net/http"
	"golang.org/x/net/html"
)

func getHref(t html.Token) (string, bool) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			return a.Val, true
		}
	}

	return "", false
}

func findUrl(t html.Token) (string, bool) {
	isAnchor := t.Data == "a"
	if !isAnchor {
		return "", false
	}

	url, ok := getHref(t)
	if !ok {
		return "", false
	}

	httpLink := strings.Index(url, "http") == 0
	if !httpLink {
		return "", false
	}

	return url, true
}

func findPageUrls(resp *http.Response) <-chan string {
	b := resp.Body

	urls := make(chan string)
	go func() {
		defer b.Close()
		z := html.NewTokenizer(b)

		for {
			tt := z.Next()
			if tt == html.ErrorToken {
				break
			}
			if tt == html.StartTagToken {
				url, ok := findUrl(z.Token())
				if ok {
					urls <- url
				}
			}
		}
		close(urls)
	}()

	return urls
}