package gol

import (
	"fmt"
	"time"

	"uk.ac.bris.cs/gameoflife/stubs"
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

// func ticker(events chan<- Event) {
// 	for {
// 		// sleep for 2 seconds
// 		time.Sleep(2 * time.Second)

// 		// send rpc call
// 		// fmt.Println("sending request alive cells")
// 	}
// }

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels) {
	// send input command
	c.ioCommand <- ioInput
	c.ioFilename <- fmt.Sprintf("%dx%d", p.ImageWidth, p.ImageHeight)

	// initialise world
	world := make([][]byte, p.ImageHeight)
	for i := range world {
		world[i] = make([]byte, p.ImageWidth)
	}

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


	// TODO return later
	ticker := time.NewTicker(2 * time.Second)
	done := make(chan bool)

	// anonymous go routine to handle keypress and tickers
	go func() {
			for {
					select {
					case <-done:
							return
					case <-ticker.C:
						res := sendToRPC(stubs.Request{}, RequestAliveCellsCount)
						c.events <- AliveCellsCount{CompletedTurns: res.CompletedTurns, CellsCount: res.AliveCellsCount}
					// TODO continue
					case key := <-c.keyPresses:
						if key == 's' {
							// get current state of the board then outputPGM file
							response := sendToRPC(stubs.Request{}, RequestCurrentGameState)

							outputPGMFile(p, c, response.CompletedTurns, response.World)
						} else if key == 'q' {
							// terminate controller without causing error on server

							// send termination to server to get the last state then close
							client.Close()
						} else if key == 'k' {
							// send rpc to cleanly kill components and return last state of the game to ouput

							// all componenets of the distributed system is shutdown cleanly and the system outputs a pgm image of the latest state
							sendToRPC(stubs.Request{}, Shutdown)

						} else if key == 'p' {
							// pause the process on the aws node and have the controller print the current turn
						}
					}
			}
	}()

	fmt.Println("sending to broker")

	// send evolveGOL to broker to handle
	response := sendToRPC(stubs.Request{World: world, Turns: p.Turns}, EvolveGoL)
	fmt.Println("received response")

	// TODO reutrn later
	ticker.Stop()
	done <- true

	fmt.Println("received response")

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
	outFileName := fmt.Sprintf("%dx%dx%d", p.ImageWidth, p.ImageHeight, p.Turns)
	fmt.Println(outFileName)
	// send to channel
	c.ioFilename <- outFileName

	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	// send imageoutput complete event
	c.events <- ImageOutputComplete{CompletedTurns: p.Turns, Filename: outFileName}

	c.events <- StateChange{p.Turns, Quitting}

	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}

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
