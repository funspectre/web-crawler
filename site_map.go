package main

import (
	"fmt"
	"sync"
)

type SiteMap struct {
	mutex sync.Mutex
	table map[string][]string
}

func (sm *SiteMap) Add(parentLink string, childLinks []string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.table[parentLink] = childLinks
}

func (sm *SiteMap) Show() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	for parentLink, childLinks := range sm.table {
		fmt.Println(parentLink)
		for _, childLink := range childLinks {
			fmt.Println(childLink)
		}
		fmt.Println("")
	}
}

func (sm *SiteMap) String() string {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	s := ""
	for parentLink, childLinks := range sm.table {
		s += fmt.Sprintln(parentLink)
		for _, childLink := range childLinks {
			s += fmt.Sprintln(childLink)
		}
		s += fmt.Sprintln("")
	}

	return s
}
