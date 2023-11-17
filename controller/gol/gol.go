package gol

import (
	"fmt"
	"net/rpc"
)

// Params provides the details of how to run the Game of Life and which image to load.
type Params struct {
	Turns       int
	Threads     int
	ImageWidth  int
	ImageHeight int
}

var client *rpc.Client

// Run starts the processing of Game of Life. It should initialise channels and goroutines.
func Run(p Params, events chan<- Event, keyPresses <-chan rune) {

	server := "127.0.0.1:8030"
	// connect to server
	fmt.Println("Connected to: ", server)

	var err error
	client, err = rpc.Dial("tcp", server)

	defer client.Close()

	// handle errors
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return 
	}
	
	ioCommand := make(chan ioCommand)
	ioIdle := make(chan bool)
	ioFilename := make(chan string)
	ioOutput := make(chan uint8, p.ImageHeight * p.ImageWidth)
	ioInput := make(chan uint8, p.ImageHeight * p.ImageWidth)

	ioChannels := ioChannels{
		command:  ioCommand,
		idle:     ioIdle,
		filename: ioFilename,
		output:   ioOutput,
		input:    ioInput,
	}
	go startIo(p, ioChannels)

	distributorChannels := distributorChannels{
		events:     events,
		ioCommand:  ioCommand,
		ioIdle:     ioIdle,
		ioFilename: ioFilename,
		ioOutput:   ioOutput,
		ioInput:    ioInput,
		keyPresses: keyPresses,
	}
	distributor(p, distributorChannels)
}

// func connect(server string) {

// }