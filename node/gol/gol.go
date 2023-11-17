package gol

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

// CloseDetectConn is a wrapper for net.Conn that detects client disconnects.
// type CloseDetectConn struct {
// 	net.Conn
// 	onClose func()
// }

// // Close overrides the Close method of net.Conn.
// func (c *CloseDetectConn) Close() error {
// 	c.onClose()
// 	return c.Conn.Close()
// }

var server *Server

// Run starts the processing of Game of Life. It should initialise channels and goroutines.
func Run() {
	// listen to connections
	// pAddr := flag.String("port", "8030", "Port to listen on")

	// fmt.Println("Worker node listening on", *pAddr)

	// flag.Parse()
	// rand.Seed(time.Now().UnixNano())

	// rpc.Register(&ControllerOperations{})

	// ln, _ := net.Listen("tcp", ":"+*pAddr)
	// defer ln.Close()

	// rpc.Accept(ln)

	// listen to connections
	// pAddr := flag.String("port", "8030", "Port to listen on")

	// fmt.Println("Worker node listening on", *pAddr)

	// flag.Parse()
	// rand.Seed(time.Now().UnixNano())

	// rpc.Register(&ControllerOperations{})

	// ln, err := net.Listen("tcp", ":"+*pAddr)
	// if err != nil {
	// 	fmt.Println("Error starting server:", err)
	// 	return
	// }
	// defer ln.Close()

	// for {
	// 	conn, err := ln.Accept()
	// 	if err != nil {
	// 		fmt.Println("Error accepting connection:", err)
	// 		continue
	// 	}

	// 	closeDetectConn := &CloseDetectConn{
	// 		Conn: conn,
	// 		onClose: func() {
	// 			fmt.Println("Client disconnected")
	// 			// Additional cleanup or handling here
	// 		},
	// 	}

	// 	go rpc.ServeConn(closeDetectConn)
	// }

	server = NewServer()

	pAddr := flag.String("port", "8030", "Port to listen on")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	// start server
	if err := server.Start(*pAddr); err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

}
