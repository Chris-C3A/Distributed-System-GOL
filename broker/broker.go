package main

import (
	"fmt"
	"net/rpc"
	"sync"

	"uk.ac.bris.cs/gameoflife/stubs"
	"uk.ac.bris.cs/gameoflife/util"
)

type ControllerOperations struct{}

var WorkerEvolveGoL = "WorkerOperations.EvolveGoL"

// global variables
var world [][]uint8
var mutex sync.Mutex
var turn int

// var terminate = false
// var workers []string
// var workersClient []*rpc.Client


// initalize Broker
func (s *ControllerOperations) Broker(req stubs.Request, res *stubs.Response) (err error) {
	// worker ports
	// workers = req.Workers

	// // connect to worker and add client in array
	// for _, worker := range workers {
	// 	// connect to worker
	// 	workerClient := connect(worker)

	// 	// add to array of worker clients
	// 	workersClient = append(workersClient, workerClient)
	// }

	return
}

func (s *ControllerOperations) EvolveGoL(req stubs.Request, res *stubs.Response) (err error) {
	// assign request world globally
	turn = 0
	world = req.World

	//! not used atm
	// numOfWorkers := req.NumOfWorkers


	fmt.Println("received from controller")

	fmt.Println("sending to all workers:", workersClient)

	dy := len(world)/len(workersClient)
	// add remainder as well
	remainderY := len(world) % len(workersClient)

	for turn < req.Turns {

		startY := 0

		// initialize slice of output channels
		workerCalls := make([]*rpc.Call, len(workersClient))
		workerResponses := make([]*stubs.Response, len(workersClient))


		// break down world and send to each worker its part
		for i, workerClient := range workersClient {
			// execute 1 turn on worker
			// send request asynchronously
			workerResponses[i] = new(stubs.Response)
			
			var request stubs.Request

			// partition world for worker + halos
			var subworld [][]byte
			var haloTop, haloBottom []byte

			if i == len(workersClient) - 1 {
				subworld = world[startY:startY+dy+remainderY]
			} else {
				subworld = world[startY:startY+dy]
			}

			// sepcial case for one worker
			if i == 0 && i == len(workersClient) - 1 {
				haloTop = world[len(world)-1]
				haloBottom = world[0]
			} else if i == 0 {
				haloTop = world[len(world)-1]
				haloBottom = world[startY+dy]
			} else if i == len(workersClient) - 1 {
				haloTop = world[startY-1]
				haloBottom = world[0]
			} else {
				haloTop = world[startY-1]
				haloBottom = world[startY+dy]
			}

			// add subworld + halo top and bottom rows
			request = stubs.Request{World: subworld, HaloTop: haloTop, HaloBottom: haloBottom, Turns: 1}

			workerCalls[i] = workerClient.Go(WorkerEvolveGoL, request, workerResponses[i], nil)

			startY += dy
		}

		// reconstruct world?

		var newWorld [][]byte
		// wait for calls to complete and reconstruct the world
		for i, workerCall := range workerCalls {
			// wait for the request to complete (check for errors later)
			<-workerCall.Done
			newWorld = append(newWorld, workerResponses[i].World...)
		}

		// copy the newWorld
		mutex.Lock()
		copy(world, newWorld)
		mutex.Unlock()

		turn++
	}

	// add final world to response
	res.World = world

	// alive cells
	res.AliveCells = util.CalculateAliveCells(world)

	fmt.Println("Sending final result to controller")

	return
}

func (s *ControllerOperations) RequestAliveCellsCount(req stubs.Request, res *stubs.Response) (err error) {
	mutex.Lock()
	res.AliveCellsCount = len(util.CalculateAliveCells(world))
	res.CompletedTurns = turn
	mutex.Unlock()

	return
}

// func (s *ControllerOperations) RequestCurrentGameState(req stubs.Request, res *stubs.Response) (err error) {
// 	mutex.Lock()
// 	res.World = world
// 	res.CompletedTurns = turn
// 	mutex.Unlock()

// 	return
// }

// func (s *ControllerOperations) Shutdown(req stubs.Request, res *stubs.Response) (err error) {
// 	terminate = true

// 	server.Stop()

// 	return
// }


func sendToWorker(workerClient *rpc.Client, req stubs.Request, function string) *stubs.Response {

	res := new(stubs.Response)

	// call rpc function
	workerClient.Call(function, req, res)
	// workerClient.Go(function, req, res,)

	return res
}
