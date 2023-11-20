package main

import "uk.ac.bris.cs/node/node"


var Test int = 1

func main() {

	node.Run()

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
