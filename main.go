package main

import (
	"runtime"
	"log"

	"github.com/wojtodzio/go_gl_camera/gfx"
	"github.com/wojtodzio/go_gl_camera/win"
	"github.com/wojtodzio/go_gl_camera/cam"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"time"
)


const (
	width  = 1000
	height = 1000
)
func init() {
	runtime.LockOSThread()
}

var cubeVertices = []float32{
	// position        // texture position
	-0.5, -0.5, -0.5,  0.0, 0.0,
	 0.5, -0.5, -0.5,  1.0, 0.0,
	 0.5,  0.5, -0.5,  1.0, 1.0,
	 0.5,  0.5, -0.5,  1.0, 1.0,
	-0.5,  0.5, -0.5,  0.0, 1.0,
	-0.5, -0.5, -0.5,  0.0, 0.0,

	-0.5, -0.5,  0.5,  0.0, 0.0,
	 0.5, -0.5,  0.5,  1.0, 0.0,
	 0.5,  0.5,  0.5,  1.0, 1.0,
	 0.5,  0.5,  0.5,  1.0, 1.0,
	-0.5,  0.5,  0.5,  0.0, 1.0,
	-0.5, -0.5,  0.5,  0.0, 0.0,

	-0.5,  0.5,  0.5,  1.0, 0.0,
	-0.5,  0.5, -0.5,  1.0, 1.0,
	-0.5, -0.5, -0.5,  0.0, 1.0,
	-0.5, -0.5, -0.5,  0.0, 1.0,
	-0.5, -0.5,  0.5,  0.0, 0.0,
	-0.5,  0.5,  0.5,  1.0, 0.0,

	 0.5,  0.5,  0.5,  1.0, 0.0,
	 0.5,  0.5, -0.5,  1.0, 1.0,
	 0.5, -0.5, -0.5,  0.0, 1.0,
	 0.5, -0.5, -0.5,  0.0, 1.0,
	 0.5, -0.5,  0.5,  0.0, 0.0,
	 0.5,  0.5,  0.5,  1.0, 0.0,

	-0.5, -0.5, -0.5,  0.0, 1.0,
	 0.5, -0.5, -0.5,  1.0, 1.0,
	 0.5, -0.5,  0.5,  1.0, 0.0,
	 0.5, -0.5,  0.5,  1.0, 0.0,
	-0.5, -0.5,  0.5,  0.0, 0.0,
	-0.5, -0.5, -0.5,  0.0, 1.0,

	-0.5,  0.5, -0.5,  0.0, 1.0,
	 0.5,  0.5, -0.5,  1.0, 1.0,
	 0.5,  0.5,  0.5,  1.0, 0.0,
	 0.5,  0.5,  0.5,  1.0, 0.0,
	-0.5,  0.5,  0.5,  0.0, 0.0,
	-0.5,  0.5, -0.5,  0.0, 1.0,
}

var cubePositions = [][]float32 {
	{ 0.0,  0.0,  -3.0},
	{ 2.0,  5.0, -15.0},
	{-1.5, -2.2, -2.5 },
	{-3.8, -2.0, -12.3},
	{ 2.4, -0.4, -3.5 },
	{-1.7,  3.0, -7.5 },
	{ 1.3, -2.0, -2.5 },
	{ 1.5,  2.0, -2.5 },
	{ 1.5,  0.2, -1.5 },
	{-1.3,  1.0, -1.5 },
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

	vao := gfx.CreateVAO(cubeVertices)

	camera := cam.NewFPSCamera(mgl32.Vec3{0, 0, 3}, mgl32.Vec3{0, 1, 0}, -90, 0)

	log.Println("Entering program loop")
	for !window.ShouldClose() {
		window.StartFrame()

		camera.Update(window.SincePreviousFrame())

		// Clear the colorbuffer
		gl.ClearColor(0.2, 0.5, 0.5, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)  // depth buffer needed for DEPTH_TEST
		program.Use()

		cameraTransform := camera.GetTransform()
		program.Uniform("camera", &cameraTransform[0])

		aspect := float32(window.Width()) / float32(window.Height())
		program.UniformProject(aspect)

		vao.Bind()

		for _, position := range cubePositions {
			worldTransform := mgl32.Translate3D(position[0], position[1], position[2])

			program.Uniform("world", &worldTransform[0])
			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}

		vao.UnBind()

		time.Sleep(1 * time.Millisecond)
	}
	return nil
}
