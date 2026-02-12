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
	maze      *Maze
	pacman    *PacMan
	ghosts    [4]*Ghost
	modeTimer *ModeTimer

	score            int
	highScore        int
	lives            int
	level            int
	ghostsEatenCombo int // resets each power pellet
	frightenedTimer  int // ticks remaining for frightened mode
	deathTimer       int // ticks remaining for death pause/respawn
}

func New() *Game {
	InitSprites()
	return &Game{
		maze:      NewMaze(),
		pacman:    NewPacMan(),
		ghosts:    NewGhosts(),
		modeTimer: NewModeTimer(1),
		lives:     3,
		level:     1,
	}
}

func (g *Game) Update() error {
	// Handle death pause
	if !g.pacman.Alive {
		g.deathTimer--
		if g.deathTimer <= 0 {
			g.respawn()
		}
		return nil
	}

	ReadInput(g.pacman)
	g.pacman.Move(g.maze)
	g.checkDotConsumption()

	// Update frightened timer
	if g.frightenedTimer > 0 {
		g.frightenedTimer--
		if g.frightenedTimer == 0 {
			for _, ghost := range g.ghosts {
				if ghost.Mode == GhostFrightened {
					ghost.Mode = GhostChase // return to normal
				}
			}
		}
	}

	// Update ghost AI
	g.modeTimer.Tick()
	globalMode := g.modeTimer.CurrentMode()
	for _, ghost := range g.ghosts {
		UpdateGhost(ghost, g.maze, g.pacman, globalMode)
	}

	// Check collisions
	g.checkGhostCollisions()

	return nil
}

// respawn resets Pac-Man and ghosts after a death.
func (g *Game) respawn() {
	g.pacman = NewPacMan()
	g.ghosts = NewGhosts()
	g.modeTimer.Reset()
	g.frightenedTimer = 0
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
		g.triggerFrightenedMode()
	}
}

// triggerFrightenedMode sets all non-eaten ghosts to frightened and reverses their direction.
func (g *Game) triggerFrightenedMode() {
	g.ghostsEatenCombo = 0
	g.frightenedTimer = 360 // 6 seconds at 60 TPS (will be adjusted by difficulty later)
	for _, ghost := range g.ghosts {
		if ghost.Mode != GhostEaten && !ghost.InHouse {
			ghost.Mode = GhostFrightened
			ghost.Dir = reverseDir(ghost.Dir)
		}
	}
}

// CheckCollision returns true if Pac-Man and a ghost are within 6 pixels of each other.
func CheckCollision(p *PacMan, gh *Ghost) bool {
	dx := p.X - gh.X
	dy := p.Y - gh.Y
	dist := math.Sqrt(dx*dx + dy*dy)
	return dist < 6
}

// ghostEatScore returns the score for eating the next ghost in the combo.
func (g *Game) ghostEatScore() int {
	// 200, 400, 800, 1600
	return 200 << g.ghostsEatenCombo
}

// checkGhostCollisions checks for Pac-Man colliding with any ghost.
func (g *Game) checkGhostCollisions() {
	for _, ghost := range g.ghosts {
		if ghost.InHouse || ghost.Mode == GhostEaten {
			continue
		}
		if !CheckCollision(g.pacman, ghost) {
			continue
		}
		if ghost.Mode == GhostFrightened {
			// Eat the ghost
			g.score += g.ghostEatScore()
			g.ghostsEatenCombo++
			ghost.Mode = GhostEaten
		} else {
			// Pac-Man dies
			g.pacman.Alive = false
			g.lives--
			g.deathTimer = 120 // pause before respawn
			return
		}
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

	// Draw ghosts.
	for _, ghost := range g.ghosts {
		g.drawGhost(screen, ghost)
	}

	// Draw Pac-Man.
	g.drawPacMan(screen)
}

// drawGhost draws a ghost sprite based on its current mode.
func (g *Game) drawGhost(screen *ebiten.Image, ghost *Ghost) {
	var sprite *ebiten.Image
	switch ghost.Mode {
	case GhostFrightened:
		sprite = sprites.GhostFrightened
	case GhostEaten:
		sprite = sprites.GhostEyes
	default:
		sprite = sprites.GhostSprites[ghost.ID]
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-6, -6)
	op.GeoM.Translate(ghost.X, ghost.Y+float64(HUDTopRows*TileSize))
	screen.DrawImage(sprite, op)
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
