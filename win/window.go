package win

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

type Window struct {
	width,
	height int

	inputManager *InputManager

	deltaTime float32
	previousFrameTime float64

	glfwWindow *glfw.Window
}

func NewWindow(width, height int, title string) (*Window, error) {
	glfwWindow, err := glfw.CreateWindow(width, height, title, nil, nil)

	if err != nil {
		return nil, err
	}

	glfwWindow.MakeContextCurrent();
	glfwWindow.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	inputManager := NewInputManager()

	glfwWindow.SetKeyCallback(inputManager.keyCallback)
	glfwWindow.SetCursorPosCallback(inputManager.mouseCallback)

	return &Window{
		width: width,
		height: height,
		glfwWindow: glfwWindow,
		inputManager: inputManager,
	}, nil
}

func (window *Window) InputManager() *InputManager {
	return window.inputManager
}

func (window *Window) ShouldClose() bool {
	return window.glfwWindow.ShouldClose()
}

func (window *Window) StartFrame() {
	/// swap with the previous rendered buffer
	window.glfwWindow.SwapBuffers()

	// poll for UI window events
	glfw.PollEvents()

	if window.inputManager.IsActive(QUIT) {
		window.glfwWindow.SetShouldClose(true)
	}

	currentFrameTime := glfw.GetTime()

	if window.previousFrameTime == 0 {
		window.deltaTime = 0
	} else {
		window.deltaTime = float32(currentFrameTime - window.previousFrameTime)
	}
	window.previousFrameTime = currentFrameTime

	window.inputManager.UpdateCursor()
}

func (window *Window) SincePreviousFrame() float32 {
	return window.deltaTime
}

func (window *Window) Width() int {
	return window.width
}

func (window *Window) Height() int {
	return window.height
}
