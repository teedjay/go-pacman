package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// PacManSpriteSize is the width/height of the Pac-Man sprite in pixels.
const PacManSpriteSize = 13

// GhostSpriteSize is the width/height of ghost sprites in pixels.
const GhostSpriteSize = 13

// Sprites holds pre-generated tile images for maze rendering.
type Sprites struct {
	Wall         *ebiten.Image
	Dot          *ebiten.Image
	PowerPellet  *ebiten.Image
	Empty        *ebiten.Image
	GhostDoor    *ebiten.Image
	PacManFrames [3]*ebiten.Image   // closed, half-open, full-open
	GhostSprites [4]*ebiten.Image   // one per ghost ID (Blinky, Pinky, Inky, Clyde)
	GhostFrightened *ebiten.Image   // blue frightened ghost
	GhostEyes    *ebiten.Image      // just eyes for eaten ghost
}

// sprites is the package-level sprite cache, initialized by InitSprites.
var sprites *Sprites

// InitSprites generates all tile sprites and stores them in the package-level cache.
func InitSprites() {
	ghostColors := [4]color.RGBA{
		{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}, // Blinky - red
		{R: 0xFF, G: 0xB8, B: 0xFF, A: 0xFF}, // Pinky - pink
		{R: 0x00, G: 0xFF, B: 0xFF, A: 0xFF}, // Inky - cyan
		{R: 0xFF, G: 0xB8, B: 0x52, A: 0xFF}, // Clyde - orange
	}
	var ghostSprites [4]*ebiten.Image
	for i := 0; i < 4; i++ {
		ghostSprites[i] = GenerateGhostSprite(ghostColors[i])
	}

	sprites = &Sprites{
		Wall:        GenerateWallTile(),
		Dot:         GenerateDotSprite(),
		PowerPellet: GeneratePowerPelletSprite(),
		Empty:       GenerateEmptyTile(),
		GhostDoor:   GenerateGhostDoorTile(),
		PacManFrames: [3]*ebiten.Image{
			GeneratePacManFrame(0),
			GeneratePacManFrame(1),
			GeneratePacManFrame(2),
		},
		GhostSprites:    ghostSprites,
		GhostFrightened: GenerateGhostFrightened(),
		GhostEyes:       GenerateGhostEyes(),
	}
}

// GenerateWallTile returns an 8x8 blue tile with a 1px black border.
func GenerateWallTile() *ebiten.Image {
	img := ebiten.NewImage(TileSize, TileSize)
	blue := color.RGBA{R: 0x21, G: 0x21, B: 0xDE, A: 0xFF}
	black := color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}

	for y := 0; y < TileSize; y++ {
		for x := 0; x < TileSize; x++ {
			if x == 0 || x == TileSize-1 || y == 0 || y == TileSize-1 {
				img.Set(x, y, black)
			} else {
				img.Set(x, y, blue)
			}
		}
	}
	return img
}

// GenerateDotSprite returns an 8x8 image with a 2x2 white square centered (pixels 3-4, 3-4).
func GenerateDotSprite() *ebiten.Image {
	img := ebiten.NewImage(TileSize, TileSize)
	white := color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}

	for y := 3; y <= 4; y++ {
		for x := 3; x <= 4; x++ {
			img.Set(x, y, white)
		}
	}
	return img
}

// GeneratePowerPelletSprite returns an 8x8 image with a 6x6 white square centered (pixels 1-6, 1-6).
func GeneratePowerPelletSprite() *ebiten.Image {
	img := ebiten.NewImage(TileSize, TileSize)
	white := color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}

	for y := 1; y <= 6; y++ {
		for x := 1; x <= 6; x++ {
			img.Set(x, y, white)
		}
	}
	return img
}

// GenerateEmptyTile returns an 8x8 fully black image.
func GenerateEmptyTile() *ebiten.Image {
	img := ebiten.NewImage(TileSize, TileSize)
	// ebiten.NewImage is already initialized to transparent black.
	// Fill with opaque black to be explicit.
	img.Fill(color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF})
	return img
}

// GenerateGhostDoorTile returns an 8x8 image with a pink horizontal bar in the middle (2 pixels tall, centered).
func GenerateGhostDoorTile() *ebiten.Image {
	img := ebiten.NewImage(TileSize, TileSize)
	pink := color.RGBA{R: 0xFF, G: 0xB8, B: 0xFF, A: 0xFF}

	// 2 pixels tall, centered vertically: rows 3 and 4
	for y := 3; y <= 4; y++ {
		for x := 0; x < TileSize; x++ {
			img.Set(x, y, pink)
		}
	}
	return img
}

// GeneratePacManFrame generates a 13x13 Pac-Man sprite frame.
// frame 0: closed mouth (full circle)
// frame 1: half-open mouth (~30 degrees)
// frame 2: full-open mouth (~60 degrees)
// The mouth faces right by default.
func GeneratePacManFrame(frame int) *ebiten.Image {
	const size = PacManSpriteSize
	const cx, cy = 6, 6 // center
	const radius = 6

	img := ebiten.NewImage(size, size)
	yellow := color.RGBA{R: 0xFF, G: 0xFF, B: 0x00, A: 0xFF}

	// Determine mouth half-angle in radians based on frame.
	var mouthHalfAngle float64
	switch frame {
	case 0:
		mouthHalfAngle = 0 // closed
	case 1:
		mouthHalfAngle = 15.0 * math.Pi / 180.0 // ~30 degrees total
	case 2:
		mouthHalfAngle = 30.0 * math.Pi / 180.0 // ~60 degrees total
	}

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := x - cx
			dy := y - cy

			// Check if pixel is within the circle.
			if dx*dx+dy*dy > radius*radius {
				continue
			}

			// For frame 0 (closed), draw all pixels in the circle.
			if frame == 0 {
				img.Set(x, y, yellow)
				continue
			}

			// Calculate angle from center to this pixel.
			// atan2 returns angle in [-pi, pi], with 0 pointing right.
			angle := math.Atan2(float64(dy), float64(dx))

			// Exclude pixels within the mouth angle range (mouth faces right, centered on angle 0).
			if angle > -mouthHalfAngle && angle < mouthHalfAngle && dx > 0 {
				continue // inside the mouth, skip this pixel
			}

			img.Set(x, y, yellow)
		}
	}

	return img
}

// GenerateGhostSprite generates a 13x13 ghost sprite with the given body color.
// Shape: rounded top (semicircle), flat sides, wavy bottom (3 bumps), with eyes.
func GenerateGhostSprite(bodyColor color.RGBA) *ebiten.Image {
	const size = GhostSpriteSize
	img := ebiten.NewImage(size, size)

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			// Top half: semicircle (center at 6, 6, radius 6)
			if y <= 6 {
				dx := x - 6
				dy := y - 6
				if dx*dx+dy*dy <= 36 {
					img.Set(x, y, bodyColor)
				}
			} else if y < 11 {
				// Middle: rectangular body from x=0 to x=12
				if x >= 0 && x <= 12 {
					img.Set(x, y, bodyColor)
				}
			} else {
				// Bottom: wavy edge with 3 bumps
				bump := false
				for _, cx := range []int{2, 6, 10} {
					dx := x - cx
					dy := y - 10
					if dx*dx+dy*dy <= 4 {
						bump = true
						break
					}
				}
				if bump {
					img.Set(x, y, bodyColor)
				}
			}
		}
	}

	drawGhostEyesOn(img)
	return img
}

// drawGhostEyesOn draws two eyes with pupils on the given image.
func drawGhostEyesOn(img *ebiten.Image) {
	white := color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	blue := color.RGBA{R: 0x21, G: 0x21, B: 0xFF, A: 0xFF}

	for _, ex := range []int{4, 8} {
		for y := 0; y < GhostSpriteSize; y++ {
			for x := 0; x < GhostSpriteSize; x++ {
				dx := x - ex
				dy := y - 4
				if dx*dx+dy*dy <= 4 {
					img.Set(x, y, white)
				}
			}
		}
		// Pupil: blue dot (looking forward/right)
		img.Set(ex+1, 4, blue)
		img.Set(ex+1, 5, blue)
	}
}

// GenerateGhostFrightened generates the blue frightened ghost sprite.
func GenerateGhostFrightened() *ebiten.Image {
	blue := color.RGBA{R: 0x21, G: 0x21, B: 0xFF, A: 0xFF}
	img := GenerateGhostSprite(blue)

	white := color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	img.Set(4, 4, white)
	img.Set(8, 4, white)

	// Wavy mouth
	for x := 2; x <= 10; x++ {
		yOff := 8
		if x%2 == 0 {
			yOff = 9
		}
		img.Set(x, yOff, white)
	}
	return img
}

// GenerateGhostEyes generates just the eyes sprite for eaten ghosts.
func GenerateGhostEyes() *ebiten.Image {
	img := ebiten.NewImage(GhostSpriteSize, GhostSpriteSize)
	drawGhostEyesOn(img)
	return img
}
