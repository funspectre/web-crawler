package main

import (
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	pendingLinks := make(chan string, 32)
	pendingJobCount := make(chan int)
	visitedLinks := LinkHash{
		visited: map[string]bool{},
		trials:  map[string]int{},
	}
	siteMap := SiteMap{
		table: map[string][]string{},
	}

	// get root URL from CLI argument
	startingUrl := os.Args[1]

	// normalize root URL
	parsedStartingUrl, err := url.Parse(startingUrl)
	if err != nil {
		log.Fatal(err)
	}
	if parsedStartingUrl.Path == "" {
		parsedStartingUrl.Path = "/"
	}

	go WatchDog(pendingLinks, pendingJobCount)

	pendingLinks <- parsedStartingUrl.String()

	for link := range pendingLinks {
		if visitedLinks.Has(link) {
			continue
		}
		pendingJobCount <- 1
		wg.Add(1)
		go Scrape(&wg, parsedStartingUrl, &visitedLinks, &siteMap, pendingLinks, pendingJobCount, link)
	}

	wg.Wait()

	// print output summary
	log.Printf("%d unique web pages found in domain: %s\n", visitedLinks.Size(), parsedStartingUrl.Host)

	path, err := filepath.Abs("./sitemap.txt")
	if err != nil {
		log.Panicln(err)
	}
	f, err := os.Create(path)
	if err != nil {
		log.Panicln(err)
	}
	defer f.Close()
	_, err = f.WriteString(siteMap.String())
	if err != nil {
		log.Panicln(err)
	}

	log.Printf("Sitemap successfully written to %s\n", path)
}
