package gol

import (
	"flag"
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

const ALIVE byte = 255
const DEAD byte = 0

var EvolveGoL = "ControllerOperations.Evolve"

type Response struct {
	Message string
}

type Request struct {
	Message string
}

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels) {

	turn := 0
	// send input command
	c.ioCommand <- ioInput

	//wait for idle
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	fmt.Println("running")

	// initalize world
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

	// TODO: Execute all turns of the Game of Life.
	for ; turn < p.Turns; turn++ {
		world = calculateNextState(p, world)
	}

	// TODO: Report the final state using FinalTurnCompleteEvent.

	// put each byte of final wrold into output channel
	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			c.ioOutput <- world[y][x]
		}
	}

	// put ioOuput into command channel
	c.ioCommand <- ioOutput

	// finalturncompleteevent
	c.events <- FinalTurnComplete{CompletedTurns: turn, Alive: calculateAliveCells(p, world)}

	// Make sure that the Io has finished any output before exiting.
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	c.events <- StateChange{turn, Quitting}

	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}

func getNumOfLiveNeighbours(world [][]byte, i int, j int) int {
	numOfLiveNeighbours := 0

	// positive modulus
	up := ((i-1)%len(world) + len(world)) % len(world)
	down := ((i+1)%len(world) + len(world)) % len(world)
	right := ((j+1)%len(world[i]) + len(world[i])) % len(world[i])
	left := ((j-1)%len(world[i]) + len(world[i])) % len(world[i])

	neighbours := [8]byte{world[up][j], world[down][j], world[i][left], world[i][right], world[up][left], world[up][right], world[down][right], world[down][left]}

	for _, neighbour := range neighbours {
		if neighbour == ALIVE {
			numOfLiveNeighbours++
		}
	}

	return numOfLiveNeighbours

}

func calculateNextState(p Params, world [][]byte) [][]byte {
	newWorld := make([][]byte, p.ImageHeight)

	for i := 0; i < len(newWorld); i++ {
		newWorld[i] = make([]byte, p.ImageWidth)
	}

	for i := 0; i < len(newWorld); i++ {
		for j := 0; j < len(newWorld[i]); j++ {
			// get number of live neighbours
			numOfLiveNeighbours := getNumOfLiveNeighbours(world, i, j)

			// rules of the game of life
			if world[i][j] == ALIVE {
				if numOfLiveNeighbours < 2 || numOfLiveNeighbours > 3 {
					newWorld[i][j] = DEAD
				} else {
					newWorld[i][j] = ALIVE
				}
			} else {
				if numOfLiveNeighbours == 3 {
					newWorld[i][j] = ALIVE
				}
			}
		}
	}

	return newWorld

}

func calculateAliveCells(p Params, world [][]byte) []util.Cell {
	var aliveCells []util.Cell

	for y := 0; y < len(world); y++ {
		for x := 0; x < len(world[y]); x++ {
			if world[y][x] == ALIVE {
				// add cell coordinates to aliveCells slice
				aliveCells = append(aliveCells, util.Cell{X: x, Y: y})
			}
		}
	}

	return aliveCells
}

func handleConnections() {
	server := flag.String("server", "127.0.0.1:8030", "IP:port string to connect to as server")
	flag.Parse()
	fmt.Println("Server: ", *server)
	client, _ := rpc.Dial("tcp", *server)
	defer client.Close()
	request := Request{Message: "Hello"}
	response := new(Response)
	client.Call(EvolveGoL, request, response)
	fmt.Println("Responded: " + response.Message)
}
