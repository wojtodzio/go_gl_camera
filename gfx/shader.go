package gfx

import (
	"io/ioutil"

	"github.com/go-gl/gl/v4.1-core/gl"

)

type Shader struct {
	glShader uint32
}

func (shader *Shader) glObject() uint32 {
	return shader.glShader
}

func (shader *Shader) Delete() {
	gl.DeleteShader(shader.glObject())
}

func NewShader(source string, shaderType uint32) (*Shader, error) {
	glShader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(glShader, 1, csources, nil)
	free()
	gl.CompileShader(glShader)

	wrappedShader := &Shader{glShader: glShader}

	err := getGlError(
		wrappedShader,
		gl.COMPILE_STATUS,
		gl.GetShaderiv,
		gl.GetShaderInfoLog,
		"Failed to compile",
	)
	if err != nil {
		return nil, err
	}

	return wrappedShader, nil
}

func NewShaderFromFile(filename string, shaderType uint32) (*Shader, error) {
	fileContent, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	shaderSource := string(fileContent) + "\x00"

	return NewShader(shaderSource, shaderType)
}
