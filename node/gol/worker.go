package gol

import (
	"flag"
	"math/rand"
	"net"
	"net/rpc"
	"time"
	"uk.ac.bris.cs/gameoflife/util"
)

const ALIVE byte = 255
const DEAD byte = 0

func (s *ControllerOperations) EvolveGoL(req Request, res *Response, p Params) (err error) {
	res.Message = calculateNextState(p, req.Message)
	return
}

type Response struct {
	Message [][]byte
}

type Request struct {
	Message [][]byte
}

type ControllerOperations struct{}

func worker(p Params) {
	handleConnections()
	turn := 0

	// initalize world
	world := make([][]byte, p.ImageHeight)
	for i := range world {
		world[i] = make([]byte, p.ImageWidth)
	}

	// TODO: Execute all turns of the Game of Life.
	for ; turn < p.Turns; turn++ {
		world = calculateNextState(p, world)
	}
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
	pAddr := flag.String("port", "8030", "Port to listen on")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	rpc.Register(&ControllerOperations{})
	ln, _ := net.Listen("tcp", ":"+*pAddr)
	defer ln.Close()
	rpc.Accept(ln)
}
