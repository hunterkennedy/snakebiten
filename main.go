package main

import (
	"image/color"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Game struct{}
type Coord struct {
	x int
	y int
}
type Orientation int
const (
	North Orientation = 0
	South Orientation = 1
	East  Orientation = 2
	West  Orientation = 3
)

const (
	screenX = 720
	screenY = 480
	tileSize = 20.0
)
var (
	moveEvery = 30
	counter = 0
	score = 0
	level = 1
	maxXTiles, maxYTiles int
	board [][]int
	appleX, appleY int
	snakeCoords []Coord
	snakeOrientation Orientation = East
	snakeLen = 1
)

// Called before the program started
func init() {
	maxXTiles = (screenX - 2 * tileSize) / tileSize
	maxYTiles = (screenY - 2 * tileSize) / tileSize
	snakeCoords = append(snakeCoords, Coord{maxXTiles/2, maxYTiles/2})
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update(screen *ebiten.Image) error {
	counter++
	handleInput()
	if (counter == moveEvery) {
		moveSnake(headOnApple())
		counter = 0
	}
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "Score: " + strconv.Itoa(score),
		screenX - 4 * tileSize, screenY - tileSize)
	ebitenutil.DebugPrintAt(screen, "Level: " + strconv.Itoa(level),
		tileSize, screenY - tileSize)
	//drawBoard(screen)
	// drawApple()
	drawSnake(screen)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenX, screenY
}

// Converts a standard coordinate pair to the coordinating tile-size screen coordinate
func coordToPixel(c Coord) (x float64, y float64) {
	newX := tileSize + float64(c.x) * tileSize
	newY := tileSize + float64(c.y) * tileSize
	return newX, newY
}

func drawBoard(screen *ebiten.Image) {
	// Draw the max size board given screenX, screenY, and tileSize
	ebitenutil.DrawRect(screen, 0, 0, screenX, screenY, color.White)
	ebitenutil.DrawRect(screen, tileSize, tileSize, screenX-(2*tileSize), 
		screenY-(2*tileSize), color.Black)
}

func drawSnake(screen *ebiten.Image) {
	for _, xy := range snakeCoords {
		x, y:= coordToPixel(xy)
		ebitenutil.DrawRect(screen, x+2, y+2, tileSize-2, tileSize-2, color.White)
	}
	
}

func handleInput() {
	w := ebiten.IsKeyPressed(ebiten.KeyW)
	a := ebiten.IsKeyPressed(ebiten.KeyA)
	s := ebiten.IsKeyPressed(ebiten.KeyS)
	d := ebiten.IsKeyPressed(ebiten.KeyD)
	if w {
		snakeOrientation = North
	}
	if a {
		snakeOrientation = West
	}
	if s {
		snakeOrientation = South
	}
	if d {
		snakeOrientation = East
	}
}

func moveSnake(onApple bool) {
	
	curCoord := snakeCoords[0]
	if snakeOrientation == North {
		curCoord.y -= 1
	}
	if snakeOrientation == West {
		curCoord.x -= 1
	}
	if snakeOrientation == South {
		curCoord.y += 1
	}
	if snakeOrientation == East {
		curCoord.x += 1
	}
	snakeCoords = append(snakeCoords, Coord{})
	copy(snakeCoords[1:], snakeCoords[0:])
	snakeCoords[0] = curCoord
	if !onApple {
		snakeCoords = snakeCoords[:snakeLen]
	} else {
		snakeLen++
	}
	
}

func headOnApple() bool {
	if snakeCoords[0].x == appleX && snakeCoords[0].y == appleY {
		return true
	}
	return false
}

func main() {
	// ebiten.SetMaxTPS(10)
	game := &Game{}
	ebiten.SetWindowSize(screenX, screenY)
	ebiten.SetWindowTitle("Snakebiten")
	// ebiten.SetWindowIcon()
	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}