// TODO:
// - Fix late collisions
// - Add portals, walls
// - Add multiple apples

package main

import (
	"image/color"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

type Game struct{}

type Orientation int

const (
	North Orientation = 0
	South Orientation = 1
	East  Orientation = 2
	West  Orientation = 3
)

const (
	screenX     = 720
	screenY     = 480
	tileSize    = 24.0
	titleString = "Snakebiten"
	playMsg     = "Press spacebar to start!"
)

var (
	moveEvery            = 15
	frameCounter         = 0
	flashFreq            = 5
	totalFlashes         = 60 / flashFreq
	score                = 0
	level                = 1
	snakeLen             = 1
	hasCollided          = false
	snakeVisible         = true
	maxXTiles, maxYTiles int
	board                [][]int
	appleCoord           Coord
	snakeCoords          DoublyLinkedList
	snakeOrientation     = East
	smallFont            font.Face
	bigFont              font.Face
	showMenu             = true
	nextSnakeOrientation = East
)

// Called before the program started
func init() {
	rand.Seed(time.Now().UnixNano())
	maxXTiles = (screenX - 2*tileSize) / tileSize
	maxYTiles = (screenY - 3*tileSize) / tileSize
	snakeCoords.PushFront(Coord{rand.Intn(maxXTiles), rand.Intn(maxYTiles)})
	appleCoord = snakeCoords.Front() // Makes the apple spawn on the head so we
	// begin with len2 snake
	tt, err := truetype.Parse(fonts.ArcadeN_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	smallFont = truetype.NewFace(tt, &truetype.Options{
		Size:    tileSize - 2,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	bigFont = truetype.NewFace(tt, &truetype.Options{
		Size:    tileSize * 2,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update(screen *ebiten.Image) error {
	if showMenu && ebiten.IsKeyPressed(ebiten.KeySpace) {
		showMenu = false
	}
	if showMenu {
		return nil
	}
	frameCounter++
	if !hasCollided {
		handleInput()
		if frameCounter == moveEvery {
			onApple := headOnApple()
			if onApple {
				appleCoord = Coord{rand.Intn(maxXTiles), rand.Intn(maxYTiles)}
				for !validApplePos() {
					appleCoord = Coord{rand.Intn(maxXTiles), rand.Intn(maxYTiles)}
				}
				score++
			}
			moveSnake(onApple)
			hasCollided = headCollision()
			frameCounter = 0
		}
	} else if totalFlashes > 0 {
		// Handle flashing
		if frameCounter%flashFreq == 0 {
			totalFlashes--
			if snakeVisible == true {
				snakeVisible = false
			} else {
				snakeVisible = true
			}
		}
		if totalFlashes == 0 {
			snakeVisible = true
		}
	} else if !snakeCoords.Front().IsNilCoord() {
		if frameCounter%flashFreq == 0 {
			if snakeVisible == true {
				snakeVisible = false
				snakeCoords.PopFront()
			} else {
				snakeVisible = true
			}
		}
	} else {
		// Restart game
		snakeCoords.PushFront(Coord{rand.Intn(maxXTiles), rand.Intn(maxYTiles)})
		appleCoord = snakeCoords.Front()
		score = 0
		level = 1
		snakeVisible = true
		hasCollided = false
		totalFlashes = 60 / flashFreq
		frameCounter = 0
	}
	return nil
}

// Gets the length of the string with the given font face
func textLen(text string, fontSize int) int {
	return len(text) * fontSize
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	drawBoard(screen)
	if !showMenu {
		drawApple(screen)
		if snakeVisible {
			drawSnake(screen)
		}
		text.Draw(screen, "Score:"+strconv.Itoa(score), smallFont,
			screenX-10*tileSize, screenY-(0.5*tileSize), color.White)
		text.Draw(screen, "Level:"+strconv.Itoa(level), smallFont,
			tileSize, screenY-(0.5*tileSize), color.White)

	} else {
		text.Draw(screen, titleString, bigFont,
			(screenX/2)-(textLen(titleString, tileSize*2)/2), (screenY-3*tileSize)/2, color.White)
		text.Draw(screen, playMsg, smallFont,
			(screenX/2)-(textLen(playMsg, tileSize-2)/2), (screenY)/2, color.White)
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenX, screenY
}

// Converts a standard coordinate pair to the coordinating tile-size screen coordinate
func coordToPixel(c Coord) (x float64, y float64) {
	newX := tileSize + float64(c.x)*tileSize
	newY := tileSize + float64(c.y)*tileSize
	return newX, newY
}

func drawBoard(screen *ebiten.Image) {
	// Draw the max size board given screenX, screenY, and tileSize
	// Outer border
	ebitenutil.DrawRect(screen, 0, 0, screenX, screenY, color.NRGBA{30, 90, 150, 255})
	ebitenutil.DrawRect(screen, tileSize, tileSize, screenX-(2*tileSize),
		screenY-(3*tileSize), color.Black)
}

func drawSnake(screen *ebiten.Image) {
	iter := snakeCoords.GetIterator()
	for iter.Next() {
		x, y := coordToPixel(iter.Get())
		ebitenutil.DrawRect(screen, x+1, y+1,
			tileSize-2, tileSize-2, color.White)
	}

}

func handleInput() {
	w := ebiten.IsKeyPressed(ebiten.KeyW)
	a := ebiten.IsKeyPressed(ebiten.KeyA)
	s := ebiten.IsKeyPressed(ebiten.KeyS)
	d := ebiten.IsKeyPressed(ebiten.KeyD)
	if w && snakeOrientation != South {
		nextSnakeOrientation = North
	}
	if a && snakeOrientation != East {
		nextSnakeOrientation = West
	}
	if s && snakeOrientation != North {
		nextSnakeOrientation = South
	}
	if d && snakeOrientation != West {
		nextSnakeOrientation = East
	}
}

func moveSnake(onApple bool) {
	snakeOrientation = nextSnakeOrientation
	curCoord := snakeCoords.Front()
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
	snakeCoords.PushFront(curCoord)
	// Remove the last item if we are not on an apple
	if !onApple {
		snakeCoords.PopBack()
	} else {
		snakeLen++
	}

}

func headOnApple() bool {
	if Equals(appleCoord, snakeCoords.Front()) {
		return true
	}
	return false
}

func drawApple(screen *ebiten.Image) {
	x, y := coordToPixel(appleCoord)
	ebitenutil.DrawRect(screen, x+1, y+1, tileSize-2, tileSize-2, color.NRGBA{255, 0, 0, 255})
}

// Checks if the head is colliding with the wall OR another component of the snake
func headCollision() bool {
	head := snakeCoords.Front()
	// Check if out of bounds
	if head.x >= maxXTiles || head.y >= maxYTiles || head.x < 0 || head.y < 0 {
		return true
	}
	iter := snakeCoords.GetIterator()
	// Bypass the head
	iter.Next()
	for iter.Next() {
		if Equals(iter.Get(), head) {
			return true
		}
	}
	return false
}

// Checks if the apple coord is inside of the snake
func validApplePos() bool {
	iter := snakeCoords.GetIterator()
	for iter.Next() {
		if Equals(appleCoord, iter.Get()) {
			return false
		}
	}
	return true
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
