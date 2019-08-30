package gfx

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	fovy = 60.0
	near = 0.1
	far = 100.0
)

type Program struct {
	glProgram uint32
	shaders []*Shader
}

func NewProgram(shaders ...*Shader) (*Program, error) {
	program := &Program{glProgram: gl.CreateProgram()}
	program.Attach(shaders...)

	if err := program.Link(); err != nil {
		return nil, err
	}

	return program, nil
}

func (program *Program) glObject() uint32 {
	return program.glProgram
}

func (program *Program) Attach(shaders ...*Shader) {
	for _, shader := range shaders {
		gl.AttachShader(program.glProgram, shader.glShader)
		program.shaders = append(program.shaders, shader)
	}
}

func (program *Program) Link() error {
	gl.LinkProgram(program.glProgram)

	return getGlError(
		program,
		gl.LINK_STATUS,
		gl.GetProgramiv,
		gl.GetProgramInfoLog,
		"PROGRAM::LINKING_FAILURE",
	)
}

func (program *Program) Delete() {
	for _, shader := range program.shaders {
		shader.Delete()
	}
	gl.DeleteProgram(program.glObject())
}

func (program *Program) Use() {
	gl.UseProgram(program.glObject())
}

func (program *Program) GetUniformLocation(name string) int32 {
	return gl.GetUniformLocation(program.glObject(), gl.Str(name + "\x00"))
}

// https://www.khronos.org/registry/OpenGL-Refpages/gl4/html/glUniform.xhtml
func (program *Program) Uniform(name string, value *float32) {
	gl.UniformMatrix4fv(
		program.GetUniformLocation(name),
		1,
		false,
		value,
	)
}

func (program *Program) UniformProject(aspect float32) {
	projectTransform := mgl32.Perspective(
		mgl32.DegToRad(fovy),
		aspect,
		near,
		far,
	)

	program.Uniform("project", &projectTransform[0])
}
