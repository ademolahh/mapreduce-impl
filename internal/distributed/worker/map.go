package worker

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"os"
	"strconv"

	"github.com/ademolahh/map-reduce-impl/internal/shared"
)

func (w *Worker) mapFunc(response shared.GetTaskResponse, mps func(string, string) []shared.KV) error {
	data, err := os.ReadFile(response.Input)
	if err != nil {
		return err

	}

	kv := mps(response.Input, string(data))

	if err := loop(kv, response.Id, response.NReduce); err != nil {
		return err
	}

	req := shared.CompleteTaskRequest{
		Task:   response.Task,
		Output: response.Id,
	}
	err = w.client.Call("Master.CompleteTask", req, &shared.CompleteTaskResponse{})

	return err
}

func loop(kv []shared.KV, id string, nr int) error {

	result := make(map[int][]shared.KV)

	for _, v := range kv {
		key := v.Word
		ptn := r(key) % nr
		result[ptn] = append(result[ptn], v)
	}

	for k, v := range result {
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}

		pattern := fmt.Sprintf("%s", id)

		tmp, err := os.CreateTemp("tmp", pattern)
		if err != nil {
			return err
		}
		tmpName := tmp.Name()
		defer os.Remove(tmpName)

		if _, err := tmp.Write(data); err != nil {
			tmp.Close()
			return err
		}

		if err := tmp.Close(); err != nil {
			return err
		}

		err = os.Rename(tmp.Name(), "mapresult/"+pattern+"-"+strconv.Itoa(k)+".json")

		if err != nil {
			return err
		}

	}

	return nil
}

func r(key string) int {
	h := fnv.New32()
	h.Write([]byte(key))
	return int(h.Sum32() & 0xfffffff)
}
