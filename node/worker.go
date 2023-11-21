package main

import (
	"sync"

	"uk.ac.bris.cs/gameoflife/stubs"
	"uk.ac.bris.cs/gameoflife/util"
)

type WorkerOperations struct{}

// global variables
var world [][]uint8
var turn = 0
var mutex sync.Mutex
var terminate = false

// will have extra top row and bottom row to use
// startY and endY will not be needed anymore since only the actual part of the world that needs to be processed is being sent


// RPC method to assign the world before evolving GoL
// run worker
// probably have a chan for the halos when received and sent back
// func InitWorker(req stubs.Request, res *stubs.Response) (err error) {
// 	fmt.Println("Worker initalized by broker")
// 	// reset globally values when evolveGoL is called
// 	turn = 0

// 	mutex.Lock()
// 	world = req.World
// 	mutex.Unlock()

// 	return 
// }

// rpc call to send and receive halos

// RPC methods
func (s *WorkerOperations) EvolveGoL(req stubs.Request, res *stubs.Response) (err error) {
	turn = 0

	// assign request world globally
	world = req.World

	worker(req.Turns, req.HaloTop, req.HaloBottom)

	// exclude halos
	// mutex.Lock()
	// exclude halos (no need)
	// world = world[1:len(world)-1]
	// mutex.Unlock()
	res.World = world

	// return halos separately
	// res.HaloTop = world[0]
	// res.HaloBottom = world[len(world)-1]

	// no need for this
	// res.AliveCells = util.CalculateAliveCells(world)

	return
}


// not used anymore
func (s *WorkerOperations) RequestAliveCellsCount(req stubs.Request, res *stubs.Response) (err error) {
	mutex.Lock()
	res.AliveCellsCount = len(util.CalculateAliveCells(world))
	res.CompletedTurns = turn
	mutex.Unlock()

	return
}

func (s *WorkerOperations) RequestCurrentGameState(req stubs.Request, res *stubs.Response) (err error) {
	mutex.Lock()
	res.World = world
	res.CompletedTurns = turn
	mutex.Unlock()

	return
}

func (s *WorkerOperations) Shutdown(req stubs.Request, res *stubs.Response) (err error) {
	terminate = true

	server.Stop()

	return
}


// worker function
func worker(turns int, haloTop, haloBottom []byte) {
	for turn < turns && !terminate {
		mutex.Lock()
		world = util.CalculateNextState(world, haloTop, haloBottom)
		turn++
		// after each turn update the halos and send to the broker
		// wait for halos
		mutex.Unlock()
	}
}
