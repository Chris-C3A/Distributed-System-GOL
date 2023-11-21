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

// RPC methods
func (s *WorkerOperations) EvolveGoL(req stubs.Request, res *stubs.Response) (err error) {
	// fmt.Println("evolving gol called by broker")
	// reset globally values when evolveGoL is called
	turn = 0
	// assign request world globally
	world = req.World

	worker(req.Turns, req.StartY, req.EndY)

	res.World = world

	res.AliveCells = util.CalculateAliveCells(world)

	// fmt.Println("sending result to broker")

	return
}

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

func worker(turns int, startY, endY int) {
	for turn < turns && !terminate {
		mutex.Lock()
		// TODO continue from here
		world = util.CalculateNextState(world, startY, endY)
		turn++
		mutex.Unlock()
	}
}
