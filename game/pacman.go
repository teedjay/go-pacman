package game

import "math"

// Direction represents a movement direction.
type Direction int

const (
	DirNone Direction = iota
	DirUp
	DirDown
	DirLeft
	DirRight
)

// PacMan represents the player-controlled Pac-Man entity.
type PacMan struct {
	X, Y       float64   // pixel position (center of sprite)
	Dir        Direction  // current movement direction
	NextDir    Direction  // queued direction from input
	Speed      float64    // pixels per tick
	AnimFrame  int        // 0, 1, 2 (closed, half, open)
	AnimTimer  int        // ticks until next frame
	Alive      bool
	DeathFrame int        // current death animation frame (0-10)
	DeathTimer int        // ticks remaining in death animation

	lastCenterTX int // tile where last center processing happened
	lastCenterTY int
}

// NewPacMan creates a new PacMan at the spawn position.
func NewPacMan() *PacMan {
	return &PacMan{
		X:              float64(PacmanSpawnX*TileSize + TileSize/2),
		Y:              float64(PacmanSpawnY*TileSize + TileSize/2),
		Dir:            DirNone,
		Speed:          1.5,
		Alive:          true,
		lastCenterTX:   -1,
		lastCenterTY:   -1,
	}
}

// TileX returns current tile column.
func (p *PacMan) TileX() int { return int(p.X) / TileSize }

// TileY returns current tile row.
func (p *PacMan) TileY() int { return int(p.Y) / TileSize }

// IsAtTileCenter returns true if position is within Speed pixels of the nearest tile center.
func (p *PacMan) IsAtTileCenter() bool {
	centerX := float64(p.TileX()*TileSize + TileSize/2)
	centerY := float64(p.TileY()*TileSize + TileSize/2)
	return math.Abs(p.X-centerX) < p.Speed && math.Abs(p.Y-centerY) < p.Speed
}

// Move updates Pac-Man's position based on current direction and maze walls.
func (p *PacMan) Move(m *Maze) {
	if p.IsAtTileCenter() {
		tileX, tileY := p.TileX(), p.TileY()
		if p.Dir != DirNone && tileX == p.lastCenterTX && tileY == p.lastCenterTY {
			goto pacMove // already processed this tile center while moving
		}
		p.lastCenterTX = tileX
		p.lastCenterTY = tileY

		// Snap to exact tile center.
		p.X = float64(tileX*TileSize + TileSize/2)
		p.Y = float64(tileY*TileSize + TileSize/2)

		// Check if NextDir leads to a passable tile; if so, switch.
		if p.NextDir != DirNone {
			nx, ny := nextTile(tileX, tileY, p.NextDir)
			if m.IsPassable(nx, ny) {
				p.Dir = p.NextDir
				p.NextDir = DirNone
			}
		}

		// Check if current Dir leads to a passable tile; if not, stop.
		if p.Dir != DirNone {
			nx, ny := nextTile(tileX, tileY, p.Dir)
			if !m.IsPassable(nx, ny) {
				p.Dir = DirNone
			}
		}
	}
pacMove:

	// Advance position based on direction.
	switch p.Dir {
	case DirUp:
		p.Y -= p.Speed
	case DirDown:
		p.Y += p.Speed
	case DirLeft:
		p.X -= p.Speed
	case DirRight:
		p.X += p.Speed
	}

	// Tunnel wrapping
	if p.X < 0 {
		p.X += float64(MazeCols * TileSize)
	} else if p.X >= float64(MazeCols*TileSize) {
		p.X -= float64(MazeCols * TileSize)
	}

	// Advance animation: cycle through frames every 4 ticks.
	// AnimTimer counts total ticks; divide by 4 to get cycle position.
	// Cycle: 0→1→2→1→0→1→2→1... (4-step cycle mapping to frames [0,1,2,1])
	if p.Dir != DirNone {
		p.AnimTimer++
		cycle := [4]int{0, 1, 2, 1}
		p.AnimFrame = cycle[(p.AnimTimer/4)%4]
	}
}

// nextTile returns the tile coordinates one step in the given direction.
func nextTile(x, y int, dir Direction) (int, int) {
	switch dir {
	case DirUp:
		return x, y - 1
	case DirDown:
		return x, y + 1
	case DirLeft:
		return x - 1, y
	case DirRight:
		return x + 1, y
	}
	return x, y
}
