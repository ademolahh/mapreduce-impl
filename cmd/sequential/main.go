package main

import (
	"fmt"
	"os"

	"github.com/ademolahh/map-reduce-impl/internal/sequential"
)

func main() {
	if err := sequential.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
