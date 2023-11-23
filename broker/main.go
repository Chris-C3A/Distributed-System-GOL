package main

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
var workersAddr []string

func main() {
	// create new server instance
	server = NewServer()

	pAddr := flag.String("port", "8030", "Port to listen on")
	
	workersAddrP := flag.String("workersAddr", "127.0.0.1:8040", "List of worker node addreses and ports to connect to")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	// get workers addressses seperated by commas
	workersAddr = strings.Split(*workersAddrP, ",")
	// workers := []string{"8040", "8050", "8060", "8070"}

	// connect to workers
	for _, worker := range workersAddr {
		// connect to worker
		// workerClient := connect(worker)

		client, err := rpc.Dial("tcp", worker)
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

}