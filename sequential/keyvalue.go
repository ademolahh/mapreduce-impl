package main

import "github.com/ademolahh/map-reduce-impl/shared"

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
