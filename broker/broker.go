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
			if i == len(workersClient) - 1 {
				request = stubs.Request{World: world, StartY: startY, EndY: startY + dy + remainderY, Turns: 1}
			} else {
				request = stubs.Request{World: world, StartY: startY, EndY: startY + dy, Turns: 1}
			}
			// request := stubs.Request{World: world, StartY: startY, EndY: startY + dy, Turns: 1}

			workerCalls[i] = workerClient.Go(WorkerEvolveGoL, request, workerResponses[i], nil)


			// response := sendToWorker(workerClient, stubs.Request{World: world, StartY: startY, EndY: startY + dy,  Turns: 1}, WorkerEvolveGoL)

			// // add response 
			// responses = append(responses, response)

			startY += dy
		}

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

	// // break down world and send to each worker its part
	// for _, workerClient := range workersClient {
	// 	response := sendToWorker(workerClient, stubs.Request{World: req.World, Turns: 1}, WorkerEvolveGoL)
	// 	responses = append(responses, response)
	// }

	// // TODO test assuming one response
	// res.World = responses[0].World

	// res.AliveCells = responses[0].AliveCells

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


// func worker(turns int) {
// 	for turn < turns && !terminate {
// 		mutex.Lock()
// 		// world = calculateNextState()
// 		turn++
// 		mutex.Unlock()
// 	}
// }

func sendToWorker(workerClient *rpc.Client, req stubs.Request, function string) *stubs.Response {

	res := new(stubs.Response)

	// call rpc function
	workerClient.Call(function, req, res)
	// workerClient.Go(function, req, res,)

	return res
}
