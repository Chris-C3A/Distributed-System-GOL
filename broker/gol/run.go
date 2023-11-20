package gol

import (
	"flag"
	"fmt"
	"math/rand"
	"net/rpc"
	"strings"
	"time"
)


var server *Server
var workersClient []*rpc.Client

// Run starts the processing of Game of Life. It should initialise channels and goroutines.
func Run() {
	server = NewServer()

	pAddr := flag.String("port", "8030", "Port to listen on")
	workersAddr := flag.String("workerPorts", "8040", "List of worker node ports to connect to")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	// get workers addressses seperated by commas
	workers := strings.Split(*workersAddr, ",")
	// workers := []string{"8040", "8050", "8060", "8070"}

	// connect to workers

	for _, worker := range workers {
		// connect to worker
		// workerClient := connect(worker)

		client, err := rpc.Dial("tcp", "127.0.0.1:"+worker)
		fmt.Println("Connected to: ", worker)

		defer client.Close()

		// handle errors
		if err != nil {
			fmt.Println("Error connecting to worker:", err)
		}

		// add to array of worker clients
		workersClient = append(workersClient, client)
	}


	// start server
	if err := server.Start(*pAddr); err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	// // close clients
	// for _, workerClient := range workersClient {
	// 	defer workerClient.Close()
	// }


}