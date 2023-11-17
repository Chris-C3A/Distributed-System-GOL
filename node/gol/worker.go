package gol

import (
	"sync"

	"uk.ac.bris.cs/gameoflife/util"
)

const ALIVE byte = 255
const DEAD byte = 0

type Response struct {
	World [][]byte
	AliveCells []util.Cell
	AliveCellsCount int
	CompletedTurns int
}

type Request struct {
	World [][]byte
	Turns int
}

type ControllerOperations struct{}

// global variables
var world [][]uint8
var turn = 0
var mutex sync.Mutex
var terminate = false

// RPC methods

func (s *ControllerOperations) EvolveGoL(req Request, res *Response) (err error) {
	// reset globally values when evolveGoL is called
	turn = 0
	// assign request world globally
	world = req.World

	worker(req.Turns)

	res.World = world

	res.AliveCells = calculateAliveCells()

	return
}

func (s *ControllerOperations) RequestAliveCellsCount(req Request, res *Response) (err error) {
	mutex.Lock()
	res.AliveCellsCount = len(calculateAliveCells())
	res.CompletedTurns = turn
	mutex.Unlock()

	return
}

func (s *ControllerOperations) RequestCurrentGameState(req Request, res *Response) (err error) {
	mutex.Lock()
	res.World = world
	res.CompletedTurns = turn
	mutex.Unlock()

	return
}

func (s *ControllerOperations) Shutdown(req Request, res *Response) (err error) {
	terminate = true

	server.Stop()

	return
}


func worker(turns int) {
	for turn < turns && !terminate {
		mutex.Lock()
		world = calculateNextState()
		turn++
		mutex.Unlock()
	}
}

func getNumOfLiveNeighbours(i int, j int) int {
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

func calculateNextState() [][]byte {
	// fmt.Println("world:", len(world), len(world[0]))
	newWorld := make([][]byte, len(world))

	for i := 0; i < len(newWorld); i++ {
		newWorld[i] = make([]byte, len(world[i]))
	}

	// fmt.Println("newWorld:", len(newWorld), len(newWorld[0]))

	for i := 0; i < len(newWorld); i++ {
		for j := 0; j < len(newWorld[i]); j++ {
			// get number of live neighbours
			numOfLiveNeighbours := getNumOfLiveNeighbours(i, j)

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

func calculateAliveCells() []util.Cell {
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