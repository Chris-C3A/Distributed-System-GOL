package main

import (
	"fmt"
	"net/rpc"
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
	initialized = false
	mutex             sync.Mutex
	client *rpc.Client
	haloWorkerAddr string
	startNode = false
	HaloExchangeWorker = "WorkerOperations.HaloExchange"
)

var haloTopChan chan []byte
var haloBottomChan chan []byte
var haloTopToSend chan []byte
var haloBottomToSend chan []byte
var done chan bool

func (s *WorkerOperations) InitWorker(req stubs.Request, res *stubs.Response) (err error) {
	// fmt.Println("Worker initialized by broker")
	turn = 0
	turns = req.Turns
	world = req.World
	haloWorkerAddr = req.HaloWorkerAddr
	startNode = false

	if !initialized {
		// initialize channels
		haloTopChan = make(chan []byte, 1)
		haloBottomChan = make(chan []byte, 1)
		haloTopToSend = make(chan []byte, 1)
		haloBottomToSend = make(chan []byte, 1)
		done = make(chan bool)
	}

	if turns == 0 {
		res.World = world
		return
	}

	go worker(turns)

	client, _ = rpc.Dial("tcp", haloWorkerAddr)

	if req.Start {
		startNode = true
		startHaloExchange()
	}

	// wait for all turns to be processed before returning world
	<- done
	client.Close()

	// Wait to receive first halos
	res.World = world

	return
}

func (s *WorkerOperations) HaloExchange(req stubs.Request, res *stubs.Response) (err error) {
	fmt.Println("called")
	// Halo Exchange part

	haloTopChan <- req.HaloTop

	// stop at startNode
	if !startNode {
		// call right neighbour

		// send bottom as halo top to the next worker
		request := stubs.Request{HaloTop: <-haloBottomToSend, Start: false}

		response := new(stubs.Response)

		// maybe use .call
		// call := client.Go(HaloExchangeWorker, request, response, nil)
		mutex.Lock()
		client.Call(HaloExchangeWorker, request, response)
		// call := client.Go(HaloExchangeWorker, request, response, nil)
		mutex.Unlock()

		// <-call.Done

		// receive bottom halo
		haloBottomChan <- response.HaloBottom
	}


	// send back bottom halo
	res.HaloBottom = <-haloTopToSend
	// fmt.Println("sending back")
	fmt.Println("ended succesfully")

	return
}

func startHaloExchange() {
	fmt.Print("Starting halo exchange chain")
		// client, _ := rpc.Dial("tcp", haloWorkerAddr)

		// send bottom as halo top to the next worker
		// haloBottomToSend <- world[len(world)-1]
		request := stubs.Request{HaloTop: <-haloBottomToSend}

		response := new(stubs.Response)

		// maybe use .call
		// call := client.Go(HaloExchangeWorker, request, response, nil)

		mutex.Lock()
		client.Call(HaloExchangeWorker, request, response)
		// call := client.Go(HaloExchangeWorker, request, response, nil)
		mutex.Unlock()

		// fmt.Println("start node worked properly")

		// <-call.Done

		// receive bottom halo

		haloBottomChan <- response.HaloBottom
}



// Worker function
func worker(turns int) {
	for turn < turns && !terminate {

		// Send halos to other workers
		haloTopToSend <- world[0]
		haloBottomToSend <- world[len(world)-1]
		// fmt.Println("sent to channels properly")


		// Receive halos to include in the next iteration calculation
		haloTop := <-haloTopChan
		haloBottom := <-haloBottomChan

		mutex.Lock()
		world = util.CalculateNextState(world, haloTop, haloBottom)

		turn++
		mutex.Unlock()
		if startNode {
			// start haloexchange after every turn
			go startHaloExchange()
		}
		// fmt.Println("proccessed turn")
	}

	// Send done channel
	done <- true
}