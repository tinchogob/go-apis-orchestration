package main

import (
	"fmt"
)

import "time"

type Search struct {
	Items []SearchItem
}

type SearchItem struct {
	ID int
}

type Site struct {
	Name string
}

type Item struct {
	ID   int
	Name string
}

func GetSearch() (*Search, error) {
	<-time.After(time.Millisecond * 1000)

	s := &Search{
		[]SearchItem{{1}, {2}, {3}},
	}

	return s, nil
}

func GetSite() (*Site, error) {
	<-time.After(time.Millisecond * 100)

	return &Site{"site"}, nil
}

func GetItem(id int) (*Item, error) {
	<-time.After(time.Millisecond * 100)
	return &Item{
		ID:   id,
		Name: fmt.Sprintf("Item-%v", id),
	}, nil
}
