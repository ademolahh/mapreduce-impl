package sequential

import (
	"flag"
	"fmt"
	"log"
	"sort"

	"github.com/ademolahh/map-reduce-impl/internal/shared"
)

const READ_FOLDER string = "words/text"
const OUT_PATH = "words/result/result.txt"

func Run() error {
	build := buildFlag()

	mps, rdc, err := shared.Fetch(build)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	log.Println("plugin loaded")

	kvs, err := runMapPhase(mps, READ_FOLDER)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	log.Println("map phase done")

	sort.Sort(kvs)

	if err := runReducePhase(kvs, rdc, OUT_PATH); err != nil {
		return fmt.Errorf("error: %w", err)
	}
	log.Println("reduce phase done")
	return nil
}

func buildFlag() string {
	var build = flag.String("build", "", "")
	flag.Parse()

	if *build == "" {
		panic("expect: go run ./sequential <build>")
	}

	return *build
}
