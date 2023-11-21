package util

func MakeWorld(width, height int) [][]byte {
	world := make([][]byte, height)

	for i := range world {
		world[i] = make([]byte, width)
	}

	return world
}

func CalculateNextState(world [][]byte, startY, endY int) [][]byte {
	height := endY - startY
	width := len(world[0])

	newWorld := MakeWorld(width, height)

	for i := startY; i < endY; i++ {
		for j := 0; j < width; j++ {
			// get number of live neighbours
			numOfLiveNeighbours := getNumOfLiveNeighbours(world, i, j)

			// rules of the game of life
			if world[i][j] == ALIVE {
				if numOfLiveNeighbours < 2 || numOfLiveNeighbours > 3 {
					newWorld[i-startY][j] = DEAD
				} else {
					newWorld[i-startY][j] = ALIVE
				}
			} else {
				if numOfLiveNeighbours == 3 {
					newWorld[i-startY][j] = ALIVE
				}
			}
		}
	}

	return newWorld

}

func CalculateAliveCells(world [][]byte) []Cell {
	var aliveCells []Cell

	for y := 0; y < len(world); y++ {
		for x := 0; x < len(world[y]); x++ {
			if world[y][x] == ALIVE {
				// add cell coordinates to aliveCells slice
				aliveCells = append(aliveCells, Cell{X: x, Y: y})
			}
		}
	}

	return aliveCells
}

// helper functions for gol logic
func getNumOfLiveNeighbours(world [][]byte, i int, j int) int {
	numOfLiveNeighbours := 0

	// positive modulus
	up := ((i-1)%len(world) + len(world)) % len(world)
	down := ((i+1)%len(world) + len(world)) % len(world)
	right := ((j+1)%len(world[i]) + len(world[i])) % len(world[i])
	left := ((j-1)%len(world[i]) + len(world[i])) % len(world[i])

	// neighbours of cell i,j
	neighbours := [8]byte{world[up][j], world[down][j], world[i][left], world[i][right], world[up][left], world[up][right], world[down][right], world[down][left]}

	for _, neighbour := range neighbours {
		if neighbour == ALIVE {
			numOfLiveNeighbours++
		}
	}

	return numOfLiveNeighbours

}
