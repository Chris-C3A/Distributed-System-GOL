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
var WorkerInit = "WorkerOperations.InitWorker"
var WokerHaloExchange = "WorkerOperations.HaloExchange"
var WorkerRequestCurrentGameState = "WorkerOperations.RequestCurrentGameState"
var WorkerShutdown = "WorkerOperatiions.Shutdown"
var WorkerStop = "WorkerOperatiions.WorkerStop"

// global variables
var world [][]uint8
var mutex sync.Mutex
var wg sync.WaitGroup
var turn int
var terminate = false
var shutdown = false
var running = true


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
	terminate = false
	shutdown = false
	running = true

	//! not used atm
	// numOfWorkers := req.NumOfWorkers
	fmt.Println("received from controller")

	fmt.Println("sending to all workers:", workersClient)

	dy := len(world)/len(workersClient)
	// add remainder as well
	remainderY := len(world) % len(workersClient)

	startY := 0

	// initialize slice of output channels
	// workerCalls := make([]*rpc.Call, len(workersClient))
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
		request = stubs.Request{World: subworld, HaloTop: haloTop, HaloBottom: haloBottom, Turns: req.Turns}

		wg.Add(1)

		// asynchronously send rpc calls
		go func(i int, workerClient *rpc.Client, req stubs.Request) {
			defer wg.Done()

			workerResponses[i] = sendToWorker(workerClient, req, WorkerInit)

		}(i, workerClient, request)

		// workerCalls[i] = workerClient.Go(WorkerInit, request, workerResponses[i], nil)

		startY += dy
	}

	turn ++

	// for all turns proccess halo exchanges

	for turn < req.Turns && !terminate {
		// wait for all rpc calls to finish
		wg.Wait()

		// receive halos
		// reset halos arrays
		var halosTop [][]byte
		var halosBottom [][]byte

		for _, workerResponse := range workerResponses {
			halosTop = append(halosTop, workerResponse.HaloTop)
			halosBottom = append(halosBottom, workerResponse.HaloBottom)
		}

		// send halos to workers
		for i, workerClient := range workersClient {
			indexTop := ((i+1) % len(workersClient) + len(workersClient)) % len(workersClient)
			indexBottom := ((i-1) % len(workersClient) + len(workersClient)) % len(workersClient)


			request := stubs.Request{HaloBottom: halosTop[indexTop], HaloTop: halosBottom[indexBottom]}

			wg.Add(1)

			// asynchronously send rpc calls
			go func(i int, workerClient *rpc.Client, req stubs.Request) {
				defer wg.Done()

				workerResponses[i] = sendToWorker(workerClient, req, WokerHaloExchange)

			}(i, workerClient, request)

		}

		// if last turn then send to get world parts of each worker
		// this is done on the server (workers)

		turn ++
	}
	
	if terminate {
		res.World = world
		res.CompletedTurns = turn-1
		return
	}

	// reconstruct final world from all workers
	var newWorld [][]byte
	// wait for calls to complete and reconstruct the world
	wg.Wait()
	for _, workerResponse := range workerResponses {
		// wait for the request to complete (check for errors later)
		newWorld = append(newWorld, workerResponse.World...)
	}

	// copy the newWorld
	mutex.Lock()
	copy(world, newWorld)
	mutex.Unlock()

	// add final world to response
	res.World = world

	// alive cells
	res.AliveCells = util.CalculateAliveCells(world)

	res.CompletedTurns = turn

	fmt.Println("Sending final result to controller")

	return
}

// change
func (s *ControllerOperations) RequestAliveCellsCount(req stubs.Request, res *stubs.Response) (err error) {
	mutex.Lock()

	wg.Add(1)
	var currentWorld [][]byte
	for _, workerClient := range workersClient {
		response := sendToWorker(workerClient, stubs.Request{}, WorkerRequestCurrentGameState)
		currentWorld = append(currentWorld, response.World...)
	}
	wg.Done()

	res.AliveCellsCount = len(util.CalculateAliveCells(currentWorld))
	// reconstruct world from all workers
	res.CompletedTurns = turn-1
	mutex.Unlock()

	return
}

func (s *ControllerOperations) RequestCurrentGameState(req stubs.Request, res *stubs.Response) (err error) {
	mutex.Lock()
	wg.Add(1)
	var currentWorld [][]byte
	for _, workerClient := range workersClient {
		response := sendToWorker(workerClient, stubs.Request{}, WorkerRequestCurrentGameState)
		currentWorld = append(currentWorld, response.World...)
	}
	wg.Done()

	res.World = currentWorld
	res.CompletedTurns = turn-1
	mutex.Unlock()

	return
}

func (s *ControllerOperations) Shutdown(req stubs.Request, res *stubs.Response) (err error) {

	// server.Stop()

	wg.Add(1)
	// send to all worker to terminate
	for _, workerClient := range workersClient {
		// response := sendToWorker(workerClient, stubs.Request{}, WorkerShutdown)
		fmt.Println("sending shutdown to workers")
		sendToWorker(workerClient, stubs.Request{}, WorkerShutdown)
	}

	wg.Done()

	shutdown = true
	terminate = true

	// could add to return final state of world before closing

	return
}

func (s *ControllerOperations) TogglePause(req stubs.Request, res *stubs.Response) (err error) {
	mutex.Lock()

	// toggle running
	running = !running

	if !running {
		// busy waiting
		wg.Add(1)
	} else {
		// stop waiting
		wg.Done()
	}

	mutex.Unlock()

	// could add to return final state of world before closing
	res.IsPaused = !running
	res.CompletedTurns = turn-1

	return
}

func (s *ControllerOperations) StopWorkers(req stubs.Request, res *stubs.Response) (err error) {
	mutex.Lock()

	// terminate broker which stops sending to workers
	terminate = true

	mutex.Unlock()

	// wg.Add(1)
	// // send to all worker to terminate
	// for _, workerClient := range workersClient {
	// 	// response := sendToWorker(workerClient, stubs.Request{}, WorkerShutdown)
	// 	sendToWorker(workerClient, stubs.Request{}, WorkerStop)
	// }

	// wg.Done()


	// could add to return final state of world before closing
	// res.CompletedTurns = turn-1
	// res.World = world

	return
}


func sendToWorker(workerClient *rpc.Client, req stubs.Request, function string) *stubs.Response {

	res := new(stubs.Response)

	// call rpc function
	workerClient.Call(function, req, res)
	// workerClient.Go(function, req, res,)

	return res
}
