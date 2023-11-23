package main

import (
	"fmt"
	"sync"

	"uk.ac.bris.cs/gameoflife/stubs"
	"uk.ac.bris.cs/gameoflife/util"
)

type WorkerOperations struct{}

var (
	world             [][]uint8
	turn              = 0
	turns							int
	terminate         = false
	mutex             sync.Mutex
)

func (s *WorkerOperations) InitWorker(req stubs.Request, res *stubs.Response) (err error) {
	fmt.Println("Worker initialized by broker")
	turn = 0
	turns = req.Turns
	world = req.World

	mutex.Lock()
	world = util.CalculateNextState(world, req.HaloTop, req.HaloBottom)
	turn++
	mutex.Unlock()

	if turn == turns {
		res.World = world
		return
	}

	// // Run worker as a goroutine
	// go worker(req.Turns)

	// Wait to receive first halos
	res.HaloTop = world[0]
	res.HaloBottom = world[len(world)-1]

	fmt.Println("Initialization successful")

	return
}

func (s *WorkerOperations) HaloExchange(req stubs.Request, res *stubs.Response) (err error) {
	fmt.Println("Halo exchange called")

	mutex.Lock()
	world = util.CalculateNextState(world, req.HaloTop, req.HaloBottom)
	turn ++
	mutex.Unlock()

	if turn == turns {
		res.World = world
		return
	}

	fmt.Println("Sending halos to broker")

	// Get halos after the iteration to send back to broker
	res.HaloTop = world[0]
	res.HaloBottom = world[len(world)-1]

	return
}


func (s *WorkerOperations) RequestCurrentGameState(req stubs.Request, res *stubs.Response) (err error) {
	mutex.Lock()
	res.World = world
	res.CompletedTurns = turn
	mutex.Unlock()
	return
}
// // Worker function
// func worker(turns int) {
// 	for turn < turns && !terminate {
// 		mutex.Lock()

// 		// Receive halos to include in the next iteration calculation
// 		haloTop := <-haloTopChan
// 		haloBottom := <-haloBottomChan

// 		world = util.CalculateNextState(world, haloTop, haloBottom)

// 		// Send halos to other workers
// 		go func() {
// 			haloTopToSend <- world[0]
// 			haloBottomToSend <- world[len(world)-1]
// 		}()

// 		turn++
// 		mutex.Unlock()
// 	}

// 	// Send done channel
// 	done = true
// }