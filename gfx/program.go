package gfx

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Program struct {
	glProgram uint32
	shaders []*Shader
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

func NewProgram(shaders ...*Shader) (*Program, error) {
	program := &Program{glProgram: gl.CreateProgram()}
	program.Attach(shaders...)

	if err := program.Link(); err != nil {
		return nil, err
	}

	return program, nil
}
