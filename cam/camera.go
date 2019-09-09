package cam

import (
	"math"

	"github.com/wojtodzio/go_gl_camera/win"

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

	inputManager *win.InputManager
}

func NewFPSCamera(position, worldUp mgl32.Vec3, yawDeg, pitchDeg float32, window *win.Window) *FPSCamera {
	return &FPSCamera{
		pitchRad: mgl32.DegToRad(pitchDeg),
		yawRad: mgl32.DegToRad(yawDeg),
		position: position,
		up: mgl32.Vec3{0, 1, 0},
		worldUp: worldUp,
		inputManager: window.InputManager(),
	}
}

func (camera *FPSCamera) Update(deltaTime float32) {
	camera.updatePosition(deltaTime)
	camera.updateDirection()
	camera.updateVectors()
}

func (camera *FPSCamera) updatePosition(deltaTime float32) {
	distanceMoved := deltaTime * moveSpeed

	if camera.inputManager.IsActive(win.FORWARD) {
		camera.position = camera.position.Add(camera.front.Mul(distanceMoved))
	}
	if camera.inputManager.IsActive(win.BACKWARD) {
		camera.position = camera.position.Sub(camera.front.Mul(distanceMoved))
	}
	if camera.inputManager.IsActive(win.LEFT) {
		camera.position = camera.position.Sub(camera.front.Cross(camera.up).Normalize().Mul(distanceMoved))
	}
	if camera.inputManager.IsActive(win.RIGHT) {
		camera.position = camera.position.Add(camera.front.Cross(camera.up).Normalize().Mul(distanceMoved))
	}
}

func (camera *FPSCamera) updateDirection() {
	deltaCursor := camera.inputManager.CursorChange()

	deltaX := mgl32.DegToRad(-cursorSensitivity * deltaCursor[0])
	deltaY := mgl32.DegToRad(cursorSensitivity * deltaCursor[1])

	camera.pitchRad += deltaY
	if camera.pitchRad > mgl32.DegToRad(89) {
		camera.pitchRad = math.Pi / 2
	} else if camera.pitchRad < -mgl32.DegToRad(89) {
		camera.pitchRad = -math.Pi / 2
	}

	camera.yawRad = float32(math.Mod(float64(camera.yawRad + deltaX), 360))
	camera.updateVectors()
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
