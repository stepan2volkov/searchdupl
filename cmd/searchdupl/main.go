package main

import (
	"fmt"
	"github.com/stepan2volkov/searchdupl/search"
)

func main()  {
	results := search.Scan("/Users/stepan/go/pkg/mod", false)
	for duplicate := range results {
		fmt.Printf("Found duplicate: \"%s\"\n", duplicate)
	}
}
