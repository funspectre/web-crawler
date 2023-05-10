package main

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

func main() {
	group := new(errgroup.Group)

	pendingLinks := make(chan string)
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

	go func() {
		for {
			select {
			case link := <-pendingLinks:
				if visitedLinks.Visited(link) {
					log.Println(fmt.Sprintf("%s has been visited already", link))
					continue
				}

				if visitedLinks.Scraping(link) {
					log.Println(fmt.Sprintf("%s is being visited currently", link))
					continue
				}

				group.Go(func() error {
					return Scrape(parsedStartingUrl, &visitedLinks, &siteMap, pendingLinks, link)
				})
			}
		}
	}()

	pendingLinks <- parsedStartingUrl.String()

	if err := group.Wait(); err != nil {
		log.Println("Error Found!")
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
