package main

import "github.com/ademolahh/map-reduce-impl/distributed/worker"

func main() {
	worker, err := worker.New()
	if err != nil {
		panic(err)
	}

	if err := worker.Work(); err != nil {
		panic(err)
	}
}
