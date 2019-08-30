package gfx

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type VAO struct {
	glVAO uint32
}

func (vao *VAO) glObject() uint32 {
	return vao.glVAO
}

// Vertex Array Object for a triangle
func CreateVAO(vertices []float32) *VAO {
	var glVAO uint32
	gl.GenVertexArrays(1, &glVAO)

	var glVBO uint32
	gl.GenBuffers(1, &glVBO)

	var glEBO uint32;
	gl.GenBuffers(1, &glEBO)

	gl.BindVertexArray(glVAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, glVBO)

	// copy vertices data into VBO
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices) * 4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// size of one whole vertex
	var stride int32 = 3 * 4 + 2 * 4
	var offset int = 0

	// position
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, gl.PtrOffset(offset))
	gl.EnableVertexAttribArray(0)
	offset += 3 * 4

	// texture position
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, stride, gl.PtrOffset(offset))
	gl.EnableVertexAttribArray(1)
	offset += 2 * 4

	// unbind the VAO
	gl.BindVertexArray(0)

	return &VAO{glVAO: glVAO}
}

func (vao *VAO) Bind() {
	gl.BindVertexArray(vao.glObject())
}

func (vao *VAO) UnBind() {
	gl.BindVertexArray(0)
}
