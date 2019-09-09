package win

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Action glfw.Key

const (
	FORWARD = Action(glfw.KeyW)
	BACKWARD = Action(glfw.KeyS)
	LEFT = Action(glfw.KeyA)
	RIGHT = Action(glfw.KeyD)
	QUIT = Action(glfw.KeyEscape)
)

type InputManager struct {
	keysPressed [glfw.KeyLast]bool

	firstCursorAction bool
	cursor mgl32.Vec2
	cursorChange mgl32.Vec2
	cursorLast mgl32.Vec2
	bufferedCursorChange mgl32.Vec2
}

func NewInputManager() *InputManager {
	return &InputManager{firstCursorAction: false}
}

func (inputManager *InputManager) IsActive(action Action) bool {
	return inputManager.keysPressed[action]
}

func (inputManager *InputManager) Cursor() mgl32.Vec2 {
	return inputManager.cursor
}

func (inputManager *InputManager) CursorChange() mgl32.Vec2 {
	return inputManager.cursorChange
}

// Update Cursor() and CursorChange()
func (inputManager *InputManager) UpdateCursor() {
	inputManager.cursorChange[0] = inputManager.bufferedCursorChange[0]
	inputManager.cursorChange[1] = inputManager.bufferedCursorChange[1]
	inputManager.cursor[0] = inputManager.cursorLast[0]
	inputManager.cursor[1] = inputManager.cursorLast[1]

	inputManager.bufferedCursorChange[0] = 0
	inputManager.bufferedCursorChange[1] = 0
}

func (inputManager *InputManager) keyCallback(_ *glfw.Window, key glfw.Key, _ int, action glfw.Action, _ glfw.ModifierKey) {
	switch action {
	case glfw.Press:
		inputManager.keysPressed[Action(key)] = true
	case glfw.Release:
		inputManager.keysPressed[Action(key)] = false
	}
}

func (inputManager *InputManager) mouseCallback(_ *glfw.Window, xPositionF64, yPositionF64 float64) {
	xPosition := float32(xPositionF64)
	yPosition := float32(yPositionF64)

	if inputManager.firstCursorAction {
		inputManager.cursorLast[0] = xPosition
		inputManager.cursorLast[1] = yPosition
		inputManager.firstCursorAction = false
	}

	inputManager.bufferedCursorChange[0] += xPosition - inputManager.cursorLast[0]
	inputManager.bufferedCursorChange[1] += yPosition - inputManager.cursorLast[1]

	inputManager.cursorLast[0] = xPosition
	inputManager.cursorLast[1] = yPosition
}
