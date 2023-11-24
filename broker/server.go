package main

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
)

// CloseDetectConn is a wrapper for net.Conn that detects client disconnects.
type CloseDetectConn struct {
	net.Conn
	onClose func()
}

// Close overrides the Close method of net.Conn.
func (c *CloseDetectConn) Close() error {
	c.onClose()
	return c.Conn.Close()
}

type Server struct {
	listener net.Listener
	server   *rpc.Server
	wg       sync.WaitGroup
}

func NewServer() *Server {
	return &Server{
		server: rpc.NewServer(),
	}
}

func (s *Server) Start(port string) error {
	rpc.Register(new(ControllerOperations))

	listener, err := net.Listen("tcp", ":"+port)
	fmt.Println("Broker listening on", port)

	if err != nil {
		return err
	}

	defer listener.Close()

	s.listener = listener

	// for {
	conn, err := s.listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
		// continue
	}

	closeDetectConn := &CloseDetectConn{
		Conn: conn,
		onClose: func() {
			fmt.Println("Client disconnected")
			// Additional cleanup or handling here
			// TODO
		},
	}

	// rpc.ServeConn(closeDetectConn)
	s.server.ServeConn(closeDetectConn)
	
	// for !terminate {

	// }

	// for !terminate {
	// }
	// go rpc.ServeConn(conn)
		// go s.server.ServeConn(closeDetectConn)
	// }
	fmt.Println("shutdown server")
	return nil
}

func (s *Server) Stop() {
	if s.listener != nil {
		_ = s.listener.Close()
	}

	// Wait for ongoing RPC calls to finish
	s.wg.Wait()
}