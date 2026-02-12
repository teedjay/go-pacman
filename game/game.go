package game

import (
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
	maze *Maze
}

func New() *Game {
	InitSprites()
	return &Game{
		maze: NewMaze(),
	}
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
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
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
