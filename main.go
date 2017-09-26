package main

import (
	"fmt"
)

func main() {
	//s, e := Example1()
	//s, e := Example2()
	s, e := Example3()
	if e != nil {
		panic(e)
	}

	fmt.Println(s)
}
