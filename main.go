package main

import (
	"golang.org/x/sync/errgroup"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

func main() {
	group := new(errgroup.Group)

	pendingLinks := make(chan string, 32)
	pendingJobCount := make(chan int)
	visitedLinks := LinkHash{
		scraping: map[string]bool{},
		visited:  map[string]bool{},
		trials:   map[string]int{},
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
		group.Go(func() error {
			return Scrape(parsedStartingUrl, &visitedLinks, &siteMap, pendingLinks, pendingJobCount, link)
		})
	}

	//wg.Wait()
	if err := group.Wait(); err != nil {
		log.Println(err)
	}

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
