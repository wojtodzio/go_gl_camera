package main

import (
	"runtime"
	"log"

	"github.com/wojtodzio/go_gl_camera/gfx"
	"github.com/wojtodzio/go_gl_camera/win"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	"time"
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
		log.Println("Failed to init glfw")
		panic(err)
	}
	defer glfw.Terminate()

	window, err := win.NewWindow(width, height, "Virtual Camera")
	if err != nil {
		log.Println("Failed to create a new window")
		panic(err)
	}

	if err := initOpenGL(); err != nil {
		log.Println("Failed to init OpenGL")
		panic(err)
	}

	if err := programLoop(window); err != nil {
		log.Println("Program crashed")
		panic(err)
	}
}

func initGlfw() error {
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
	// ensures that triangles that are "behind" others do not draw over top of them
	gl.Enable(gl.DEPTH_TEST)

	return nil
}

func programLoop(window *win.Window) error {
	vertexShader, err := gfx.NewShaderFromFile("shaders/basic.vert", gl.VERTEX_SHADER)
	if err != nil {
		return err
	}
	fragmentShader, err := gfx.NewShaderFromFile("shaders/basic.frag", gl.FRAGMENT_SHADER)
	if err != nil {
		return err
	}

	program, err := gfx.NewProgram(vertexShader, fragmentShader)
	if err != nil {
		return err
	}

	defer program.Delete()

	log.Println("Entering program loop")
	for !window.ShouldClose() {
		window.StartFrame()

		// Clear the colorbuffer
		gl.ClearColor(0.2, 0.5, 0.5, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)  // depth buffer needed for DEPTH_TEST
		program.Use()

		// TODO: Remove, it's a quick hack to not kill my machine's performance
		time.Sleep(10 * time.Millisecond)

		gl.BindVertexArray(0)
	}
	return nil
}
