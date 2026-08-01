package shared

import (
	"time"
)

type KV struct {
	Word  string `json:"word"`
	Value string `json:"value"`
}

type KeyValue []KV

func (kv KeyValue) Len() int {
	return len(kv)
}

func (kv KeyValue) Less(i, j int) bool {
	return kv[i].Word < kv[j].Word
}

func (kv KeyValue) Swap(i, j int) {
	kv[i], kv[j] = kv[j], kv[i]
}

type StateType int

const (
	_ StateType = iota
	Idle
	Inprogress
	Done
)

type TaskType int

const (
	_ TaskType = iota
	Map
	Reduce
)

type Task struct {
	Input string
	Start time.Time
	Type  TaskType
	State StateType
	R     string
}

type GetTaskRequest struct{}
type GetTaskResponse struct {
	Input string // file for map and reduce id for reduce
	Task  TaskType
	Id    string
}

type CompleteTaskRequest struct {
	Task   TaskType
	Output string
}
type CompleteTaskResponse struct{}
