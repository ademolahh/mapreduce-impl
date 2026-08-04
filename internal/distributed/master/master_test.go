package master

import (
	"testing"
	"time"

	"github.com/ademolahh/map-reduce-impl/internal/shared"
)

var Status = map[shared.StateType]string{
	shared.Idle:       "Idle",
	shared.Inprogress: "Inprogress",
	shared.Done:       "Done",
}

func TestNew(t *testing.T) {
	files := []string{"a.txt"}
	nr := 1

	m := New(files, nr)

	if nr != m.nReduce {
		t.Fatalf("expected %d got %d", nr, m.nReduce)
	}

	if len(files) != m.mapCount {
		t.Fatalf("expected %d got %d", len(files), m.mapCount)
	}

	if nr != m.reduceCount {
		t.Fatalf("expected %d got %d", nr, m.reduceCount)

	}

	if len(m.task) != len(files)+nr {
		t.Fatalf("expected %d got %d", len(m.task), len(files)+nr)
	}

	tasks := m.task

	for _, v := range tasks {
		if shared.Idle != v.State {
			t.Errorf("expected %s got %s", Status[shared.Idle], Status[v.State])
		}
	}
}

func TestGetTaskReturnsMapFirst(t *testing.T) {
	files := []string{"a.txt"}
	nr := 1

	m := New(files, nr)

	var req shared.GetTaskRequest
	var res shared.GetTaskResponse

	if err := m.GetTask(req, &res); err != nil {
		t.Fatal(err)
	}

	if res.Input != "a.txt" {
		t.Fatalf("expected %s got %s", "a.txt", res.Input)
	}

	task := m.task["map-0"]

	if task == nil {
		t.Fatal("invalid task")
	}

	if task.State != shared.Inprogress {
		t.Fatalf("expected %s got %s", Status[shared.Inprogress], Status[task.State])
	}

}

func TestCompleted(t *testing.T) {
	files := []string{"a.txt"}
	nr := 1

	m := New(files, nr)

	var req shared.GetTaskRequest
	var res shared.GetTaskResponse

	if err := m.GetTask(req, &res); err != nil {
		t.Fatal(err)
	}

	completeReq := shared.CompleteTaskRequest{
		Task:   shared.Reduce,
		Output: "map-0",
	}
	var completeRes shared.CompleteTaskResponse
	if err := m.CompleteTask(completeReq, &completeRes); err != nil {
		t.Fatal(err)
	}

	if m.task["map-0"].State != shared.Done {
		t.Fatalf("expected %s got %s", Status[shared.Done], Status[m.task["map-0"].State])

	}

}

func TestGetReduceTaskAfterMapComplete(t *testing.T) {
	files := []string{"a.txt"}
	nr := 1

	m := New(files, nr)

	var req shared.GetTaskRequest
	var res shared.GetTaskResponse

	if err := m.GetTask(req, &res); err != nil {
		t.Fatal(err)
	}

	if m.task["map-0"].State != shared.Inprogress {
		t.Fatalf("expected %s got %s", Status[shared.Inprogress], Status[m.task["map-0"].State])
	}

	completeReq := shared.CompleteTaskRequest{
		Task:   shared.Map,
		Output: "map-0",
	}
	var completeRes shared.CompleteTaskResponse
	if err := m.CompleteTask(completeReq, &completeRes); err != nil {
		t.Fatal(err)
	}

	var reduceReq shared.GetTaskRequest
	var reduceRes shared.GetTaskResponse

	if err := m.GetTask(reduceReq, &reduceRes); err != nil {
		t.Fatal(err)
	}

	if reduceRes.Input == "" {
		t.Fatal("empty task")
	}

	reduceTask := m.task[reduceRes.Id]

	if reduceTask == nil {
		t.Fatal("invalid task")
	}

	if reduceTask.State != shared.Inprogress {
		t.Fatalf("expected %s got %s", Status[shared.Inprogress], Status[reduceTask.State])
	}

	completeReduceReq := shared.CompleteTaskRequest{
		Task:   shared.Reduce,
		Output: reduceRes.Id,
	}
	var completeReduceRes shared.CompleteTaskResponse
	if err := m.CompleteTask(completeReduceReq, &completeReduceRes); err != nil {
		t.Fatal(err)
	}

	if reduceTask.State != shared.Done {
		t.Fatalf("expected %s got %s", Status[shared.Done], Status[reduceTask.State])
	}

	if m.reduceCount != 0 {
		t.Fatalf("expected %d got %d", 0, m.reduceCount)
	}
}

func TestCompleteTaskInvalidState(t *testing.T) {
	files := []string{"a.txt"}
	nr := 1

	m := New(files, nr)

	completeReq := shared.CompleteTaskRequest{
		Task:   shared.Map,
		Output: "map-0",
	}
	var completeRes shared.CompleteTaskResponse
	if err := m.CompleteTask(completeReq, &completeRes); err == nil {
		t.Fatal("expected error")
	}
}

func TestDone(t *testing.T) {
	m := New([]string{"a.txt"}, 1)

	m.task["reduce-0"].State = shared.Inprogress

	req := shared.CompleteTaskRequest{
		Task:   shared.Reduce,
		Output: "reduce-0",
	}

	if err := m.CompleteTask(req, &shared.CompleteTaskResponse{}); err != nil {
		t.Fatal(err)
	}

	if !m.Done() {
		t.Fatal("expected done")
	}
}

func TestCheckAfterRest(t *testing.T) {
	m := New([]string{"a.txt"}, 1)

	task := m.task["map-0"]
	task.State = shared.Inprogress
	task.Start = time.Now().Add(-11 * time.Second)

	go m.Checker()

	time.Sleep(2 * time.Second)

	if task.State != shared.Idle {
		t.Fatalf("expected %s got %s", Status[shared.Idle], Status[task.State])
	}
}
