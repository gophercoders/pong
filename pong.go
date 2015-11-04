package main

import (
	// This is the graphics library we are going to use. It is called the
	// Simple Direct Media Library. SDL for short. We need this to create the
	// window and to provide the drawing functions we need.
	"fmt"
	"math"
	"strconv"

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

// computersBat is the graphic used to represent the computers bat
var computersBat *sdl.Texture
var ball *sdl.Texture

// There is one graphic for each score, so one graphic for one, one for two, one for three...
// We have 12 graphics in total for the numbers 0 to 11 (inclusive). We will always have
// exactly 12 graphics because the game rules say that the game is over when a players score
// reaches 11.
// The natural way to store these grapincs is in an array. You can think of an array as a row of
// "boxes" in memory. Each box can hold the same "thing" and each box is the same size.
// Each box is given a number, it's index, starting at ZERO. We can use the index to retreieve the
// contents of any of the boxes in memory. Each of our boxes is big enough to store one graphic.

// The scoresGfx - 0 to 11 - are stored in an array. We need 12 spaces in the array because we count the zero.
var scoresGfx [12]*sdl.Texture

// The game over graphic
var gameOverGfx *sdl.Texture

// ---- Game State variables ----

// the balls speed in pixels per second
// This should never change during the game. We can make sure of this
// if we use define a constant value. Go provents us form changing the
// value of a constant - it's an illegal action - it breaks the rile of go.
const BallSpeed = 550

// This is the speed of the computers bat, in pixels per second
const ComputersBatSpeed = 350

// This is the number of points a player has to score to win the game
const WinningScore = 11

// The quit flag this is used to control the main game loop.
// If quit is true then the user wants to finish the game. This will
// break the main game loop.
var quit bool

// The paused flag is used to control if the game is paused.
// If the paused flag is true then the game is paused. The game will stop moving
// the ball and both players bats.
var paused bool

// The game over flag controls how the game ends.
// If a players score is equal to te WinningScore the game is over so
// the game over flag is true. When the game over flag is true the AI stops
// the ball stops moving, and the player cannot move or pause the game.
// The game over game over flag starts as false
var gameOver bool

// my bats x and y position on the screen in pixels
var myBatX int
var myBatY int

// my bats width and height. This is the width and height of the bat graphic in pixels
var myBatW int
var myBatH int

// the computers bats x and y position on the screen in pixels
var computersBatX int
var computersBatY int

// the computers bats width and height. This is the width and height of the bat graphic in pixels
var computersBatW int
var computersBatH int

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

// the scores for each player
var myScore int
var computersScore int

// the position of my score on the screen in pixels
var myScoreX int
var myScoreY int

// the position of the computers score on the screen in pixels
var computersScoreX int
var computersScoreY int

// the size of the score graphic in pixels
var scoreW int
var scoreH int

// the position of the game over graphic in screen pixels
var gameOverX int
var gameOverY int

// the gameOver graphic size
var gameOverW int
var gameOverH int

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
	// initially the gmae is not over
	gameOver = false
	// The scores start at zero
	myScore = 0
	computersScore = 0
	// load the game graphics
	loadGraphics()
	initialiseMyBatPosition()
	initialiseComputersBatPosition()
	initialiseBallPosition()
	initialiseBallDirection()
	initialiseScorePositions()
	initialiseGameOverPosition()
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
	if computersBat != nil {
		computersBat.Destroy()
	}
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
				// if the game is in the game over state we must also ignore
				// the key press
				if gameOver == true {
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
				// if the game is in the game over state we must also ignore
				// the key press
				if gameOver == true {
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
			// We must always respond to the paused key being pressed - if the
			// game is not over.
			// If the game is running the pause key pauses the game.
			// But if the game is paused, we must still respond to the paused key.
			// This is the only way to unpause the game.
			if isKeyPause(event) {
				// if the game is in the game over state we must also ignore
				// the key press
				if gameOver == true {
					return
				}
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
	// if the game has finished we must do nothing
	if gameOver == true {
		return
	}
	// update the balls state
	updateBallState()
	// move the computer players bat
	updateComputersBatPosition()
	// now check for collisions between the ball/walls and the ball/bats
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

// This is the games artifical intelligence.
// The rules that the move the computers bat are simple:
//    Chase the ball when it is going towards the AI player.
//    Only chase the ball when it is on its side of the playing field.
//    Move towards the center when the ball is moving away from the AI player.
func updateComputersBatPosition() {
	// work out the frame time
	var frameTime = float64(1) / float64(60)
	// calculate how far the computers bat could have moved in the frame time
	var deltaY float64
	deltaY = ComputersBatSpeed * frameTime

	var middleOfBatY float64
	middleOfBatY = float64(computersBatY) + float64(computersBatH/2)
	var middleOfBallY float64
	middleOfBallY = ballY + float64(ballH/2)
	var middleOfTheScreenY float64
	middleOfTheScreenY = float64(windowHeight / 2)
	// if the ball is on the computers half of the screen and the ball is moving
	// towards the computers bat then we want to move the computers bat towards the ball
	if ballX > float64(windowWidth/2) && ballDirX > 0 {
		// if the middle of the bat is above the middle of the ball
		// we want to move the bat down
		if middleOfBatY < middleOfBallY {
			//move the bat down, but the maximum about we can
			computersBatY = int((float64(computersBatY) + deltaY))
		} else if middleOfBatY > middleOfBallY {
			// if the middle of the bat is below the middle of the ball
			// we want to move move the bat up
			computersBatY = int(float64(computersBatY) - deltaY)
		}
	} else {
		// move to the center when the ball is far away or moving away
		// if the bat is above the middle od the screen move down
		if middleOfBatY < middleOfTheScreenY {
			computersBatY = int(float64(computersBatY) + deltaY)
		} else if middleOfBatY > middleOfTheScreenY {
			// if the bat is below the middle of the screen move up
			computersBatY = int(float64(computersBatY) - deltaY)
		}
	}
	// check the computers bat has not moved off the screen
	// Check the top first
	if computersBatY < 0 {
		computersBatY = 0
	}
	// now check the bottom
	if computersBatY+computersBatH > windowHeight {
		computersBatY = windowHeight - computersBatH
	}
}

func checkForCollisions() {
	checkForBallWallCollisions()
	// check to see if we hit the players bat
	var hitPlayersBat bool
	hitPlayersBat = checkForBallPayersBatCollisions()
	if hitPlayersBat == true {
		// the ball hit the players bat, so reflect it along a new direction
		reflectBallFromPlayersBat()
	}
	// check to see if the ball hit the computers bat
	var computersBatHit bool
	computersBatHit = checkForBallComputersBatCollisions()
	if computersBatHit == true {
		// the ball hit the players bat, so reflect it along a new direction
		reflectBallFromComputersBat()
	}
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
		// the ball hit the left wall, so the computer scored a point
		computersScore = computersScore + 1
		// now we need to reset the game state so that the ball starts
		// in the middle again.
		resetGameState()
		if computersScore == WinningScore {
			gameOver = true
		}
	} else if ballX > playingFieldRight {
		// we hit the right wall so the player scored a point
		myScore = myScore + 1
		resetGameState()
		if myScore == WinningScore {
			gameOver = true
		}
	}
}

func resetGameState() {
	// We want to reset the gmae state after a point is scored.
	// So we want to put the ball back in the centre of the screen again
	// and then "serve" it towrds one of the players.
	// Set the bals initial poistion
	initialiseBallPosition()
	// Now we need to set the balls direction
	initialiseBallDirection()
}
func checkForBallPayersBatCollisions() bool {
	// Did the ball collide with the players bat?
	// We need to look for an overlap between the bounding box of the ball and the
	// bounding box of the bat. If we find an overalp we need to return
	// true, if not we need to return false

	// if the right of the ball is less than the left of the bat - no collision
	if ballX+float64(ballW) < float64(myBatX) {
		return false
	}
	// is the left of the ball is greater than the right of the bat - no collision
	if ballX > float64(myBatX+myBatW) {
		return false
	}
	// if the bottom of the ball is less than the top of the bat - no collision
	if ballY+float64(ballH) < float64(myBatY) {
		return false
	}
	// if the top of the ball is greater then the bottom of the bat - no collision
	if ballY > float64(myBatY+myBatH) {
		return false
	}
	// otherwise some part of the bat and ball overlap - so there is a collision
	return true
}

// reflect the ball back from the players bat. If the ball hit above the
// middle of the bat reflect the ball upwards. If the ball hit below the middle
// of the bat reflect downwards.
// The direction of reflection depends on where the ball hit the bat. The
// further away from the middle of the bat the greater the angle of reflection.
// The angle is determined by the vector that represents the balls direction
// At the top of the bat the reflected direction is [1,2]. At the bottom of the
// bat the reflected direction is [1,-2].
func reflectBallFromPlayersBat() {
	// The ball has hit the players bat. So make the left most part of the
	// ball - ballX - line up with the right most part of the bat - myBatX + myBatW
	ballX = float64(myBatX + myBatW)
	// we need to know where the center of the ball is (in the Y axis) so we
	// can work out where it hit on the bat.
	var ballCentreY float64
	ballCentreY = ballY + float64(ballH/2)

	// Now we need to work out where the ball hit on the bat, the hitpoint.
	var hitPoint float64
	hitPoint = ballCentreY

	// clip the hit point so that it is within the bat
	// If the hitpoint is above the bat, make it the top of thr bat
	if hitPoint < float64(myBatY) {
		hitPoint = float64(myBatY)
	} else if hitPoint > float64(myBatY+myBatH) {
		// if the hitpoint is below the bottom of the bat make it the bottomof the bat
		hitPoint = float64(myBatY + myBatH)
	}
	// Scale the hitpoint so that it is between zero and the height of the bat.
	// This is easy we just need to subtract the Y coordinate of the top of the bat.
	// This in effect translates the hit point to a position relative to the
	// screens origin.
	hitPoint = hitPoint - float64(myBatY)
	// We want to ensure that everything that hits above the centre line of the bat is
	// reflected upwards but everything that hits below the centre line is reflected
	// downwards. When we reflect upwards we want the Y coordinate to be negative
	// - remeber that a negative change in the Y coordiante moves up the screen -
	// When we want to go downwards we need a positive change in the Y coordinate.
	// We can arrange this by simply subtracting half the bats height from the
	// current value of the hit point. If we hit on the top half of the bat, then
	// hitpoint will be negative so the ball will be reflected upwards. This is
	// is exactly what we want.
	hitPoint = hitPoint - float64(myBatH/2)
	// now we need to calcualte the exaclt vector that we need to reflect
	// the ball along. This is easy we just need to scale the hitpoint so that
	// it lies between -2.0 and +2.0.
	var scaledReflection float64
	scaledReflection = 2.0 * (hitPoint / (float64(myBatH) / 2.0))
	// now we know the vector we want to change the balls direction to, we can set it.
	// The 1 just means the vector always goes to the right
	setBallDirection(1, scaledReflection)
}

func checkForBallComputersBatCollisions() bool {
	// Did the ball collide with the computers bat?
	// We need to look for an overlap between the bounding box of the ball and the
	// bounding box of the bat. If we find an overalp we need to return
	// true, if not we need to return false

	// if the right of the ball is less than the left of the bat - no collision
	if ballX+float64(ballW) < float64(computersBatX) {
		return false
	}
	// is the left of the ball is greater than the right of the bat - no collision
	if ballX > float64(computersBatX+computersBatW) {
		return false
	}
	// if the bottom of the ball is less than the top of the bat - no collision
	if ballY+float64(ballH) < float64(computersBatY) {
		return false
	}
	// if the top of the ball is greater then the bottom of the bat - no collision
	if ballY > float64(computersBatY+computersBatH) {
		return false
	}
	// otherwise some part of the bat and ball overlap - so there is a collision
	return true
}

func reflectBallFromComputersBat() {
	// The ball has hit the players bat. So make the left most part of the
	// ball - ballX - line up with the left most part of the bat - myBatX -
	// but set back by the balls width
	ballX = float64(computersBatX - ballW)
	// we need to know where the center of the ball is (in the Y axis) so we
	// can work out where it hit on the bat.
	var ballCentreY float64
	ballCentreY = ballY + float64(ballH/2)

	// Now we need to work out where the ball hit on the bat, the hitpoint.
	var hitPoint float64
	hitPoint = ballCentreY

	// clip the hit point so that it is within the bat
	// If the hitpoint is above the bat, make it the top of thr bat
	if hitPoint < float64(computersBatY) {
		hitPoint = float64(computersBatY)
	} else if hitPoint > float64(computersBatY+computersBatH) {
		// if the hitpoint is below the bottom of the bat make it the bottomof the bat
		hitPoint = float64(computersBatY + computersBatH)
	}
	// Scale the hitpoint so that it is between zero and the height of the bat.
	// This is easy we just need to subtract the Y coordinate of the top of the bat.
	// This in effect translates the hit point to a position relative to the
	// screens origin.
	hitPoint = hitPoint - float64(computersBatY)
	// We want to ensure that everything that hits above the centre line of the bat is
	// reflected upwards but everything that hits below the centre line is reflected
	// downwards. When we reflect upwards we want the Y coordinate to be negative
	// - remeber that a negative change in the Y coordiante moves up the screen -
	// When we want to go downwards we need a positive change in the Y coordinate.
	// We can arrange this by simply subtracting half the bats height from the
	// current value of the hit point. If we hit on the top half of the bat, then
	// hitpoint will be negative so the ball will be reflected upwards. This is
	// is exactly what we want.
	hitPoint = hitPoint - float64(computersBatH/2)
	// now we need to calcualte the exaclt vector that we need to reflect
	// the ball along. This is easy we just need to scale the hitpoint so that
	// it lies between -2.0 and +2.0.
	var scaledReflection float64
	scaledReflection = 2.0 * (hitPoint / (float64(computersBatH) / 2.0))
	// now we know the vector we want to change the balls direction to, we can set it.
	// The -1 just means the vector always goes to the left
	setBallDirection(-1, scaledReflection)
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
	renderComputersBat()
	renderScore()
	// if the game is over render the gameOver graphic
	if gameOver == true {
		renderGameOver()
	} else {
		// otherwise we need to draw the ball
		renderBall()
	}
	// Show the game window window.
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
	loadComputersBatGraphic()
	setSizeOfComputersBat()
	loadBallGraphic()
	setSizeOfBall()
	loadScores()
	setSizeOfScore()
	loadGameOverGraphic()
	setSizeOfGameOverGraphic()
}

func loadMyBatGraphic() {
	myBat = loadGraphic("./assets/graphics/bat.png")
}

func loadComputersBatGraphic() {
	computersBat = loadGraphic("./assets/graphics/bat.png")
}

func loadBallGraphic() {
	ball = loadGraphic("./assets/graphics/ball.png")
}

func loadScores() {
	var scoreGfxFilename string
	var i int
	for i = 0; i < 12; i++ {
		// create the score name dynamiclly
		// The filenames look like this:
		// ./assets/graphics/1.png or ./assets/graphics/10.png
		// So we can use the value of the loop counter - i- to help generate the
		// file name in the loop.
		// The strconv.Itoa function converts a decimal integer number to a string
		// the plus - +'s - join the strings together
		scoreGfxFilename = "./assets/graphics/" + strconv.Itoa(i) + ".png"
		// load the graphic and store it the i'th position in the scores array
		// So the first score is at scoresGfx[0] because i started at zero. The next
		// graphic is at scoresGfx[1] because the next value of i is one.
		scoresGfx[i] = loadGraphic(scoreGfxFilename)
	}
}

func loadGameOverGraphic() {
	gameOverGfx = loadGraphic("./assets/graphics/GameOver.png")
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

func initialiseComputersBatPosition() {
	computersBatX = windowWidth - (windowWidth / 10) - computersBatW/2
	computersBatY = windowHeight/2 - myBatH/2
}

func initialiseBallPosition() {
	ballX = float64(windowWidth/2 - ballW/2)
	ballY = float64(windowHeight/2 - ballH/2)
}

func initialiseScorePositions() {
	myScoreX = windowWidth/4 - scoreW/2
	myScoreY = windowHeight / 8
	computersScoreX = (windowWidth/4)*3 - scoreW/2
	computersScoreY = windowHeight / 8
}

func initialiseGameOverPosition() {
	gameOverX = windowWidth/2 - gameOverW/2
	gameOverY = windowHeight/2 - gameOverH/2
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

func setSizeOfComputersBat() {
	var w, h int32
	var err error
	_, _, w, h, err = computersBat.Query()
	if err != nil {
		fmt.Print("Failed to query texture: ")
		fmt.Println(err)
		panic(err)
	}
	computersBatW = int(w)
	computersBatH = int(h)
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

func setSizeOfScore() {
	var w, h int32
	var err error
	_, _, w, h, err = scoresGfx[0].Query() // All the score graphics are the same size, so we can use the first one
	if err != nil {
		fmt.Print("Failed to query texture: ")
		fmt.Println(err)
		panic(err)
	}
	scoreW = int(w)
	scoreH = int(h)
}

func setSizeOfGameOverGraphic() {
	var w, h int32
	var err error
	_, _, w, h, err = gameOverGfx.Query()
	if err != nil {
		fmt.Print("Failed to query texture: ")
		fmt.Println(err)
		panic(err)
	}
	gameOverW = int(w)
	gameOverH = int(h)
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

func renderComputersBat() {

	var src, dst sdl.Rect

	src.X = 0
	src.Y = 0
	src.W = int32(computersBatW)
	src.H = int32(computersBatH)

	dst.X = int32(computersBatX)
	dst.Y = int32(computersBatY)
	dst.W = int32(computersBatW)
	dst.H = int32(computersBatH)

	renderer.Copy(computersBat, &src, &dst)

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

func renderScore() {
	renderMyScore()
	renderComputersScore()
}

func renderMyScore() {
	var src, dst sdl.Rect

	src.X = 0
	src.Y = 0
	src.W = int32(scoreW)
	src.H = int32(scoreH)

	dst.X = int32(myScoreX)
	dst.Y = int32(myScoreY)
	dst.W = int32(scoreW)
	dst.H = int32(scoreH)

	renderer.Copy(scoresGfx[myScore], &src, &dst)
}

func renderComputersScore() {
	var src, dst sdl.Rect

	src.X = 0
	src.Y = 0
	src.W = int32(scoreW)
	src.H = int32(scoreH)

	dst.X = int32(computersScoreX)
	dst.Y = int32(computersScoreY)
	dst.W = int32(scoreW)
	dst.H = int32(scoreH)

	renderer.Copy(scoresGfx[computersScore], &src, &dst)
}

func renderGameOver() {
	var src, dst sdl.Rect

	src.X = 0
	src.Y = 0
	src.W = int32(gameOverW)
	src.H = int32(gameOverH)

	dst.X = int32(gameOverX)
	dst.Y = int32(gameOverY)
	dst.W = int32(gameOverW)
	dst.H = int32(gameOverH)

	renderer.Copy(gameOverGfx, &src, &dst)
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
