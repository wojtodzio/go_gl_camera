package cam

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	moveSpeed = 5.00
	cursorSensitivity = 0.05
)

type FPSCamera struct {
	// Eular Angles
	pitchRad,
	yawRad float32

	// Camera attributes
	position,
	front,
	up,
	right,
	worldUp mgl32.Vec3
}

func NewFPSCamera(position, worldUp mgl32.Vec3, yawDeg, pitchDeg float32) *FPSCamera {
	return &FPSCamera{
		pitchRad: mgl32.DegToRad(pitchDeg),
		yawRad: mgl32.DegToRad(yawDeg),
		position: position,
		up: mgl32.Vec3{0, 1, 0},
		worldUp: worldUp,
	}
}

func (camera *FPSCamera) Update(deltaTime float32) {
	camera.updatePosition(deltaTime)
	camera.updateVectors()
}

func (camera *FPSCamera) updatePosition(deltaTime float32) {
	distanceMoved := deltaTime * moveSpeed

	camera.position = camera.position.Add(camera.front.Mul(distanceMoved))
}

func (camera *FPSCamera) updateVectors() {
	// x, y, z
	cosPitch := math.Cos(float64(camera.pitchRad))
	sinPitch := math.Sin(float64(camera.pitchRad))
	cosYaw := math.Cos(float64(camera.yawRad))
	sinYaw := math.Sin(float64(camera.yawRad))
	camera.front[0] = float32(cosPitch * cosYaw)
	camera.front[1] = float32(sinPitch)
	camera.front[2] = float32(cosPitch * sinYaw)
	camera.front = camera.front.Normalize()

	// Gram-Schmidt process to figure out right and up vectors
	camera.right = camera.worldUp.Cross(camera.front).Normalize()
	camera.up = camera.right.Cross(camera.front).Normalize()
}

// Retruns a matrix to transform from the world coordinates to camera's coordinates
func (camera *FPSCamera) GetTransform() mgl32.Mat4 {
	cameraTarget := camera.position.Add(camera.front)

	return mgl32.LookAt(
		camera.position.X(), camera.position.Y(), camera.position.Z(),
		cameraTarget.X(), cameraTarget.Y(), cameraTarget.Z(),
		camera.up.X(), camera.up.Y(), camera.up.Z(),
	)
}
