package gol

import (
	"fmt"
	"net/rpc"

	"uk.ac.bris.cs/gameoflife/util"
)

type distributorChannels struct {
	events     chan<- Event
	ioCommand  chan<- ioCommand
	ioIdle     <-chan bool
	ioFilename chan<- string
	ioOutput   chan<- uint8
	ioInput    <-chan uint8
}

type Response struct {
	World [][]byte
	AliveCells []util.Cell
	// CompletedTurns int
}

type Request struct {
	World [][]byte
	Turns int
}

const ALIVE byte = 255
const DEAD byte = 0

var EvolveGoL = "ControllerOperations.EvolveGoL"

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels) {
	// TODO: Create a 2D slice to store the world.
		// send input command
	c.ioCommand <- ioInput
	c.ioFilename <- fmt.Sprintf("%dx%d", p.ImageWidth, p.ImageHeight)

	// initialise world
	world := make([][]byte, p.ImageHeight)
	for i := range world {
		world[i] = make([]byte, p.ImageWidth)
	}

	// store image bytes into world
	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			world[y][x] = <-c.ioInput
		}
	}

	fmt.Println("sending to server")

	response := handleConnections(world, p)
	fmt.Println("received response")

	// TODO: Report the final state using FinalTurnCompleteEvent.
	// FinalTurnComplete Event
	c.events <- FinalTurnComplete{CompletedTurns: p.Turns, Alive: response.AliveCells}

	// put each byte of final world into output channel
	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			c.ioOutput <- response.World[y][x]
		}
	}

	// put ioOutput into command channel
	c.ioCommand <- ioOutput
	// output file
	outFileName := fmt.Sprintf("%dx%dx%d", p.ImageWidth, p.ImageHeight, p.Turns)
	// send to channel
	c.ioFilename <- outFileName

	// send imageoutput complete event
	c.events <- ImageOutputComplete{CompletedTurns: p.Turns, Filename: outFileName}

	// Make sure that the Io has finished any output before exiting.
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	c.events <- StateChange{p.Turns, Quitting}

	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}

func handleConnections(world [][]uint8, p Params) *Response {
	// server := flag.String("server", "127.0.0.1:8030", "IP:port string to connect to as server")
	// flag.Parse()
	server := "127.0.0.1:8030"
	fmt.Println("Server: ", server)

	client, _ := rpc.Dial("tcp", server)

	defer client.Close()

	request := Request{World: world, Turns: p.Turns}

	response := new(Response)

	client.Call(EvolveGoL, request, response)

	return response
}
