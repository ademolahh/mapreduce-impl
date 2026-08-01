package main

import (
	"flag"

	"github.com/ademolahh/map-reduce-impl/internal/distributed/worker"
)

func main() {
	worker, err := worker.New()
	if err != nil {
		panic(err)
	}

	programPath := flag.String("program", "", "path to MapReduce application plugin (.so file)")
	flag.Parse()

	if err := worker.Work(*programPath); err != nil {
		panic(err)
	}
}
