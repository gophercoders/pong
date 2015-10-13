package main

import (
	// This is the graphics library we are going to use. It is called the
	// Simple Direct Media Library. SDL for short. We need this to create the
	// window and to provide the drawing functions we need.
	"fmt"
	"math"

	"github.com/gophercoders/random"
	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

// These are the variables for the graphics library
// They have to be outside of the main function because the functions at the
// end of the file need them.
// This the window we are going to draw into
var window *sdl.Window

// This the abstraction to the graphics hardware inside the computer
// that actually does the drawing
var renderer *sdl.Renderer

// These variabels are important. They are the width and height of the window
// If you change these you will change the size of the image
var windowWidth int
var windowHeight int

// Image is a reusable surface used to load the game graohics.
var image *sdl.Surface

// myBat is the gaphic used to represent the the players bat
var myBat *sdl.Texture
var ball *sdl.Texture

// ---- Game State variables ----

// the balls speed in pixels per second
// This should never change during the game. We can make sure of this
// if we use define a constant value. Go provents us form changing the
// value of a constant - it's an illegal action - it breaks the rile of go.
const BallSpeed = 550

// The quit flag this is used to control the main game loop.
// If quit is true then the user wants to finish the game. This will
// break the main game loop.
var quit bool

// The paused flag is used to control if the game is paused.
// If the paused flag is true then the game is paused. The game will stop moving
// the ball and both players bats.
var paused bool

// my bats x and y position on the screen in pixels
var myBatX int
var myBatY int

// my bats width and height. This is the width and height of the grapic in pixels
var myBatW int
var myBatH int

// the balls x and y position on the screen in pixels
// We don't store these as int types! We want to know the exact position, so
// we can record fractions of a pixel. We will convert the number to int type
// just before we draw the ball on screen.
var ballX float64
var ballY float64

// the balls width and height. This is the width and height of the grapic in pixels
var ballW int
var ballH int

// the balls direction in x (horizontal) and y (vertical) across the screen
// Again we want to know the exact direction of travel, so we store these as
// float64 types not ints.
var ballDirX float64
var ballDirY float64

// The programs main function
func main() {
	// ---- This is the start of Owen's graphics setup code ----

	// First we have to initalise the SDL library, before we can use it
	sdl.Init(sdl.INIT_EVERYTHING)
	// defer is a go keyword and a special feature.
	// This means that go will automatically call the function sdl.Quit() before
	// the program exits for us. We don't have to remember to put this at the end!
	defer sdl.Quit()

	// if you want to change these try 800 for the width and 600 for the height
	windowWidth = 1024
	windowHeight = 768

	// Now we have to create the window we want to use.
	// We need to tell the SDL library how big to make the window of the correct
	// size - that's what the bit in the brackets does
	window = createWindow(windowWidth, windowHeight)
	// automatically destroy the window when the program finishes
	defer window.Destroy()
	// Now we have a window we need to create a renderer so we can draw into
	// it. In this case we want to use the first graphics card that supports faster
	// drawing
	renderer = createRenderer(window)
	// automatically destroy the renderer when the program exits.
	defer renderer.Destroy()

	// Set a black i.e. RGBA (0,0,0,0) background colour and clear the window
	renderer.SetDrawColor(0, 0, 0, 0)
	renderer.Clear()
	// ---- This is the end of Owen's graphics setup code ----

	// defer any cleanup actions
	defer cleanup()
	// initialise the games variables.
	initialise()
	// render everything initially so that we can see the game before it starts
	render()
	// now start the main game loop of the game.
	gameMainLoop()
}

// Initialise sets the inital values of the game state variables.
// Initialise must be called before the games main loop starts.
func initialise() {
	// initially set the quit flag to false.
	quit = false
	// initially the game is not paused
	paused = false
	// load the game graphics
	loadGraphics()
	initialiseMyBatPosition()
	initialiseBallPosition()
	initialiseBallDirection()
}

func initialiseBallDirection() {
	// pick some random numbers to determine if the ball will move up or down
	// and left or right initially.
	var n int
	n = random.GetRandomNumberInRange(1, 10)
	var up bool
	if isOddNumber(n) {
		up = true // we want the ball to move up - decreasing Y coordinate
	} else {
		up = false // we want the ball to move down - increasing y coordinate
	}

	n = random.GetRandomNumberInRange(1, 10)
	var left bool
	if isOddNumber(n) {
		left = true // we want the ball to move left - decreasing X coordiiate
	} else {
		left = false //we want the ball to move right - increasing X coordinate
	}
	// pick two random mumbers for the initial direction
	ballDirX = float64(random.GetRandomNumberInRange(1, 10))
	ballDirY = float64(random.GetRandomNumberInRange(1, 10))
	// are we moving left?
	if left {
		ballDirX = ballDirX * -1
	} // otherwise the ball is moving right so ballDirX should be positive
	// are we moving up?
	if up {
		// yes - make the number negative
		ballDirY = ballDirY * -1
	} // otherwise the ball is moving dwon so ballDirY should be positive
	// the vector now needs to be normalised
	setBallDirection(ballDirX, ballDirY)
}

func setBallDirection(newDirectionX, newDirectionY float64) {
	// nornalise the direction vector
	var length float64
	// To mornalise the vector multiply each side by itself, and then add then
	// results together
	length = float64(newDirectionX*newDirectionX + newDirectionY*newDirectionY)
	// then take the square root
	length = math.Sqrt(length)
	// We want to keep the balls speed constant so
	// the balls new position (in each direction) is the balls speed (in each direction)
	// multiplied by _scalled_ new direction (in each direction)
	ballDirX = BallSpeed * (newDirectionX / length)
	ballDirY = BallSpeed * (newDirectionY / length)
}

func isOddNumber(number int) bool {
	// use the modulus operator (%) to divide the number by two and
	// return the _remainder_
	// An even number has no remainder so false is returned.
	// An odd number has a remainder so true is returned
	if number%2 == 0 {
		return false
	}
	return true
}

// GameMainLoop controls the game. It performs three manin tasks. The first task
// is to get the users input. The second task is to update the games state based
// on the user input and the rules of the game. The final task is to update, or
// render, the changes to the screen.
func gameMainLoop() {
	for quit == false {
		getInput()
		// if the game is not paused, then we must update the games state
		if paused == false {
			updateState()
		}
		render()
	}
}

func cleanup() {
	if myBat != nil {
		myBat.Destroy()
	}
	if ball != nil {
		ball.Destroy()
	}
}

// GetInput gets the users input and updates the game state variables that realte
// to the users input, for example, the direction that the user wants to move their
// bat in.
func getInput() {
	var event sdl.Event
	event = sdl.PollEvent()
	if event != nil {
		if isQuitEvent(event) {
			quit = true
		}
		if isKeyDownEvent(event) {
			if isKeyUp(event) {
				// If the game is paused we must iognore the up cursor key.
				// We do not want to move the players bat up if the game is
				// paused. We can do this by leaving the getInput function
				// early, by using the return keyword.
				if paused == true {
					return
				}
				// if the game is not paused we can process the key press
				myBatY = myBatY - myBatH/4
				// make sure we do not go off the top of the screen!
				if myBatY < 0 {
					myBatY = 0
				}
			}
			if isKeyDown(event) {
				// If the game is paused we must iognore the down cursor key.
				// We do not want to move the players bat up if the game is
				// paused. We can do this by leaving the getInput function
				// early, by using the return keyword.
				if paused == true {
					return
				}
				// if the game is not paused we can process the key
				myBatY = myBatY + myBatH/4
				// make sure we do not go off the bottom of the screen
				// we have to account for the heigh of the bat when we do this
				// becase myBatY is the Y coordinate of the top right of the bat,
				// but the bottom right (or left) will go of the bottom of the
				// screen first.
				if myBatY+myBatH > windowHeight {
					myBatY = windowHeight - myBatH
				}
			}
			// We must always respond to the paused key being pressed.
			// If the game is running the pause key pauses the game.
			// But if the game is paused, we must still respond to the paused key.
			// This is the only way to unpause the game.
			if isKeyPause(event) {
				if paused == true {
					paused = false
				} else {
					paused = true
				}
			}
		}
	}
}

func isQuitEvent(event sdl.Event) bool {
	var ok bool
	_, ok = event.(*sdl.QuitEvent)
	return ok
}

func isKeyDownEvent(event sdl.Event) bool {
	var ok bool
	_, ok = event.(*sdl.KeyDownEvent)
	return ok
}

func isKeyUp(event sdl.Event) bool {
	var keyDownEvt *sdl.KeyDownEvent
	var ok bool
	keyDownEvt, ok = event.(*sdl.KeyDownEvent)
	if !ok {
		panic("KeyDownEvent type assertion failed!")
	}
	return (keyDownEvt.Keysym.Sym == sdl.K_UP)
}

func isKeyDown(event sdl.Event) bool {
	var keyDownEvt *sdl.KeyDownEvent
	var ok bool
	keyDownEvt, ok = event.(*sdl.KeyDownEvent)
	if !ok {
		panic("KeyDownEvent type assertion failed!")
	}
	return (keyDownEvt.Keysym.Sym == sdl.K_DOWN)
}

func isKeyPause(event sdl.Event) bool {
	var keyDownEvt *sdl.KeyDownEvent
	var ok bool
	keyDownEvt, ok = event.(*sdl.KeyDownEvent)
	if !ok {
		panic("KeyDownEvent type assertion failed!")
	}
	return (keyDownEvt.Keysym.Sym == sdl.K_PAUSE)
}

// UpdateGameState updates the game state variables based on the user input and
// the rules of the game.
func updateState() {
	// update the balls state
	updateBallState()
	checkForCollisions()
}

func updateBallState() {
	// just update the position.....
	var frameTime = float64(1) / float64(60)
	// work out how far the ball moved during the last "frame"
	// Easy - just the direction times the frameTime
	var xDelta = ballDirX * frameTime
	var yDelta = ballDirY * frameTime
	// the balls new position is the last position + the delta for this frame
	ballX = ballX + xDelta
	ballY = ballY + yDelta
}

func checkForCollisions() {
	checkForBallWallCollisions()
}

func checkForBallWallCollisions() {
	// did the ball bit the top or the bottom wall or the left or right wall?
	var playingFieldBottom = float64(windowHeight) - float64(ballH)
	var playingFieldRight = float64(windowWidth) - float64(ballW)
	// Check the top first
	if ballY < 0 { // the top of the wall is a Y cooordinate on the screen
		// stop the ball from going off the top of the screen
		ballY = 0.0
		// yes we hit the top, so reflect the ball back by changing
		ballDirY = ballDirY * -1
	} else if ballY > playingFieldBottom {
		// we hit the bottom so stop the ball from going off the bottom of the
		// screen
		ballY = playingFieldBottom
		// now reflect the ball back
		ballDirY = ballDirY * -1

	}
	//check for left wall next
	if ballX < 0 { // the left edge of the x coordinate of the screen
		// we hit the left wall
		// stop the ball from going off the left edge
		ballX = 0.0
		// now reflect the ball back
		ballDirX = ballDirX * -1
	} else if ballX > playingFieldRight {
		// we hit the right wall
		ballX = playingFieldRight
		// now reflect back
		ballDirX = ballDirX * -1
	}
}

// Render updates the screen, based on the new positions of the bats and the ball.
func render() {
	var fps uint32
	fps = 60
	var delay uint32
	delay = 1000 / fps

	var frameStart uint32
	frameStart = sdl.GetTicks()

	renderer.Clear()
	renderMyBat()
	renderBall()
	// Show the empty window window we have just created.
	renderer.Present()

	var frameTime uint32
	frameTime = sdl.GetTicks() - frameStart
	if frameTime < delay {
		sdl.Delay(delay - frameTime)
	}
}

func loadGraphics() {
	loadMyBatGraphic()
	setSizeOfMyBat()
	loadBallGraphic()
	setSizeOfBall()
}

func loadMyBatGraphic() {
	myBat = loadGraphic("./assets/graphics/bat.png")
}

func loadBallGraphic() {
	ball = loadGraphic("./assets/graphics/ball.png")
}

func loadGraphic(filename string) *sdl.Texture {
	var err error

	image, err = img.Load(filename)
	if err != nil {
		fmt.Print("Failed to load PNG: ")
		fmt.Println(err)
		panic(err)
	}
	defer image.Free()
	var graphic *sdl.Texture
	graphic, err = renderer.CreateTextureFromSurface(image)
	if err != nil {
		fmt.Print("Failed to create texture: ")
		fmt.Println(err)
		panic(err)
	}
	return graphic
}

func initialiseMyBatPosition() {
	myBatX = windowWidth/10 - myBatW/2
	myBatY = windowHeight/2 - myBatH/2
}

func initialiseBallPosition() {
	ballX = float64(windowWidth/2 - ballW/2)
	ballY = float64(windowHeight/2 - ballH/2)
}

func setSizeOfMyBat() {
	var w, h int32
	var err error
	_, _, w, h, err = myBat.Query()
	if err != nil {
		fmt.Print("Failed to query texture: ")
		fmt.Println(err)
		panic(err)
	}
	myBatW = int(w)
	myBatH = int(h)
}

func setSizeOfBall() {
	var w, h int32
	var err error
	_, _, w, h, err = ball.Query()
	if err != nil {
		fmt.Print("Failed to query texture: ")
		fmt.Println(err)
		panic(err)
	}
	ballW = int(w)
	ballH = int(h)
}

func renderMyBat() {

	var src, dst sdl.Rect

	src.X = 0
	src.Y = 0
	src.W = int32(myBatW)
	src.H = int32(myBatH)

	dst.X = int32(myBatX)
	dst.Y = int32(myBatY)
	dst.W = int32(myBatW)
	dst.H = int32(myBatH)

	renderer.Copy(myBat, &src, &dst)

}

func renderBall() {

	var src, dst sdl.Rect

	src.X = 0
	src.Y = 0
	src.W = int32(ballW)
	src.H = int32(ballH)

	dst.X = int32(ballX)
	dst.Y = int32(ballY)
	dst.W = int32(ballW)
	dst.H = int32(ballH)

	renderer.Copy(ball, &src, &dst)

}

// CheckQuit checks if the user has clicked the window's close button.
// If the user has then the quit variable is set it true. CheckQuit returns
// the value of the quit variable.
func checkQuit() bool {
	var event sdl.Event
	event = sdl.PollEvent()

	if event != nil {
		switch event.(type) {
		case *sdl.QuitEvent:
			quit = true
		}
	}
	return quit
}

// Create the graphics window using the SDl library or crash trying
func createWindow(w, h int) *sdl.Window {
	var window *sdl.Window
	var err error

	window, err = sdl.CreateWindow("Pong Game", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		w, h, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	return window
}

// Create the graphics renderer or crash trying
func createRenderer(w *sdl.Window) *sdl.Renderer {
	var r *sdl.Renderer
	var err error
	r, err = sdl.CreateRenderer(w, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	return r
}
