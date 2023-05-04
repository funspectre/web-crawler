package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
	"strings"
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

func Scrape(rootUrl *url.URL, visitedLinks *LinkHash, siteMap *SiteMap, pendingLinks chan<- string, link string) error {
	//log.Printf("%d unique links visited", visitedLinks.Size())

	//log.Printf("Attempting to visit %s", link)
	visitedLinks.Try(link)
	res, err := http.Get(link)
	if err != nil {
		visitedLinks.Failed(link)
		if visitedLinks.Tries(link) < 3 {
			pendingLinks <- link
		}
		return err
	}
	defer res.Body.Close()

	contentType := res.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		visitedLinks.Add(link)

		//log.Println(fmt.Sprintf("%s is of Content-Type: %s", link, contentType))
		return nil
	}

	visitedLinks.Add(link)

	if res.StatusCode != 200 {
		if res.StatusCode == 404 {
			log.Println(fmt.Sprintf("Document at %s not found", link))
			return nil
		}

		return errors.New(fmt.Sprintf("%s status code error: %d %s", link, res.StatusCode, res.Status))
	}

	//log.Printf("Page fetched successfully %s", link)

	//log.Printf("Parsing document at %s", link)
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	//log.Printf("Document at %s parsed successfully", link)

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

	siteMap.Add(link, links)

	//log.Printf("%d eligible links in document at %s scraped successfully", linksCount, link)

	return nil
}
