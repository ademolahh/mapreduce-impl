package worker

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"net/rpc"
	"os"
	"strconv"

	"github.com/ademolahh/map-reduce-impl/shared"
)

type Worker struct {
	client *rpc.Client
}

func New() (*Worker, error) {
	client, err := rpc.Dial("unix", "/tmp/mr-1")
	if err != nil {
		return nil, err
	}

	return &Worker{client: client}, nil
}

func (w *Worker) Work() error {
	mps, rdc, err := shared.Fetch("distributed/worker/wordcount.so")
	if err != nil {
		return err
	}
	for {
		var response shared.GetTaskResponse
		err := w.client.Call("Master.GetTask", &shared.GetTaskRequest{}, &response)
		if err == io.ErrUnexpectedEOF {
			break
		}

		if err != nil {
			fmt.Println("Get Call Failed", err)
			break
		}

		input := response.Input

		if response.Task == shared.Map {
			data, err := os.ReadFile(input)
			if err != nil {
				fmt.Println("read failed", err)
				break
			}
			kv := mps(response.Input, string(data))

			if err := Loop(kv, response.Id); err != nil {
				fmt.Println("Intermediate file failed", err)
				break
			}

			req := shared.CompleteTaskRequest{
				Task:   response.Task,
				Output: response.Id,
			}
			err = w.client.Call("Master.CompleteTask", req, &shared.CompleteTaskResponse{})
			if err == io.ErrUnexpectedEOF {
				break
			}

			if err != nil {
				fmt.Println("Completion Call Failed", err)
				break
			}

		} else {
			break
			rdc("", []string{})
		}

	}

	return nil
}

func Loop(kv []shared.KV, id string) error {

	result := make(map[int][]shared.KV)

	for _, v := range kv {
		key := v.Word
		ptn := r(key) % 3
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
