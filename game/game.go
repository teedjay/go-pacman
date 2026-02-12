package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	TileSize     = 8
	MazeCols     = 28
	MazeRows     = 31
	HUDTopRows   = 3
	HUDBotRows   = 2
	ScreenWidth  = MazeCols * TileSize                              // 224
	ScreenHeight = (MazeRows + HUDTopRows + HUDBotRows) * TileSize // 288
	Scale        = 3
)

type Game struct {
	maze   *Maze
	pacman *PacMan

	score     int
	highScore int
	lives     int
	level     int
}

func New() *Game {
	InitSprites()
	return &Game{
		maze:   NewMaze(),
		pacman: NewPacMan(),
		lives:  3,
		level:  1,
	}
}

func (g *Game) Update() error {
	ReadInput(g.pacman)
	g.pacman.Move(g.maze)
	g.checkDotConsumption()
	return nil
}

// checkDotConsumption checks if Pac-Man is on a dot or power pellet and consumes it.
func (g *Game) checkDotConsumption() {
	tx, ty := g.pacman.TileX(), g.pacman.TileY()
	tile := g.maze.TileAt(tx, ty)
	if tile == TileDot {
		g.maze.ConsumeDot(tx, ty)
		g.score += 10
	} else if tile == TilePowerPellet {
		g.maze.ConsumeDot(tx, ty)
		g.score += 50
		// TODO: trigger frightened mode (Task 9)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw maze tiles.
	for y := 0; y < MazeRows; y++ {
		for x := 0; x < MazeCols; x++ {
			var tile *ebiten.Image
			switch g.maze.TileAt(x, y) {
			case TileWall:
				tile = sprites.Wall
			case TileDot:
				tile = sprites.Dot
			case TilePowerPellet:
				tile = sprites.PowerPellet
			case TileEmpty:
				tile = sprites.Empty
			case TileGhostHouse:
				tile = sprites.Empty
			case TileGhostDoor:
				tile = sprites.GhostDoor
			default:
				tile = sprites.Empty
			}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*TileSize), float64(y*TileSize+HUDTopRows*TileSize))
			screen.DrawImage(tile, op)
		}
	}

	// Draw Pac-Man.
	g.drawPacMan(screen)
}

// drawPacMan draws the Pac-Man sprite with appropriate rotation/flip for its direction.
func (g *Game) drawPacMan(screen *ebiten.Image) {
	p := g.pacman
	frame := sprites.PacManFrames[p.AnimFrame]

	op := &ebiten.DrawImageOptions{}

	// Step 1: Center the sprite at origin (center pixel is at 6,6 in 13x13 image).
	op.GeoM.Translate(-6, -6)

	// Step 2: Apply rotation/flip based on direction.
	switch p.Dir {
	case DirLeft:
		// Flip horizontally: scale X by -1.
		op.GeoM.Scale(-1, 1)
	case DirUp:
		// Rotate -90 degrees (counter-clockwise).
		op.GeoM.Rotate(-math.Pi / 2)
	case DirDown:
		// Rotate +90 degrees (clockwise).
		op.GeoM.Rotate(math.Pi / 2)
	case DirRight, DirNone:
		// No transform needed; mouth faces right by default.
	}

	// Step 3: Translate to Pac-Man's pixel position + HUD offset.
	op.GeoM.Translate(p.X, p.Y+float64(HUDTopRows*TileSize))

	screen.DrawImage(frame, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
