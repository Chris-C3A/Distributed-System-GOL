package gol

import (
	"fmt"
	"time"

	"uk.ac.bris.cs/gameoflife/stubs"
	"uk.ac.bris.cs/gameoflife/util"
)

type distributorChannels struct {
	events     chan<- Event
	ioCommand  chan<- ioCommand
	ioIdle     <-chan bool
	ioFilename chan<- string
	ioOutput   chan<- uint8
	ioInput    <-chan uint8
	keyPresses <-chan rune
}


// rpc method to call
var InitializeBroker = "ControllerOperations.Broker"
var EvolveGoL = "ControllerOperations.EvolveGoL"
var RequestAliveCellsCount = "ControllerOperations.RequestAliveCellsCount"
var RequestCurrentGameState = "ControllerOperations.RequestCurrentGameState"
var Shutdown = "ControllerOperations.Shutdown"
var TogglePause = "ControllerOperations.TogglePause"
var StopWorkers = "ControllerOperations.StopWorkers"


// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels) {
	// send input command
	c.ioCommand <- ioInput
	// send filename of pgm file to read
	c.ioFilename <- fmt.Sprintf("%dx%d", p.ImageWidth, p.ImageHeight)

	// initialise world
	world := util.MakeWorld(p.ImageWidth, p.ImageHeight)

	// store image bytes into world
	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			world[y][x] = <-c.ioInput
		}
	}

	// send cell flipped events for initial world
	// for _, cell := range calculateAliveCells(p, world) {
	// 	c.events <- CellFlipped{CompletedTurns: 0, Cell: cell}
	// }

	// to render blank sdl window
	c.events <- TurnComplete{CompletedTurns: 0}

	// after every turn send the state of the board to run cellFlipped
	// only send alive cells to run cellFlipped
	// or some kind of that logic

	// initialise the 2 second ticker
	ticker := time.NewTicker(2 * time.Second)
	done := make(chan bool)

	// anonymous go routine to handle keypress and tickers
	go func() {
			for {
					select {
					case <-done:
							return
					case <-ticker.C:
						// request alive cells from broker
						res := sendToRPC(stubs.Request{}, RequestAliveCellsCount)

						c.events <- AliveCellsCount{CompletedTurns: res.CompletedTurns, CellsCount: res.AliveCellsCount}
					case key := <-c.keyPresses:
						// TODO continue
						if key == 's' {
							// get current state of the board then outputPGM file
							response := sendToRPC(stubs.Request{}, RequestCurrentGameState)
							outputPGMFile(p, c, response.CompletedTurns, response.World)

						} else if key == 'q' {
							fmt.Println("q key press entered")
							// terminate controller without causing error on server
							// causes error on controller

							// TODO
							// send rpc that controller is closing to stop execution of workers
							sendToRPC(stubs.Request{}, StopWorkers)

							// send termination to server to get the last state then close
							// client.Close()
						} else if key == 'k' {
							// send rpc to cleanly kill components and return last state of the game to ouput
							// TODO
							fmt.Println("killing components")

							// all componenets of the distributed system is shutdown cleanly and the system outputs a pgm image of the latest state
							sendToRPC(stubs.Request{}, Shutdown)

						} else if key == 'p' {
							// pause the process on the aws node and have the controller print the current turn
							// send rpc to broker to toggle pause functionality
							response := sendToRPC(stubs.Request{}, TogglePause)

							if response.IsPaused {
								fmt.Println("Paused at ", response.CompletedTurns)
							} else {
								fmt.Println("Continuing")
							}

						}
					}
			}
	}()

	fmt.Println("sending to broker")

	// send evolveGOL to broker to handle
	// with number of workers to use
	// if its more than the available workers it will use al the available ones
	response := sendToRPC(stubs.Request{World: world, Turns: p.Turns, NumOfWorkers: p.Threads}, EvolveGoL)

	fmt.Println("received response")

	// stop ticker and send on done channel
	ticker.Stop()
	done <- true

	// FinalTurnComplete Event
	c.events <- FinalTurnComplete{CompletedTurns: p.Turns, Alive: response.AliveCells}

	// put each byte of final world into output channel
	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			c.ioOutput <- response.World[y][x]
		}
	}

	// Make sure that the Io has finished any output before exiting.

	// put ioOutput into command channel
	c.ioCommand <- ioOutput
	// output file
	outFileName := fmt.Sprintf("%dx%dx%d", p.ImageWidth, p.ImageHeight, response.CompletedTurns)
	fmt.Println(outFileName)

	// send to channel
	c.ioFilename <- outFileName

	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	// send imageoutput complete event
	c.events <- ImageOutputComplete{CompletedTurns: response.CompletedTurns, Filename: outFileName}

	c.events <- StateChange{response.CompletedTurns, Quitting}

	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}

// helper functions
func sendToRPC(req stubs.Request, function string) *stubs.Response {

	res := new(stubs.Response)

	// call rpc function
	client.Call(function, req, res)

	return res
}

func outputPGMFile(p Params, c distributorChannels, turn int, world [][]byte) {
	// put each byte of final world into output channel
	for y := 0; y < p.ImageHeight; y++ {
		for x := 0; x < p.ImageWidth; x++ {
			c.ioOutput <- world[y][x]
		}
	}

	// put ioOutput into command channel
	c.ioCommand <- ioOutput
	// output file
	outFileName := fmt.Sprintf("%dx%dx%d", p.ImageWidth, p.ImageHeight, turn)
	// send to channel
	c.ioFilename <- outFileName

	// send image output complete event
	c.events <- ImageOutputComplete{CompletedTurns: turn, Filename: outFileName}
}
