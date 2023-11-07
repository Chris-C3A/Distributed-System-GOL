package gol

import (
	"flag"
	"math/rand"
	"net"
	"net/rpc"
	"time"
)

type distributorChannels struct {
	events     chan<- Event
	ioCommand  chan<- ioCommand
	ioIdle     <-chan bool
	ioFilename chan<- string
	ioOutput   chan<- uint8
	ioInput    <-chan uint8
}

func (s *ControllerOperations) EvolveGoL(req Request, res *Response) (err error) {
	res.Message = s.EvolveGoL(req.Message, "omak")
	return
}

type Response struct {
	Message string
}

type Request struct {
	Message string
}

type ControllerOperations struct{}

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels) {

	pAddr := flag.String("port", "8030", "Port to listen on")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	rpc.Register(&ControllerOperations{})
	ln, _ := net.Listen("tcp", ":"+*pAddr)
	defer ln.Close()
	rpc.Accept(ln)
	turn := 0

	// Make sure that the Io has finished any output before exiting.
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	c.events <- StateChange{turn, Quitting}

	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}
