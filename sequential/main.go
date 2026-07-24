package main

import (
	"fmt"
	"os"
	"plugin"
	"sort"

	"github.com/ademolahh/map-reduce-impl/shared"
)

type KeyValue []shared.KV

func (kv KeyValue) Len() int {
	return len(kv)
}

func (kv KeyValue) Less(i, j int) bool {
	return kv[i].Word < kv[j].Word
}

func (kv KeyValue) Swap(i, j int) {
	kv[i], kv[j] = kv[j], kv[i]
}

func main() {
	mps, rdc := fetch("a.so")
	file := "words/text/a.txt"
	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	content := string(data)

	var k KeyValue

	k = mps(file, content)

	sort.Sort(k)

	mr, err := os.Create("words/result/result.txt")
	if err != nil {
		panic(err)
	}
	defer mr.Close()

	var values []string
	for i := 0; i < len(k)-1; i++ {
		if k[i].Word == k[i+1].Word {
			values = append(values, k[i].Value)
		} else {
			if len(values) == 0 {
				values = append(values, k[i].Value)
			}
			result := rdc(k[i].Word, values)
			fmt.Fprintf(mr, "%s %s\n", k[i].Word, result)
			values = []string{}
			values = append(values, "1")
		}

	}

}

func fetch(program string) (func(string, string) []shared.KV, func(string, []string) string) {
	pg, err := plugin.Open(program)
	if err != nil {
		panic(err)
	}

	mps, err := pg.Lookup("Map")
	if err != nil {
		panic(err)
	}

	mpFunc := mps.(func(string, string) []shared.KV)

	rdc, err := pg.Lookup("Reduce")
	if err != nil {
		panic(err)
	}

	rdcFunc := rdc.(func(string, []string) string)

	return mpFunc, rdcFunc

}
