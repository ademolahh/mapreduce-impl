package master

import (
	"fmt"
	"sync"
	"time"
)

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
	R     int
}

type Master struct {
	mu          sync.Mutex
	mapCount    int
	reduceCount int
	task        map[string]*Task
}

func New(files []string, nr int) *Master {
	task := make(map[string]*Task)

	for i, file := range files {
		t := fmt.Sprintf("map-%d", i)
		task[t] = &Task{
			Input: file,
			Type:  Map,
			State: Idle,
		}
	}

	for i := range nr {
		r := fmt.Sprintf("reduce-%d", i)
		task[r] = &Task{
			Type:  Reduce,
			State: Idle,
			R:     i,
		}
	}

	return &Master{
		mapCount:    len(files),
		reduceCount: nr,
		task:        task,
	}
}

type GetTaskRequest struct{}
type GetTaskResponse struct{}

func (m *Master) GetTask(req GetTaskRequest, res *GetTaskResponse) error {
	return nil
}

type CompleteTaskRequest struct{}
type CompleteTaskResponse struct{}

func (m *Master) CompleteTask(req CompleteTaskRequest, res CompleteTaskResponse) error {
	return nil
}
