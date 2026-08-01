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

	ms := master.New(files, 5)

	if err := master.Serve(ms); err != nil {
		panic(err)
	}

}
