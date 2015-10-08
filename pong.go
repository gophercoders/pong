package main

import (
	// This is the graphics library we are going to use. It is called the
	// Simple Direct Media Library. SDL for short. We need this to create the
	// window and to provide the drawing functions we need.
	"fmt"

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

// ---- Game State variables ----
// The quit flag this is used to control the main game loop.
// If quit is true then the user wants to finish the game. This will
// break the main game loop.
var quit bool

// my bats x and y position on the screen in pixels
var myBatX int
var myBatY int

// my bats width and height. This is the width and height of the grapic in pixels
var myBatW int
var myBatH int

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
	// now start the main game loop of the game.
	gameMainLoop()
}

// Initialise sets the inital values of the game state variables.
// Initialise must be called before the games main loop starts.
func initialise() {
	// initially set the quit flag to false.
	quit = false
	// load the game graphics
	loadGraphics()
	initialiseMyBatPosition()
}

// GameMainLoop controls the game. It performs three manin tasks. The first task
// is to get the users input. The second task is to update the games state based
// on the user input and the rules of the game. The final task is to update, or
// render, the changes to the screen.
func gameMainLoop() {
	for quit == false {
		getInput()
		updateState()
		render()
	}
}

func cleanup() {
	if myBat != nil {
		myBat.Destroy()
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
				myBatY = myBatY - myBatH/4
				// make sure we do not go off the top of the screen!
				if myBatY < 0 {
					myBatY = 0
				}
			}
			if isKeyDown(event) {
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
		}
		/*
			switch event.(type) {
			//		case *sdl.QuitEvent:
			//			quit = true
			case *sdl.KeyDownEvent:
				// which key was presses?
				var keyDownEvt *sdl.KeyDownEvent
				var ok bool
				keyDownEvt, ok = event.(*sdl.KeyDownEvent)
				if !ok {
					panic("KeyDownEvent type assertion failed!")
				}
				switch keyDownEvt.Keysym.Sym {
				case sdl.K_UP:
					myBatY = myBatY - myBatH/4
					// make sure we do not go off the top of the screen!
					if myBatY < 0 {
						myBatY = 0
					}
				case sdl.K_DOWN:
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
			}
		*/
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

// UpdateGameState updates the game state variables based on the user input and
// the rules of the game.
func updateState() {

}

// Render updates the screen, based on the new positions of the bats and the ball.
func render() {

	renderer.Clear()
	renderMyBat()
	// Show the empty window window we have just created.
	renderer.Present()
}

func loadGraphics() {
	loadMyBatGraphic()
	setSizeOfMyBat()
}

func loadMyBatGraphic() {
	myBat = loadBatGraphic("./assets/graphics/bat.png")
}

func loadBatGraphic(filename string) *sdl.Texture {
	var err error

	image, err = img.Load(filename)
	if err != nil {
		fmt.Print("Failed to load PNG: ")
		fmt.Println(err)
		panic(err)
	}
	defer image.Free()
	var bat *sdl.Texture
	bat, err = renderer.CreateTextureFromSurface(image)
	if err != nil {
		fmt.Print("Failed to create texture: ")
		fmt.Println(err)
		panic(err)
	}
	return bat
}

func initialiseMyBatPosition() {
	myBatX = windowWidth/10 - myBatW/2
	myBatY = windowHeight/2 - myBatH/2
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

	renderer.Clear()
	renderer.SetDrawColor(255, 0, 0, 255)
	var r sdl.Rect
	r.X = 0
	r.Y = 0
	r.W = int32(windowWidth)
	r.H = int32(windowHeight)
	renderer.FillRect(&r)
	renderer.Copy(myBat, &src, &dst)

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

func getPlayersInput() {
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
