package game

import (
	"image/color"
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

	sound *SoundManager

	state      GameState
	stateTimer int
	tickCount  int

	score            int
	highScore        int
	lives            int
	level            int
	ghostsEatenCombo int // resets each power pellet
	frightenedTimer  int // ticks remaining for frightened mode
}

func New() *Game {
	InitSprites()
	sm := NewSoundManager()
	return &Game{
		sound:     sm,
		maze:      NewMaze(),
		pacman:    NewPacMan(),
		ghosts:    NewGhosts(),
		modeTimer: NewModeTimer(1),
		state:     StateTitle,
		lives:     3,
		level:     1,
	}
}

func (g *Game) Update() error {
	g.tickCount++

	switch g.state {
	case StateTitle:
		g.updateTitle()
	case StateReady:
		g.updateReady()
	case StatePlaying:
		g.updatePlaying()
	case StateDeath:
		g.updateDeath()
	case StateLevelClear:
		g.updateLevelClear()
	case StateGameOver:
		g.updateGameOver()
	}
	return nil
}

func (g *Game) updateTitle() {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.score = 0
		g.lives = 3
		g.level = 1
		g.maze.Reset()
		g.pacman = NewPacMan()
		g.ghosts = NewGhosts()
		g.modeTimer = NewModeTimer(1)
		g.frightenedTimer = 0
		g.applyDifficulty()
		g.state = StateReady
		g.stateTimer = 120 // 2 seconds
	}
}

// applyDifficulty sets speeds and timers based on current level.
func (g *Game) applyDifficulty() {
	d := GetDifficulty(g.level)
	g.pacman.Speed = d.PacManSpeed
	for _, ghost := range g.ghosts {
		ghost.Speed = d.GhostSpeed
	}
}

func (g *Game) updateReady() {
	g.stateTimer--
	if g.stateTimer <= 0 {
		g.state = StatePlaying
	}
}

func (g *Game) updatePlaying() {
	ReadInput(g.pacman)
	g.pacman.Move(g.maze)
	g.checkDotConsumption()

	// Update frightened timer
	if g.frightenedTimer > 0 {
		g.frightenedTimer--
		if g.frightenedTimer == 0 {
			for _, ghost := range g.ghosts {
				if ghost.Mode == GhostFrightened {
					ghost.Mode = GhostChase
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

	// Check level clear
	if g.maze.RemainingDots() == 0 {
		g.state = StateLevelClear
		g.stateTimer = 120 // 2 seconds of flashing
		g.sound.PlayLevelClear()
	}
}

func (g *Game) updateDeath() {
	g.stateTimer--
	if g.stateTimer <= 0 {
		if g.lives <= 0 {
			if g.score > g.highScore {
				g.highScore = g.score
			}
			g.state = StateGameOver
			g.stateTimer = 180 // 3 seconds
		} else {
			g.pacman = NewPacMan()
			g.ghosts = NewGhosts()
			g.modeTimer.Reset()
			g.frightenedTimer = 0
			g.state = StateReady
			g.stateTimer = 120
		}
	}
}

func (g *Game) updateLevelClear() {
	g.stateTimer--
	if g.stateTimer <= 0 {
		g.level++
		g.maze.Reset()
		g.pacman = NewPacMan()
		g.ghosts = NewGhosts()
		g.modeTimer = NewModeTimer(g.level)
		g.frightenedTimer = 0
		g.applyDifficulty()
		g.state = StateReady
		g.stateTimer = 120
	}
}

func (g *Game) updateGameOver() {
	g.stateTimer--
	if g.stateTimer <= 0 {
		g.state = StateTitle
	}
}

// checkDotConsumption checks if Pac-Man is on a dot or power pellet and consumes it.
func (g *Game) checkDotConsumption() {
	tx, ty := g.pacman.TileX(), g.pacman.TileY()
	tile := g.maze.TileAt(tx, ty)
	if tile == TileDot {
		g.maze.ConsumeDot(tx, ty)
		g.score += 10
		g.sound.PlayChomp()
	} else if tile == TilePowerPellet {
		g.maze.ConsumeDot(tx, ty)
		g.score += 50
		g.sound.PlayPowerUp()
		g.triggerFrightenedMode()
	}
}

// triggerFrightenedMode sets all non-eaten ghosts to frightened and reverses their direction.
func (g *Game) triggerFrightenedMode() {
	g.ghostsEatenCombo = 0
	g.frightenedTimer = GetDifficulty(g.level).FrightenedTicks
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
			g.score += g.ghostEatScore()
			g.ghostsEatenCombo++
			ghost.Mode = GhostEaten
			g.sound.PlayGhostEaten()
		} else {
			// Pac-Man dies
			g.pacman.Alive = false
			g.lives--
			g.state = StateDeath
			g.stateTimer = 120
			g.sound.PlayDeath()
			return
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	white := color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}

	switch g.state {
	case StateTitle:
		DrawText(screen, "GO PAC-MAN", 62, 100, white)
		DrawText(screen, "PRESS SPACE", 56, 140, white)
		DrawText(screen, "TO START", 68, 155, white)
		return

	case StateGameOver:
		g.drawMaze(screen)
		DrawText(screen, "GAME OVER", 65, 160, white)
		DrawHUD(screen, g.score, g.highScore, g.lives, g.level)
		return
	}

	// All other states draw the maze and entities
	g.drawMaze(screen)

	// Draw ghosts (not during death)
	if g.state != StateDeath {
		for _, ghost := range g.ghosts {
			g.drawGhost(screen, ghost)
		}
	}

	// Draw Pac-Man (not during death after animation would play)
	if g.pacman.Alive {
		g.drawPacMan(screen)
	}

	// Draw HUD
	DrawHUD(screen, g.score, g.highScore, g.lives, g.level)

	// State-specific overlays
	switch g.state {
	case StateReady:
		DrawText(screen, "READY!", 85, 164, color.RGBA{R: 0xFF, G: 0xFF, B: 0x00, A: 0xFF})
	case StateLevelClear:
		// Flash walls: alternate white/blue every 15 ticks
		// (handled in drawMaze via tickCount)
	}
}

// drawMaze draws all maze tiles.
func (g *Game) drawMaze(screen *ebiten.Image) {
	for y := 0; y < MazeRows; y++ {
		for x := 0; x < MazeCols; x++ {
			var tile *ebiten.Image
			switch g.maze.TileAt(x, y) {
			case TileWall:
				// Flash walls during level clear
				if g.state == StateLevelClear && (g.stateTimer/15)%2 == 1 {
					tile = sprites.Empty // flash to black
				} else {
					tile = sprites.Wall
				}
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
	op.GeoM.Translate(-6, -6)

	switch p.Dir {
	case DirLeft:
		op.GeoM.Scale(-1, 1)
	case DirUp:
		op.GeoM.Rotate(-math.Pi / 2)
	case DirDown:
		op.GeoM.Rotate(math.Pi / 2)
	case DirRight, DirNone:
	}

	op.GeoM.Translate(p.X, p.Y+float64(HUDTopRows*TileSize))
	screen.DrawImage(frame, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
