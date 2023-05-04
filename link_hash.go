package main

import "sync"

type LinkHash struct {
	mutex    sync.Mutex
	scraping map[string]bool
	visited  map[string]bool
	trials   map[string]int
}

func (lh *LinkHash) Add(link string) {
	lh.mutex.Lock()
	defer lh.mutex.Unlock()
	lh.visited[link] = true
}
func (lh *LinkHash) IsScraping(link string) bool {
	lh.mutex.Lock()
	defer lh.mutex.Unlock()
	return lh.scraping[link] == true
}

func (lh *LinkHash) Has(link string) bool {
	lh.mutex.Lock()
	defer lh.mutex.Unlock()
	return lh.visited[link] == true
}
func (lh *LinkHash) Try(link string) {
	lh.mutex.Lock()
	defer lh.mutex.Unlock()
	lh.trials[link]++
	lh.scraping[link] = true
}

func (lh *LinkHash) Failed(link string) {
	lh.mutex.Lock()
	defer lh.mutex.Unlock()
	lh.scraping[link] = false
}

func (lh *LinkHash) Tries(link string) int {
	lh.mutex.Lock()
	defer lh.mutex.Unlock()
	return lh.trials[link]
}

func (lh *LinkHash) Size() int {
	lh.mutex.Lock()
	defer lh.mutex.Unlock()
	return len(lh.visited)
}
