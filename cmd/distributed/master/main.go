package main

import (
	"os"

	"github.com/ademolahh/map-reduce-impl/distributed/master"
)

func main() {
	folder := "words/text"
	dirs, err := os.ReadDir(folder)
	if err != nil {
		panic(err)
	}

	var files []string
	for _, dir := range dirs {
		files = append(files, "words/text"+"/"+dir.Name())
	}

	if err := master.Serve(files, 3); err != nil {
		panic(err)
	}
}
