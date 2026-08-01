package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dirs, _ := os.ReadDir("shared")
	for _, dir := range dirs {
		name := dir.Name()
		name = strings.TrimSuffix(name, filepath.Ext(name))
		fmt.Println(name)

	}
}
