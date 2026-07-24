package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
)

const READ_FOLDER string = "words/text"
const OUT_PATH = "words/result/result.txt"

func main() {
	build := buildFlag()

	mps, rdc, err := fetch(build)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	log.Println("plugin loaded")

	kvs, err := runMapPhase(mps, READ_FOLDER)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	log.Println("map phase done")

	sort.Sort(kvs)

	if err := runReducePhase(kvs, rdc, OUT_PATH); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	log.Println("reduce phase done")
}

func buildFlag() string {
	var build = flag.String("build", "", "")
	flag.Parse()

	if *build == "" {
		panic("expect: go run ./sequential <build>")
	}

	return *build
}
