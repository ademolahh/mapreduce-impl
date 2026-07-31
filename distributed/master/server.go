package master

import (
	"fmt"
	"io"
	"net"
	"net/rpc"
	"os"
)

func Serve(files []string, nr int) error {
	master := New(files, nr)
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
