package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
)

const maxNestingLevel = 3

type Link struct {
	url    string
	status bool
}

func doResearch(url string, out chan<- Link, nestingLevel int) {
	resp, err := http.Get(url)

	isInvalid := err != nil
	out <- Link{url, !isInvalid}
	if isInvalid {
		return
	}

	if nestingLevel < maxNestingLevel {
		wg := sync.WaitGroup{}
		for url = range findPageUrls(resp) {
			wg.Add(1)
			go func() {
				doResearch(url, out, nestingLevel+1)
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func research(url string, out chan<- Link) {
	doResearch(url, out, 0)
}

func main() {
	foundUrls := make(map[string]Link)
	seedUrls := os.Args[1:]

	links := make(chan Link)
	chFinished := make(chan bool)

	for _, url := range seedUrls {
		go func() {
			research(url, links)
			chFinished <- true
		}()
	}

	for c := 0; c < len(seedUrls); {
		select {
		case link := <-links:
			foundUrls[link.url] = link
		case <-chFinished:
			c++
		}
	}

	fmt.Println("\nFound", len(foundUrls), "unique urls:\n")

	for url, link := range foundUrls {
		fmt.Printf(" - [%s] - %s \n", link.status, url)
	}

	close(links)
}
