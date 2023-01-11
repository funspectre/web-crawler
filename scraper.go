package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

func normalizeUrl(rootUrl *url.URL, link string) *url.URL {
	u, err := url.Parse(strings.TrimSpace(link))
	if err != nil {
		log.Println(err)
		return nil
	}

	u = rootUrl.ResolveReference(u)

	if u.Path == "" {
		u.Path = "/"
	}

	u.Fragment = ""

	return u
}

func end(wg *sync.WaitGroup, pendingJobCount chan<- int) {
	wg.Done()
	pendingJobCount <- -1
	return
}

func Scrape(wg *sync.WaitGroup, rootUrl *url.URL, visitedLinks *LinkHash, siteMap *SiteMap, pendingLinks chan<- string, pendingJobCount chan<- int, link string) {
	log.Printf("%d unique links visited", visitedLinks.Size())

	defer func() {
		if r := recover(); r != nil {
			m := fmt.Sprintf("Recovered from error while processing %s:\n", link)
			log.Println(m, r)
		}
	}()
	defer end(wg, pendingJobCount)

	if visitedLinks.Has(link) {
		log.Panicf("%s has been visited already", link)
	}

	log.Printf("Attempting to visit %s", link)
	visitedLinks.Try(link)
	res, err := http.Get(link)
	if err != nil {
		if visitedLinks.Tries(link) < 3 {
			pendingLinks <- link
		}
		log.Panic(err)
	}
	defer res.Body.Close()

	contentType := res.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		visitedLinks.Add(link)
		log.Panicf("%s is of Content-Type: %s", link, contentType)
	}

	visitedLinks.Add(link)

	if res.StatusCode != 200 {
		log.Panicf("status code error: %d %s", res.StatusCode, res.Status)
	}

	log.Printf("Page fetched successfully %s", link)

	log.Printf("Parsing document at %s", link)
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Document at %s parsed successfully", link)

	linksCount := 0

	links := make([]string, 0)

	seenLinks := make(map[string]bool)

	// get all anchor tags with qoquery
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		// fetch all their hrefs
		link, _ := s.Attr("href")

		u := normalizeUrl(rootUrl, link)

		// filter URLs from a different host
		if !strings.HasPrefix(u.Scheme, "http") || u.Host != rootUrl.Host {
			return
		}

		resolvedLink := u.String()

		if !seenLinks[resolvedLink] {
			seenLinks[resolvedLink] = true
			links = append(links, resolvedLink)
		}

		if !visitedLinks.Has(resolvedLink) {
			pendingLinks <- resolvedLink
		}

		linksCount++
	})

	siteMap.add(link, links)

	log.Printf("%d eligible links in document at %s scraped successfully", linksCount, link)
}
