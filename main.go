package main

import (
	"runtime"
	"log"

	_ "github.com/wojtodzio/go_gl_camera/gfx"
	"github.com/wojtodzio/go_gl_camera/win"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

)


const (
	width  = 1000
	height = 1000
)

func init() {
	runtime.LockOSThread()
}

func main() {
	if err := initGlfw(); err != nil {
		log.Fatalln("Failed to init glfw")
		panic(err)
	}
	defer glfw.Terminate()

	window, err := win.NewWindow(width, height, "Virtual Camera")
	if err != nil {
		log.Fatalln("Failed to create a new window")
		panic(err)
	}

	if err := initOpenGL(); err != nil {
		log.Fatalln("Failed to init OpenGL")
		panic(err)
	}

	for !window.ShouldClose() {
	}
}

func initGlfw() (error) {
	if err := glfw.Init(); err != nil {
		return err
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	return nil
}

func initOpenGL() error {
	if err := gl.Init(); err != nil {
		return err
	}

	// Enable depth testing
	// https://www.khronos.org/opengl/wiki/Depth_Test
	gl.Enable(gl.DEPTH_TEST)

	return nil
}
