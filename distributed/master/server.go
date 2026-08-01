package master

import (
	"fmt"
	"io"
	"net"
	"net/rpc"
	"os"
	"time"
)

func Serve(master *Master) error {

	if err := rpc.Register(master); err != nil {
		return err
	}

	socketName := "/tmp/mr-1"
	os.Remove(socketName)

	listener, err := net.Listen("unix", socketName)
	if err != nil {
		return err
	}
	fmt.Println("Started")

	go func() {
		for !master.Done() {
			time.Sleep(1 * time.Second)
		}

		os.Exit(1)
	}()

	go master.Checker()

	for {
		conn, err := listener.Accept()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		go rpc.ServeConn(conn)
	}

	return nil
}
