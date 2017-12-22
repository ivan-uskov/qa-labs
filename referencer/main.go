package main

import (
	"fmt"
	"net/http"
	"os"
	"net/url"
	"sync"
	"time"
	"strings"
)

const maxNestingLevel = 20

type Link struct {
	url    *url.URL
	status bool
}

func (l *Link) Status() string {
	if l.status {
		return "Good"
	} else {
		return "Bad"
	}
}

type Researcher struct {
	u       *url.URL
	visited *sync.Map
}

func newResearcher(u string) *Researcher {
	ur, err := url.Parse(u)
	if err != nil || ur.Host == "" || ur.Scheme == "" {
		return nil
	}

	return &Researcher{ur, &sync.Map{}}
}

func (r *Researcher) prepareUrl(urlStr string) (*url.URL, bool) {
	u, err := url.Parse(urlStr)
	if err != nil || (u.Host+u.Path) == "" {
		return nil, false
	}
	if u.Host == "" {
		u.Host = r.u.Host
	}
	if u.Scheme == "" {
		u.Scheme = r.u.Scheme
	}
	if (u.Host != r.u.Host) && ((`www.`+u.Host) != r.u.Host) && (u.Host != (`www.`+r.u.Host)) {
		return nil, false
	}

	u2, e := url.Parse(u.String())
	if e != nil {
		return nil, false
	}

	return u2, true
}

func (r *Researcher) doResearch(u *url.URL, out chan<- Link, nestingLevel int) {
	_, exists := r.visited.LoadOrStore(u.Path, nil)
	if exists {
		return
	}

	resp, err := http.Get(u.String())
	isInvalid := err != nil || resp.StatusCode != 200
	out <- Link{u, !isInvalid}
	if isInvalid  || !strings.Contains(resp.Header.Get("Content-Type"), `text/html`) {
		return
	}

	if nestingLevel < maxNestingLevel {
		wg := sync.WaitGroup{}
		for _, u := range getUniquePageUrls(resp, r.prepareUrl) {
			wg.Add(1)
			go func() {
				r.doResearch(u, out, nestingLevel+1)
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func (r *Researcher) research() <-chan Link {
	out := make(chan Link)
	go func() {
		r.doResearch(r.u, out, 0)
		close(out)
	}()

	return out
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println(`Expected url argument`)
		os.Exit(1)
	}

	start := time.Now()
	r := newResearcher(args[0])
	if r == nil {
		fmt.Println(`Invalid url specified, expected <scheme>://<host>, got: ` + args[0])
		os.Exit(1)
	}

	for link := range r.research() {
		fmt.Printf(" - [%s] - %s \n", link.Status(), link.url.Path)
	}

	fmt.Printf("Execution time %s", time.Since(start).String())
}
