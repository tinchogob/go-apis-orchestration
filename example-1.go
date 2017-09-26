package main

import (
	"bytes"
	"fmt"
	"sync"
)

func Example1() (string, error) {

	rm := NewResourceManager()

	se, e := rm.GetSearch()
	if e != nil {
		return "", e
	}

	si, e := rm.GetSite()
	if e != nil {
		return "", e
	}

	items, itemErrs := rm.GetItems()
	for _, e := range itemErrs {
		if e != nil {
			return "", e
		}
	}

	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Search over site %v had %v items: ", len(se.Items), si.Name))

	for _, it := range items {
		buffer.WriteString(it.Name + " ")
	}

	return buffer.String(), nil
}

type resourceManager struct {
	search    *Search
	searchErr error
	searchWg  *sync.WaitGroup
	site      *Site
	siteErr   error
	siteWg    *sync.WaitGroup
	items     []*Item
	itemsErr  []error
	itemsWg   *sync.WaitGroup
}

func NewResourceManager() *resourceManager {
	rm := &resourceManager{
		searchWg: new(sync.WaitGroup),
		siteWg:   new(sync.WaitGroup),
		itemsWg:  new(sync.WaitGroup),
	}

	rm.searchWg.Add(1)
	go rm.getSearch()

	rm.siteWg.Add(1)
	go rm.getSite()

	rm.itemsWg.Add(1)
	go rm.getItems()

	return rm
}

func (rm *resourceManager) getSearch() {
	rm.search, rm.searchErr = GetSearch()
	rm.searchWg.Done()
}

func (rm *resourceManager) getSite() {
	rm.site, rm.siteErr = GetSite()
	rm.siteWg.Done()
}

func (rm *resourceManager) getItems() {
	s, e := rm.GetSearch()

	if e != nil {
		return
	}

	rm.items = make([]*Item, len(s.Items))
	rm.itemsErr = make([]error, len(s.Items))

	anotherItemsWg := new(sync.WaitGroup)
	anotherItemsWg.Add(len(s.Items))

	for idx, i := range s.Items {
		go func(index, id int) {
			item, err := GetItem(id)
			rm.items[index] = item
			rm.itemsErr[index] = err
			anotherItemsWg.Done()
		}(idx, i.ID)
	}

	anotherItemsWg.Wait()
	rm.itemsWg.Done()
}

func (rm *resourceManager) GetSearch() (*Search, error) {
	rm.searchWg.Wait()
	return rm.search, rm.searchErr
}

func (rm *resourceManager) GetSite() (*Site, error) {
	rm.siteWg.Wait()
	return rm.site, rm.siteErr
}

func (rm *resourceManager) GetItems() ([]*Item, []error) {
	rm.itemsWg.Wait()
	return rm.items, rm.itemsErr
}
