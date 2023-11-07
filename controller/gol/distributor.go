package gol

import (
	"flag"
	"fmt"
	"net/rpc"
)

type distributorChannels struct {
	events     chan<- Event
	ioCommand  chan<- ioCommand
	ioIdle     <-chan bool
	ioFilename chan<- string
	ioOutput   chan<- uint8
	ioInput    <-chan uint8
}

const ALIVE byte = 255
const DEAD byte = 0

var EvolveGoL = "ControllerOperations.Evolve"

type Response struct {
	Message Params
}

type Request struct {
	Message Params
}

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels) {
	handleConnections(p)

	// TODO: Create a 2D slice to store the world.

	turn := 0

	// TODO: Execute all turns of the Game of Life.

	// TODO: Report the final state using FinalTurnCompleteEvent.

	// Make sure that the Io has finished any output before exiting.
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	c.events <- StateChange{turn, Quitting}

	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}

func handleConnections(p Params) {
	server := flag.String("server", "127.0.0.1:8030", "IP:port string to connect to as server")
	flag.Parse()
	fmt.Println("Server: ", *server)
	client, _ := rpc.Dial("tcp", *server)
	defer client.Close()
	request := Request{Message: p}
	response := new(Response)
	client.Call(EvolveGoL, request, response)
}
