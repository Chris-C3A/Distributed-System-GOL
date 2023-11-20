package node

import (
	"fmt"
	"sync"

	"uk.ac.bris.cs/node/util"
)

const ALIVE byte = 255
const DEAD byte = 0

// TODO put in stubs
type Response struct {
	World [][]byte
	AliveCells []util.Cell
	AliveCellsCount int
	CompletedTurns int
}

type Request struct {
	World [][]byte
	Turns int
	StartY int
	EndY int
	Workers []string
}

type WorkerOperations struct{}

// global variables
var world [][]uint8
var turn = 0
var mutex sync.Mutex
var terminate = false

// RPC methods

func (s *WorkerOperations) EvolveGoL(req Request, res *Response) (err error) {
	fmt.Println("evolving gol called by broker")
	// reset globally values when evolveGoL is called
	turn = 0
	// assign request world globally
	world = req.World

	worker(req.Turns, req.StartY, req.EndY)

	res.World = world

	res.AliveCells = calculateAliveCells()

	fmt.Println("sending result to broker")

	return
}

func (s *WorkerOperations) RequestAliveCellsCount(req Request, res *Response) (err error) {
	mutex.Lock()
	res.AliveCellsCount = len(calculateAliveCells())
	res.CompletedTurns = turn
	mutex.Unlock()

	return
}

func (s *WorkerOperations) RequestCurrentGameState(req Request, res *Response) (err error) {
	mutex.Lock()
	res.World = world
	res.CompletedTurns = turn
	mutex.Unlock()

	return
}

func (s *WorkerOperations) Shutdown(req Request, res *Response) (err error) {
	terminate = true

	server.Stop()

	return
}


func worker(turns int, startY, endY int) {
	for turn < turns && !terminate {
		mutex.Lock()
		// TODO continue from here
		world = calculateNextState(startY, endY)
		turn++
		mutex.Unlock()
	}
}

// GOL Logic

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

func calculateNextState(startY, endY int) [][]byte {
	// fmt.Println("world:", len(world), len(world[0]))
	height := endY - startY
	width := len(world[0])

	newWorld := make([][]byte, height)

	for i := 0; i < height; i++ {
		newWorld[i] = make([]byte, width)
	}

	for i := startY; i < endY; i++ {
		for j := 0; j < width; j++ {
			// get number of live neighbours
			numOfLiveNeighbours := getNumOfLiveNeighbours(i, j)

			// rules of the game of life
			if world[i][j] == ALIVE {
				if numOfLiveNeighbours < 2 || numOfLiveNeighbours > 3 {
					newWorld[i-startY][j] = DEAD
				} else {
					newWorld[i-startY][j] = ALIVE
				}
			} else {
				if numOfLiveNeighbours == 3 {
					newWorld[i-startY][j] = ALIVE
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