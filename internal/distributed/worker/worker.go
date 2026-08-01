package worker

import (
	"fmt"
	"io"
	"net/rpc"
	"os"
	"time"

	"github.com/ademolahh/map-reduce-impl/internal/shared"
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

func (w *Worker) Work(programPath string) error {
	mps, rdc, err := shared.Fetch(programPath)
	mapResultFolder := "mapresult"
	if err != nil {
		return err
	}
	for {
		var response shared.GetTaskResponse
		err := w.client.Call("Master.GetTask", &shared.GetTaskRequest{}, &response)
		if err == rpc.ErrShutdown || err == io.ErrUnexpectedEOF {
			fmt.Println("master shutdown...")
			break
		}

		if err != nil {
			fmt.Println("Get Call Failed", err)
			continue
		}

		if response.Task == shared.Map {
			err = w.mapFunc(response, mps)
			if err == rpc.ErrShutdown || err == io.ErrUnexpectedEOF {
				break
			}
			continue

		} else {
			dirs, err := os.ReadDir(mapResultFolder)
			if err != nil {
				continue
			}

			if response.Input == "" {
				time.Sleep(3 * time.Second)
				continue
			}

			if err := w.reduceFunc(response, mapResultFolder, dirs, rdc); err != nil {
				fmt.Println("reduce call failed")
				continue
			}

		}

	}

	return nil
}
