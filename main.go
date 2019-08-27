package main

import (
	"fmt"
	"strings"
	"runtime"
	"log"
	"io/ioutil"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)


const (
	width  = 1000
	height = 1000
)

var vertices = [...]float32{
	-0.5, -0.5, -0.5, 0.0, 0.0,
	0.5,  -0.5, -0.5, 1.0, 0.0,
	0.5,   0.5, -0.5, 1.0, 1.0,
	0.5,   0.5, -0.5, 1.0, 1.0,
	-0.5,  0.5, -0.5, 0.0, 1.0,
	-0.5, -0.5, -0.5, 0.0, 0.0,

	-0.5, -0.5,  0.5, 0.0, 0.0,
	 0.5, -0.5,  0.5, 1.0, 0.0,
	 0.5,  0.5,  0.5, 1.0, 1.0,
	 0.5,  0.5,  0.5, 1.0, 1.0,
	-0.5,  0.5,  0.5, 0.0, 1.0,
	-0.5, -0.5,  0.5, 0.0, 0.0,

	-0.5,  0.5,  0.5, 1.0, 0.0,
	-0.5,  0.5, -0.5, 1.0, 1.0,
	-0.5, -0.5, -0.5, 0.0, 1.0,
	-0.5, -0.5, -0.5, 0.0, 1.0,
	-0.5, -0.5,  0.5, 0.0, 0.0,
	-0.5,  0.5,  0.5, 1.0, 0.0,

	0.5,  0.5,  0.5, 1.0, 0.0,
	0.5,  0.5, -0.5, 1.0, 1.0,
	0.5, -0.5, -0.5, 0.0, 1.0,
	0.5, -0.5, -0.5, 0.0, 1.0,
	0.5, -0.5,  0.5, 0.0, 0.0,
	0.5,  0.5,  0.5, 1.0, 0.0,

	-0.5, -0.5, -0.5, 0.0, 1.0,
	 0.5, -0.5, -0.5, 1.0, 1.0,
	 0.5, -0.5,  0.5, 1.0, 0.0,
	 0.5, -0.5,  0.5, 1.0, 0.0,
	-0.5, -0.5,  0.5, 0.0, 0.0,
	-0.5, -0.5, -0.5, 0.0, 1.0,

	-0.5,  0.5, -0.5, 0.0, 1.0,
	 0.5,  0.5, -0.5, 1.0, 1.0,
	 0.5,  0.5,  0.5, 1.0, 0.0,
	 0.5,  0.5,  0.5, 1.0, 0.0,
	-0.5,  0.5,  0.5, 0.0, 0.0,
	-0.5,  0.5, -0.5, 0.0, 1.0,
}


func init() {
	runtime.LockOSThread()
}

func main() {
	window := initGlfw()
	defer glfw.Terminate()

	program := initOpenGL()
	program = program

	for !window.ShouldClose() {
	}
}

func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	// Setup window
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Virtual Camera", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	gl.Viewport(0, 0, width, height)

	// Enable depth testing
	// https://www.khronos.org/opengl/wiki/Depth_Test
	gl.Enable(gl.DEPTH_TEST)

	// Enable blending for alpha layer (transparency)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	vertexShader, err := compileShader("shaders/basic.vert", gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := compileShader("shaders/basic.frag", gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

func compileShader(file string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	source, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	csources, free := gl.Strs(string(source) + "\x00")
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
