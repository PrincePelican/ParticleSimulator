package sandsimulation

const (
	screenWidth  = 1280
	screenHeight = 720
	gridWidth    = 640
	gridHeight   = 360
	cellSize     = 2
	clickSize    = 20
)

func ContinueFall(sandGrid [gridWidth][gridHeight]bool, gridInfo chan *[gridWidth][gridHeight]bool) {
	for y := len(sandGrid[0]) - 2; y >= 0; y-- {
		for x := 1; x < len(sandGrid)-1; x++ {
			if sandGrid[x][y] && !sandGrid[x][y+1] {
				sandGrid[x][y] = false
				sandGrid[x][y+1] = true
			} else if sandGrid[x][y] && sandGrid[x][y+1] {
				if !sandGrid[x-1][y+1] {
					sandGrid[x][y] = false
					sandGrid[x-1][y+1] = true
				} else if !sandGrid[x+1][y+1] {
					sandGrid[x][y] = false
					sandGrid[x+1][y+1] = true
				}
			}
		}
	}

	gridInfo <- &sandGrid
}

func CheckUnmovableSand(sandGrid [gridWidth][gridHeight]bool, yAxisToCheck int) int {
	for x := 1; x < len(sandGrid)-1; x++ {
		if !sandGrid[x][yAxisToCheck] {
			return yAxisToCheck
		}
	}
	return yAxisToCheck + 1
}
