package main

import (
	"bytes"
	"errors"
	"fmt"
)

func Example3() (string, error) {

	n1_calls := 2
	c1 := make(chan response, n1_calls)

	getSearch3(c1)
	getSite3(c1)

	var search *Search
	var site *Site

	for n := 0; n < n1_calls; n++ {

		resp := <-c1

		if resp.e != nil {
			return "", resp.e
		}

		switch v := resp.r.(type) {
		case *Search:
			search = resp.r.(*Search)
		case *Site:
			site = resp.r.(*Site)
		default:
			return "", errors.New(fmt.Sprintf("error casting: %v", v))
		}
	}

	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Search over site %v had %v items: ", site.Name, len(search.Items)))

	n2_calls := len(search.Items)
	c2 := make(chan response, n2_calls)

	getItems3(search.Items, c2)
	for n := 0; n < n2_calls; n++ {
		resp := <-c2

		if resp.e != nil {
			return "", resp.e
		}

		if item, ok := resp.r.(*Item); ok {
			buffer.WriteString(item.Name + " ")
		} else {
			return "", errors.New("cast to item failed")
		}

	}

	return buffer.String(), nil
}

type response struct {
	r interface{}
	e error
}

func getSearch3(c chan response) {
	go func() {
		s, e := GetSearch()
		c <- response{s, e}
	}()
}

func getSite3(c chan response) {
	go func() {
		s, e := GetSite()
		c <- response{s, e}
	}()
}

func getItems3(items []SearchItem, c chan response) {
	for _, item := range items {
		go func(id int) {
			i, e := GetItem(id)
			c <- response{i, e}
		}(item.ID)
	}
}
