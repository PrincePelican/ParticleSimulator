package sandsimulation

const (
	screenWidth  = 1280
	screenHeight = 720
	gridWidth    = 320
	gridHeight   = 180
	cellSize     = 4
	clickSize    = 10
)

func ContinueFall(sandGrid [gridWidth][gridHeight]bool, gridInfo chan [gridWidth][gridHeight]bool) {
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
	gridInfo <- sandGrid
}
