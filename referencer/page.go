package main

import (
	"net/http"
	"net/url"
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

	u, ok := getHref(t)
	if !ok {
		return "", false
	}

	return u, true
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
				u, ok := findUrl(z.Token())
				if ok {
					urls <- u
				}
			}
		}
		close(urls)
	}()

	return urls
}

type Preparer func(u string) (*url.URL, bool)

func getUniquePageUrls(res *http.Response, p Preparer) map[string]*url.URL {
	urls := make(map[string]*url.URL)
	for u := range findPageUrls(res) {
		prepared, ok := p(u)
		if ok {
			urls[prepared.Path] = prepared
		}
	}

	return urls
}
