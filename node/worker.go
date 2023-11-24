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
	threads = 4
	mutex             sync.Mutex
)

func (s *WorkerOperations) InitWorker(req stubs.Request, res *stubs.Response) (err error) {
	fmt.Println("Worker initialized by broker")
	turn = 0
	turns = req.Turns
	world = req.World

	mutex.Lock()
	// world = util.CalculateNextState(world, req.HaloTop, req.HaloBottom)
	// world = runParallelWorkers(8)

	world = runParallelWorkers(threads, req.HaloTop, req.HaloBottom)

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
	// world = util.CalculateNextState(world, req.HaloTop, req.HaloBottom)
	world = runParallelWorkers(threads, req.HaloTop, req.HaloBottom)
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

func runParallelWorkers(threads int, haloTop, haloBottom []byte) [][]byte {
	// create world + halos to use
	combinedWorld := append([][]byte{haloTop}, world...)
	combinedWorld = append(combinedWorld, haloBottom)

	startY := 0
	dy := len(world) / threads

	// remainder pixels when height doesn't fully divide with the number of threads
	remainderY := len(world) % threads

	// initialize slice of output channels
	workerOuts := make([]chan [][]byte, threads)

	// create each worker channel
	for i := range workerOuts {
		workerOuts[i] = make(chan [][]byte)
	}

	// start go routine for each thread
	for i := 0; i < threads; i++ {
		// at the last thread add the remainder pixels to be processed
		if i == threads-1 {
			go worker(combinedWorld, workerOuts[i], startY, startY+dy+remainderY)
		} else {
			// otherwise process pixels normally
			go worker(combinedWorld, workerOuts[i], startY, startY+dy)
		}

		// update startY
		startY += dy
	}

	var newWorld [][]byte

	// recollect world data from the worker channels
	for i := 0; i < threads; i++ {
		receivedData := <-workerOuts[i]

		// append data to newWorld
		newWorld = append(newWorld, receivedData...)
	}

	// // copy newWorld into world to process next state
	// mutex.Lock()
	// copy(world, newWorld)
	// mutex.Unlock()
	return newWorld

}

// worker go routine function
func worker(combinedWorld[][]byte, out chan [][]byte, startY, endY int) {
	height := endY - startY
	newWorld := util.MakeWorld(len(world[0]), height)

	// calculates next state
	// uses global variable world

	for i := startY; i < endY; i++ {
		for j := 0; j < len(combinedWorld[i]); j++ {
			// get number of live neighbours
			numOfLiveNeighbours := util.GetNumOfLiveNeighbours(combinedWorld, i+1, j)

			// rules of the game of life
			if combinedWorld[i+1][j] == util.ALIVE {
				if numOfLiveNeighbours < 2 || numOfLiveNeighbours > 3 {
					newWorld[i-startY][j] = util.DEAD
				} else {
					newWorld[i-startY][j] = util.ALIVE
				}
			} else {
				if numOfLiveNeighbours == 3 {
					newWorld[i-startY][j] = util.ALIVE
				}
			}
		}
	}

	// send the final state to out channel
	out <- newWorld
}