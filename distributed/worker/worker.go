package worker

import (
	"net/rpc"
)

type Worker struct {
	client *rpc.Client
}

func New() *Worker {
	client, err := rpc.Dial("unix", "var/tmp/mr-1")
	if err != nil {
		panic(err)
	}

	return &Worker{client: client}
}
