package master

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/ademolahh/map-reduce-impl/shared"
)

type Master struct {
	mu          sync.Mutex
	mapCount    int
	reduceCount int
	task        map[string]*shared.Task
}

func New(files []string, nr int) *Master {
	task := make(map[string]*shared.Task)

	for i, file := range files {
		t := fmt.Sprintf("map-%d", i)
		task[t] = &shared.Task{
			Input: file,
			Type:  shared.Map,
			State: shared.Idle,
		}
	}

	for i := range nr {
		r := fmt.Sprintf("reduce-%d", i)
		task[r] = &shared.Task{
			Type:  shared.Reduce,
			State: shared.Idle,
			R:     strconv.Itoa(i),
		}
	}

	return &Master{
		mapCount:    len(files),
		reduceCount: nr,
		task:        task,
	}
}

func (m *Master) GetTask(req shared.GetTaskRequest, res *shared.GetTaskResponse) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	var input string
	var task shared.TaskType
	var taskId string

	for k, t := range m.task {
		if m.mapCount > 0 {
			if t.Type == shared.Map && t.State == shared.Idle {
				// response
				input = t.Input
				task = shared.Map

				// master details
				t.Start = time.Now()
				t.State = shared.Inprogress
				taskId = k

				break
			}
		} else {
			if t.State == shared.Idle && t.Type == shared.Reduce {
				// response
				input = t.R
				task = shared.Reduce

				// master details
				t.Start = time.Now()
				t.State = shared.Inprogress
				taskId = k

				break
			}
		}
	}

	*res = shared.GetTaskResponse{
		Input: input,
		Task:  task,
		Id:    taskId,
	}
	return nil
}

func (m *Master) CompleteTask(req shared.CompleteTaskRequest, res *shared.CompleteTaskResponse) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.task[req.Output].State != shared.Inprogress {
		return fmt.Errorf("invalid state")
	}

	m.task[req.Output].State = shared.Done
	if req.Task == shared.Map {
		m.mapCount--
	} else {
		m.reduceCount--
	}

	return nil
}

func (m *Master) Checker() {
	for {
		time.Sleep(1 * time.Second)
		m.mu.Lock()
		for _, t := range m.task {
			if t.State == shared.Inprogress && time.Since(t.Start) > 10*time.Second {
				t.State = shared.Idle
			}
		}
		m.mu.Unlock()
	}
}
