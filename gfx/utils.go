package gfx

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type glWrapper interface {
	glObject() uint32
}

type getObjIv func(uint32, uint32, *int32)
type getObjInfoLog func(uint32, int32, *int32, *uint8)

func getGlError(
	wrappedObject glWrapper,
	checkTrueParam uint32,
	getObjIvFn getObjIv,
	getObjInfoLogFn getObjInfoLog,
	failMsg string,
) error {
	var success int32
	getObjIvFn(wrappedObject.glObject(), checkTrueParam, &success)

	if success == gl.FALSE {
		var logLength int32
		getObjIvFn(wrappedObject.glObject(), gl.INFO_LOG_LENGTH, &logLength)

		log := gl.Str(strings.Repeat("\x00", int(logLength)))
		getObjInfoLogFn(wrappedObject.glObject(), logLength, nil, log)

		return fmt.Errorf("%s: %s", failMsg, gl.GoStr(log))
	}

	return nil
}
