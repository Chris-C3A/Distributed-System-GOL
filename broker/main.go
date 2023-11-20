package main

import (
	"uk.ac.bris.cs/broker/gol"
)

type ControllerOperations struct{}

func main() {
	// // listen to connections
	// pAddr := flag.String("port", "8030", "Port to listen on")

	// fmt.Println("Worker node listening on", *pAddr)

	// flag.Parse()
	// rand.Seed(time.Now().UnixNano())

	// rpc.Register(&ControllerOperations{})

	// ln, _ := net.Listen("tcp", ":"+*pAddr)
	// defer ln.Close()

	// rpc.Accept(ln)
	gol.Run()
}