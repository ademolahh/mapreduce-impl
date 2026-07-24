package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
)

func main() {
	build := buildFlag()

	mps, rdc, err := fetch(build)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	const FOLDER string = "words/text"
	kvs, err := runMapPhase(mps, FOLDER)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	sort.Sort(kvs)

	if err := runReducePhase(kvs, rdc, "words/result/result.txt"); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func buildFlag() string {
	var build = flag.String("build", "", "")
	flag.Parse()

	if *build == "" {
		panic("expect: go run ./sequential <build>")
	}

	return *build
}
