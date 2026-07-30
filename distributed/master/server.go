package master

import (
	"fmt"
	"io"
	"net"
	"net/rpc"
	"os"
	"strconv"
)

func Serve(files []string, nr int) error {
	master := New(files, nr)
	if err := rpc.Register(master); err != nil {
		return err
	}

	socketName := c()
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

func c() string {
	sct := "/tmp/mr-"
	sct += strconv.Itoa(os.Getuid())
	return sct
}
