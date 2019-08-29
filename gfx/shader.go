package gfx

import (
	"fmt"
	"strings"
	"io/ioutil"

	"github.com/go-gl/gl/v4.1-core/gl"

)

func CompileShaderFromFile(filename string, shaderType uint32) (uint32, error) {
	fileContent, err := ioutil.ReadFile(filename)

	if err != nil {
		return 0, err
	}

	shaderSource := string(fileContent) + "\x00"

	return CompileShader(shaderSource, shaderType)
}

func CompileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
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
