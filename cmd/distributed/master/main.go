package main

import (
	"os"

	"github.com/ademolahh/map-reduce-impl/distributed/master"
)

func main() {
	dirs, err := os.ReadDir("words/text")
	if err != nil {
		panic(err)
	}

	var files []string
	for _, dir := range dirs {
		files = append(files, dir.Name())
	}

	if err := master.Serve(files, 3); err != nil {
		panic(err)
	}
}
