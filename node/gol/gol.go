package gol

import (
	"flag"
	"math/rand"
	"net"
	"net/rpc"
	"time"
)

// Params provides the details of how to run the Game of Life and which image to load.
// type Params struct {
// 	Turns       int
// 	Threads     int
// 	ImageWidth  int
// 	ImageHeight int
// }

// Run starts the processing of Game of Life. It should initialise channels and goroutines.
func Run() {
	pAddr := flag.String("port", "8030", "Port to listen on")

	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	rpc.Register(&ControllerOperations{})

	ln, _ := net.Listen("tcp", ":"+*pAddr)
	defer ln.Close()

	rpc.Accept(ln)
}
