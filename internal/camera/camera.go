package camera

import "math"

type Camera struct {
	// CameraX, Y is the x position of the camera
	CameraX, CameraY float64
}

// NewCamera creates a new camera
func NewCamera(x, y float64) *Camera {
	return &Camera{
		CameraX: x,
		CameraY: y,
	}
}

// FollowPlayer makes the camera follow the player
func (c *Camera) FollowPlayer(playerX, playerY, screenWidth, screenHeight float64) {
	c.CameraX = -playerX + screenWidth/2
	c.CameraY = -playerY + screenHeight/2
}

// constrain restricts the camera to the map
func (c *Camera) Constrain(screenWidth, screenHeight, mapWidth, mapHeight float64) {
	c.CameraX = math.Min(c.CameraX, 0.0)
	c.CameraY = math.Min(c.CameraY, 0.0)

	c.CameraX = math.Max(c.CameraX, screenWidth-mapWidth)
	c.CameraY = math.Max(c.CameraY, screenHeight-mapHeight)
}
