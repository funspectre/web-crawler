package main

import "sync"

type LinkHash struct {
	mutex   sync.Mutex
	visited map[string]bool
	trials  map[string]int
}

func (lh *LinkHash) Add(link string) {
	lh.mutex.Lock()
	defer lh.mutex.Unlock()
	lh.visited[link] = true
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
