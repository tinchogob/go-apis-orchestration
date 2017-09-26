package main

import (
	"bytes"
	"fmt"
)

func Example2() (string, error) {
	searchCall := getSearch()
	siteCall := getSite()

	searchResp := <-searchCall
	if searchResp.err != nil {
		return "", searchResp.err
	}

	itemsCalls := getItems(searchResp.search.Items)

	siteResp := <-siteCall
	if siteResp.err != nil {
		return "", siteResp.err
	}

	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Search over site %v had %v items: ", len(searchResp.search.Items), siteResp.site.Name))

	for _, itemCall := range itemsCalls {
		itemResp := <-itemCall
		if itemResp.err != nil {
			return "", itemResp.err
		}

		buffer.WriteString(itemResp.item.Name + " ")
	}

	return buffer.String(), nil
}

type searchResponse struct {
	search *Search
	err    error
}

func getSearch() chan searchResponse {
	c := make(chan searchResponse, 1)
	go func() {
		s, e := GetSearch()
		c <- searchResponse{s, e}
	}()
	return c
}

type siteResponse struct {
	site *Site
	err  error
}

func getSite() chan siteResponse {
	c := make(chan siteResponse, 1)
	go func() {
		s, e := GetSite()
		c <- siteResponse{s, e}
	}()
	return c
}

type itemResponse struct {
	item *Item
	err  error
}

func getItem(id int) chan itemResponse {
	c := make(chan itemResponse, 1)
	go func() {
		i, e := GetItem(id)
		c <- itemResponse{i, e}
	}()
	return c
}

func getItems(items []SearchItem) []chan itemResponse {

	var calls []chan itemResponse
	for _, item := range items {
		calls = append(calls, getItem(item.ID))
	}

	return calls
}
