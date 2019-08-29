package win

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Window struct {
	width int
	height int
	glfwWindow *glfw.Window
}

func NewWindow(width, height int, title string) (*Window, error) {
	glfwWindow, err := glfw.CreateWindow(width, height, title, nil, nil)

	if err != nil {
		return nil, err
	}

	glfwWindow.MakeContextCurrent();
	glfwWindow.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	return &Window{
		width: width,
		height: height,
		glfwWindow: glfwWindow,
	}, nil
}

func (w *Window) ShouldClose() bool {
	return w.glfwWindow.ShouldClose()
}
