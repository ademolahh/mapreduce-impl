package master

import (
	"io"
	"net"
	"net/rpc"
)

func Serve() error {
	master := New([]string{}, 1)
	if err := rpc.Register(master); err != nil {
		return err
	}

	listener, err := net.Listen("unix", "var/tmp/mr-1")
	if err != nil {
		return err
	}

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
