package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Sprites holds pre-generated tile images for maze rendering.
type Sprites struct {
	Wall        *ebiten.Image
	Dot         *ebiten.Image
	PowerPellet *ebiten.Image
	Empty       *ebiten.Image
	GhostDoor   *ebiten.Image
}

// sprites is the package-level sprite cache, initialized by InitSprites.
var sprites *Sprites

// InitSprites generates all tile sprites and stores them in the package-level cache.
func InitSprites() {
	sprites = &Sprites{
		Wall:        GenerateWallTile(),
		Dot:         GenerateDotSprite(),
		PowerPellet: GeneratePowerPelletSprite(),
		Empty:       GenerateEmptyTile(),
		GhostDoor:   GenerateGhostDoorTile(),
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
