package main

import (
	"fmt"

	"math/rand"

	"os"

	"os/exec"

	"runtime"

	"time"

	"image/color"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	WIDTH        = kingpin.Flag("width", "The width of the grid").Short('w').Default("80").Int()
	HEIGHT       = kingpin.Flag("height", "The height of the grid").Short('h').Default("40").Int()
	ITERATIONS   = kingpin.Flag("iterations", "Number of iterations. Any negative number will use the default, infinity").Short('i').Default("-1").Int()
	FPS          = kingpin.Flag("fps", "Frames per second, how log to wait until the next iteration is displayed").Short('f').Default("10").Int()
	PERCENTAGE   = kingpin.Flag("percentage", "Percentage of living cells at the start").Short('p').Default("33").Int()
	INITIALSTATE = kingpin.Flag("initial-state", "Initial starting pattern for cells, [Random, VerticalBlinker, DiagonalBlinker]").Short('s').Default("Random").String()
)

type Game struct {
	generation   int
	grid, buffer [][]uint8
	xSize, ySize int
}

func NewMatrix(rows int, cols int) [][]uint8 {
	m := make([][]uint8, rows)
	for r := range m {
		m[r] = make([]uint8, cols)
	}
	return m
}

var (
	counter int        = 0
	black   color.RGBA = color.RGBA{95, 95, 95, 255}
	white   color.RGBA = color.RGBA{233, 233, 233, 255}
)

func initialiseGameBoard(height, width int) *Game {
	grid := NewMatrix(height, width)
	buffer := NewMatrix(height, width)
	return &Game{generation: 0, xSize: width, ySize: height, grid: grid, buffer: buffer}
}

func (g *Game) SetInitialState(state string, livingPercentage int) bool {
	if state == "Random" || state == "random" {
		g.RandInit(livingPercentage)
	} else if state == "VerticalBlinker" || state == "verticalblinker" {
		g.VerticalGlider()
	} else if state == "DiagonalBlinker" || state == "diagonalblinker" {
		g.DiagonalGlider()
	} else {
		return false
	}
	return true
}

func (g *Game) RandInit(livingPercentage int) {

	// Number of living cells
	numLiving := livingPercentage * g.xSize * g.ySize / 100
	count := 0
	r := rand.New(rand.NewSource(time.Now().Unix()))
	// This has the (Extremely small) potential to be very slow if users start supplying high percentages
	for count < numLiving {
		// Generate a random index into string
		xRand := r.Intn(g.xSize)
		yRand := r.Intn(g.ySize)
		if g.grid[yRand][xRand] != 1 {
			g.grid[yRand][xRand] = 1
			count++
		}
	}

	fmt.Printf("Total Cells: %d\nNumber of Living cells: %d\n", g.xSize*g.ySize, numLiving)
}

func (g *Game) Blinker() {
	g.grid[9][20] = 1
	g.grid[10][20] = 1
	g.grid[11][20] = 1
}

func (g *Game) DiagonalGlider() {
	// grid[WIDTH-1][HEIGHT-1] = 1
	g.grid[10][13] = 1
	g.grid[12][13] = 1
	g.grid[12][12] = 1
	g.grid[12][14] = 1
	g.grid[11][14] = 1
}

func (g *Game) Iterate() {
	for y := 1; y < g.ySize-2; y++ {
		for x := 1; x < g.xSize-2; x++ {
			// Because we are only using two copies of the Matrix we need to make sure we clear the current buffer
			g.buffer[y][x] = 0

			n := g.grid[y-1][x-1] + g.grid[y-1][x] + g.grid[y-1][x+1] + g.grid[y][x-1] + g.grid[y][x+1] + g.grid[y+1][x-1] + g.grid[y+1][x] + g.grid[y+1][x+1]

			if g.grid[y][x] == 0 && n == 3 {
				g.buffer[y][x] = 1
			} else if n > 3 || n < 2 {
				g.buffer[y][x] = 0
			} else {
				g.buffer[y][x] = g.grid[y][x]
			}
		}
	}

	// TOP
	// for x := 1; x < g.xSize-2; x++ {
	// 	g.buffer[0][x] = 0

	// 	n := g.grid[0][x-1] + g.grid[0][x+1] + g.grid[1][x-1] + g.grid[1][x] + g.grid[1][x+1] + g.grid[g.ySize-1][x-1] + g.grid[g.ySize-1][x] + g.grid[g.ySize-1][x+1]

	// 	// Now implement the rules of life
	// 	if g.grid[0][x] == 0 && n == 3 {
	// 		g.buffer[0][x] = 1
	// 	} else if n > 3 || n < 2 {
	// 		g.buffer[0][x] = 0
	// 	} else {
	// 		g.buffer[0][x] = g.grid[0][x]
	// 	}
	// }

	// // BOTTOM
	// for x := 1; x < g.xSize-2; x++ {
	// 	g.buffer[g.ySize-1][x] = 0

	// 	n := g.grid[g.ySize-1][x-1] + g.grid[g.ySize-1][x+1] + g.grid[0][x-1] + g.grid[0][x] + g.grid[0][x+1] + g.grid[g.ySize-2][x-1] + g.grid[g.ySize-2][x] + g.grid[g.ySize-2][x+1]

	// 	// Now implement the rules of life
	// 	if g.grid[g.ySize-1][x] == 0 && n == 3 {
	// 		g.buffer[g.ySize-1][x] = 1
	// 	} else if n > 3 || n < 2 {
	// 		g.buffer[g.ySize-1][x] = 0
	// 	} else {
	// 		g.buffer[g.ySize-1][x] = g.grid[g.ySize-1][x]
	// 	}
	// }

	temp := g.buffer
	g.buffer = g.grid
	g.grid = temp
}

func (g *Game) PrintGame() {
	//Top margin
	fmt.Print("╔")
	for x := 1; x <= g.xSize; x++ {
		fmt.Print("══")
	}
	fmt.Println("╗")

	//Rows
	for y := 0; y < g.ySize; y++ {
		fmt.Print("║")
		//Collumns
		for x := 0; x < g.xSize; x++ {
			if g.grid[y][x] == 1 {
				fmt.Print("██")
			} else {
				fmt.Print("  ")
			}
		}
		fmt.Println("║")
	}

	//Bottom margin
	fmt.Print("╚")
	for x := 1; x <= g.xSize; x++ {
		fmt.Print("══")
	}
	fmt.Println("╝")
}

func (g *Game) VerticalGlider() {
	g.grid[5][19] = 1
	g.grid[5][18] = 1
	g.grid[5][17] = 1
	g.grid[6][17] = 1
	g.grid[6][20] = 1
	g.grid[7][17] = 1
	g.grid[8][17] = 1
	g.grid[9][20] = 1
	g.grid[9][18] = 1
}

func add(x *int) {
	*x += 7
}

func main() {
	kingpin.Version("1.0.0")
	kingpin.Parse()

	// Extract Input Variables
	i := *ITERATIONS
	sleepTime := time.Duration(1000 / *FPS) * time.Millisecond
	game := initialiseGameBoard(*HEIGHT, *WIDTH)
	if game.SetInitialState(*INITIALSTATE, *PERCENTAGE) {
		// Main Game Loop
		for i > 0 {
			i--

			cmd := exec.Command("clear")
			if runtime.GOOS == "windows" {
				cmd = exec.Command("cmd", "/c", "cls")
			}
			cmd.Stdout = os.Stdout
			cmd.Run()
			game.PrintGame()
			fmt.Printf("Generation: %d\n", 0+*ITERATIONS-i)
			game.Iterate()
			time.Sleep(sleepTime)
		}
	} else {
		fmt.Println("Incorrect Initial State Provided...\nPlease run with the --help flag for options.")
	}
}
