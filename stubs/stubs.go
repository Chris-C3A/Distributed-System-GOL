package stubs

import "uk.ac.bris.cs/gameoflife/util"

type ControllerOperations struct{}

type Response struct {
	World [][]byte
	AliveCells []util.Cell
	AliveCellsCount int
	CompletedTurns int
	HaloTop []byte
	HaloBottom []byte
}

type Request struct {
	World [][]byte
	Turns int
	HaloTop []byte
	HaloBottom []byte
	StartY int
	EndY int
	HaloWorkerAddr string
	NumOfWorkers int
	Start bool
}