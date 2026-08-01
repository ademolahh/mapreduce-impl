package main

import (
	"flag"
	"os"

	"github.com/ademolahh/map-reduce-impl/internal/distributed/master"
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

	nr := flag.Int("nr", 0, "number of reduce task")
	flag.Parse()

	ms := master.New(files, *nr)

	if err := master.Serve(ms); err != nil {
		panic(err)
	}

}
