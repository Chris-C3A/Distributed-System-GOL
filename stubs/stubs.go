package stubs

import "uk.ac.bris.cs/gameoflife/util"

type ControllerOperations struct{}

type Response struct {
	World [][]byte
	AliveCells []util.Cell
	AliveCellsCount int
	CompletedTurns int
}

type Request struct {
	World [][]byte
	Turns int
	StartY int
	EndY int
	NumOfWorkers int
}