package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)


var server *Server

func main() {
	// create new server instance
	server = NewServer()

	pAddr := flag.String("port", "8040", "Port to listen on")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	// start server
	if err := server.Start(*pAddr); err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
}
