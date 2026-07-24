package main

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/ademolahh/map-reduce-impl/shared"
)

func Map(fileName, content string) []shared.KV {
	var result []shared.KV

	for w := range strings.FieldsSeq(content) {
		w = strings.ToLower(w)
		w = strings.TrimFunc(w, unicode.IsPunct)
		res := shared.KV{Word: w, Value: "1"}
		result = append(result, res)
	}

	return result
}

func Reduce(word string, values []string) string {
	return strconv.Itoa(len(values))
}
