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

// Game defined for Ebiten; No modifications necessary
type Game struct{}

// Orientation is an enum-ish type for NSEW
type Orientation int

const (
	// North = Up
	North Orientation = 0
	// South = Down
	South Orientation = 1
	// East = Right
	East Orientation = 2
	// West = Left
	West Orientation = 3
)

const (
	screenX     = 720
	screenY     = 480
	tileSize    = 24.0
	titleString = "Snakebiten"
	playMsg     = "Press spacebar to start!"
	controls    = "Use WASD to change direction"
	moveEvery   = 20
	flashFreq   = 3
)

var (
	highScore            = 0
	frameCounter         = 0
	totalFlashes         = 60 / flashFreq
	score                = 0
	hasCollided          = false
	snakeVisible         = true
	maxXTiles, maxYTiles int
	appleCoord           Coord
	snakeCoords          DoublyLinkedList
	snakeOrientation     = East
	smallFont            font.Face
	bigFont              font.Face
	showMenu             = true // Is menu on the screen?
	nextSnakeOrientation = East
	promptVisible        = true // Is the hlper text visible? (For flashing)
	maxScore             int
)

// Called before the program starts; set up initial variables
func init() {
	rand.Seed(time.Now().UnixNano())
	maxXTiles = (screenX - 2*tileSize) / tileSize
	maxYTiles = (screenY - 3*tileSize) / tileSize
	maxScore = maxXTiles * maxYTiles
	c, d := genSnakeCoordDir()
	snakeCoords.PushFront(c)
	nextSnakeOrientation = d
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
	frameCounter++
	// Logic to hide menu and start the game
	if showMenu && ebiten.IsKeyPressed(ebiten.KeySpace) {
		showMenu = false
		snakeVisible = true
		promptVisible = false
		frameCounter = 0
	}
	// Logic for flashing the "press space to play" prompt
	if showMenu {
		if frameCounter%(8*flashFreq) == 0 {
			if promptVisible {
				promptVisible = false
			} else {
				promptVisible = true
			}
		}
		// If the menu is shown don't do any game logic
		return nil
	}
	// Primary game logic
	if !hasCollided {
		// Get input and enqueue the next move
		handleInput()
		if frameCounter == moveEvery {
			onApple := headOnApple()
			if onApple {
				// Need to generate a valid apple position
				appleCoord = Coord{rand.Intn(maxXTiles), rand.Intn(maxYTiles)}
				for !validApplePos() {
					appleCoord = Coord{rand.Intn(maxXTiles), rand.Intn(maxYTiles)}
				}
				// Hit an apple so increment the score
				score++
			}
			next := getNextSnakePosition()
			if coordCollides(next) {
				// Collided!
				hasCollided = true
			} else {
				// Move the snake to the next position
				snakeCoords.PushFront(next)
				// Remove the last item if we are not on an apple
				if !onApple {
					snakeCoords.PopBack()
				}
			}
			// Reset frame counter, just to keep this number small
			frameCounter = 0
		}
	} else if totalFlashes > 0 {
		// We know we have collided, and that we still need to flash before
		// dismantling the snake
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
		// Handles dismantling of the snake
		if frameCounter%flashFreq == 0 {
			if snakeVisible == true {
				snakeVisible = false
				snakeCoords.PopFront()
			} else {
				snakeVisible = true
			}
		}
	} else {
		// Restart game, reset running vars
		c, d := genSnakeCoordDir()
		snakeCoords.PushFront(c)
		nextSnakeOrientation = d
		appleCoord = snakeCoords.Front()
		highScore = score
		score = 0
		snakeVisible = true
		hasCollided = false
		totalFlashes = 60 / flashFreq
		frameCounter = 0
		showMenu = true
		promptVisible = true
	}
	return nil
}

// Gets the length of the string with the given font face size
// I don't understand the ebiten implementation of this function so I
// wrote my own
func textLen(text string, fontSize int) int {
	return len(text) * fontSize
}

// Generates a random starting position and direction for the snake
// This avoid the issue of being spawned against and facing a wall
func genSnakeCoordDir() (Coord, Orientation) {
	d := West
	c := Coord{rand.Intn(maxXTiles), rand.Intn(maxYTiles)}
	if c.x < maxXTiles/2 {
		d = East
	}
	return c, d
}

// Converts a diretion to a string for the status message
// First place I realized I could use a switch
func dirToString(o Orientation) string {
	switch o {
	case North:
		return "UP"
	case East:
		return "RIGHT"
	case South:
		return "DOWN"
	default:
		return "LEFT"
	}
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the board and outline
	drawBoard(screen)
	if !showMenu {
		// Game is running, perform usual actions:
		drawApple(screen)
		if snakeVisible {
			// Must do this check for when the snake flashes
			drawSnake(screen)
		}
		scoreString := "Score:" + strconv.Itoa(score)
		scoreXOffset := screenX - (tileSize + textLen(scoreString, tileSize-2))
		text.Draw(screen, scoreString, smallFont,
			scoreXOffset,
			screenY-(0.5*tileSize), color.White)
		// Leaving this in in case I decide to implement proper levels,
		// portals, etc.
		// text.Draw(screen, levelString, smallFont,
		// 	scoreXOffset-(2*tileSize+textLen(levelString, tileSize-2)),
		// 	screenY-(0.5*tileSize), color.White)
		text.Draw(screen, dirToString(nextSnakeOrientation), smallFont,
			tileSize,
			screenY-(tileSize*0.5), color.White)
	} else {
		// Show the menu, title, prompt, etcc
		text.Draw(screen, titleString, bigFont,
			(screenX/2)-(textLen(titleString, tileSize*2)/2), (screenY-3*tileSize)/2, color.White)
		if promptVisible {
			text.Draw(screen, playMsg, smallFont,
				(screenX/2)-(textLen(playMsg, tileSize-2)/2), screenY/2, color.White)
		}
		if highScore != 0 {
			// Show the high score if we have a high score recorded
			highScoreString := "High-Score: " + strconv.Itoa(highScore)
			text.Draw(screen, highScoreString, smallFont,
				(screenX/2)-(textLen(highScoreString, tileSize-2)/2),
				(screenY/2)+3*tileSize, color.White)
		} else {
			// Show the controls iff they have not played yet
			text.Draw(screen, controls, smallFont,
				(screenX/2)-(textLen(controls, tileSize-2)/2),
				(screenY/2)+4*tileSize, color.White)
		}
		if highScore >= maxScore {
			// Probably a little underwhelming for beating a game of snake
			congrats := "Congratulations!"
			beat := "You beat the game!"
			text.Draw(screen, congrats, smallFont,
				(screenX/2)-(textLen(congrats, tileSize-2)/2),
				(screenY/2)+5*tileSize, color.White)
			text.Draw(screen, beat, smallFont,
				(screenX/2)-(textLen(beat, tileSize-2)/2),
				(screenY/2)+6*tileSize, color.White)
		}
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// Const for now, makes it easier to work with
	return screenX, screenY
}

// Converts a standard coordinate pair to the coordinating
// tile-size screen coordinate
func coordToPixel(c Coord) (x float64, y float64) {
	newX := tileSize + float64(c.x)*tileSize
	newY := tileSize + float64(c.y)*tileSize
	return newX, newY
}

// Draws the blue border and black board within
func drawBoard(screen *ebiten.Image) {
	// Draw the max size board given screenX, screenY, and tileSize
	// Outer border
	ebitenutil.DrawRect(screen, 0, 0, screenX, screenY, color.NRGBA{30, 90, 150, 255})
	// Inner board
	ebitenutil.DrawRect(screen, tileSize, tileSize, screenX-(2*tileSize),
		screenY-(3*tileSize), color.Black)
}

// Converts i to a uint8 less than or equal to 255, capped at 255
func getColor(i int) uint8 {
	ret := i
	if ret > 255 {
		return uint8(255)
	}
	return uint8(ret)
}

// Draws the snake with a nice gradient, green to white
// This was intended to make the snake easier to see and
// easier to distinguish parts of the snake, but probably
// needs more colors to properly do that
func drawSnake(screen *ebiten.Image) {
	iter := snakeCoords.GetIterator()
	i := 0
	for iter.Next() {
		x, y := coordToPixel(iter.Get())
		ebitenutil.DrawRect(screen, x+1, y+1,
			tileSize-2, tileSize-2, color.NRGBA{getColor(i * 5), 255, getColor(i * 5), 255})
		i++
	}
}

// Intended to be called every frame. Changes the NEXT movement so you can
// actually change your mind between movement ticks
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

// Gets where the snake will go with the current snake orientation
// Also sets snakeOrientation
func getNextSnakePosition() Coord {
	snakeOrientation = nextSnakeOrientation
	curCoord := snakeCoords.Front()
	if snakeOrientation == North {
		curCoord.y--
	}
	if snakeOrientation == West {
		curCoord.x--
	}
	if snakeOrientation == South {
		curCoord.y++
	}
	if snakeOrientation == East {
		curCoord.x++
	}
	return curCoord
}

// Just returns if the head is on the same coord as the apple
func headOnApple() bool {
	if Equals(appleCoord, snakeCoords.Front()) {
		return true
	}
	return false
}

// Draws the apple in the same way the snake is drawn, but red
func drawApple(screen *ebiten.Image) {
	x, y := coordToPixel(appleCoord)
	ebitenutil.DrawRect(screen, x+1, y+1, tileSize-2, tileSize-2,
		color.NRGBA{255, 0, 0, 255})
}

// Checks if c is colliding with the wall OR another component of the snake
func coordCollides(c Coord) bool {
	// head := snakeCoords.Front()
	// Check if out of bounds
	if c.x >= maxXTiles || c.y >= maxYTiles || c.x < 0 || c.y < 0 {
		return true
	}
	iter := snakeCoords.GetIterator()
	// Bypass the head
	iter.Next()
	for iter.Next() {
		if Equals(iter.Get(), c) {
			return true
		}
	}
	return false
}

// Checks if the apple coord is inside of the snake
// Coords will be generated inside of the bounds so we
// only need to check if it's in the snake
func validApplePos() bool {
	iter := snakeCoords.GetIterator()
	for iter.Next() {
		if Equals(appleCoord, iter.Get()) {
			return false
		}
	}
	return true
}

// Running the game!
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
