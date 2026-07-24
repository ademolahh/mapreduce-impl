package main

import (
	"fmt"
	"os"
)

func main() {
	dir, _ := os.ReadDir("words/text")
	for _, o := range dir {
		fmt.Println(o.Name())
	}

}
