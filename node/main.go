package main

import (
	"runtime"

	"uk.ac.bris.cs/gameoflife/gol"
)

// main is the function called when starting Game of Life with 'go run .'
func main() {
	runtime.LockOSThread()

	gol.Run()

	// open rpc server on certain port...

	// listen to connections

	// 
	// pAddr := flag.String("port", "8030", "Port to listen on")

	// flag.Parse()
	// rand.Seed(time.Now().UnixNano())

	// rpc.Register(&ControllerOperations{})

	// ln, _ := net.Listen("tcp", ":"+*pAddr)
	// defer ln.Close()

	// rpc.Accept(ln)
}
