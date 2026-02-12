package game

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
	X, Y      float64   // pixel position (center of sprite)
	Dir       Direction  // current movement direction
	NextDir   Direction  // queued direction from input
	Speed     float64    // pixels per tick
	AnimFrame int        // 0, 1, 2 (closed, half, open)
	AnimTimer int        // ticks until next frame
	Alive     bool
}

// NewPacMan creates a new PacMan at the spawn position.
func NewPacMan() *PacMan {
	return &PacMan{
		X:     float64(PacmanSpawnX*TileSize + TileSize/2),
		Y:     float64(PacmanSpawnY*TileSize + TileSize/2),
		Dir:   DirNone,
		Speed: 1.5,
		Alive: true,
	}
}

// TileX returns current tile column.
func (p *PacMan) TileX() int { return int(p.X) / TileSize }

// TileY returns current tile row.
func (p *PacMan) TileY() int { return int(p.Y) / TileSize }
