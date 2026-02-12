package game

// GhostMode represents the current behavior mode of a ghost.
type GhostMode int

const (
	GhostChase GhostMode = iota
	GhostScatter
	GhostFrightened
	GhostEaten
)

// GhostID identifies each of the four ghosts.
type GhostID int

const (
	Blinky GhostID = iota // red
	Pinky                 // pink
	Inky                  // cyan
	Clyde                 // orange
)

// Ghost represents a ghost entity in the game.
type Ghost struct {
	ID        GhostID
	X, Y      float64
	Dir       Direction
	Mode      GhostMode
	Speed     float64
	SpawnX    int // tile coords for initial position
	SpawnY    int
	ScatterX  int // scatter corner target tile
	ScatterY  int
	InHouse   bool // still in ghost house
	ExitTimer int  // ticks until ghost leaves house
}

// NewGhosts creates all four ghosts at their starting positions.
func NewGhosts() [4]*Ghost {
	return [4]*Ghost{
		{
			ID: Blinky, Speed: 1.3,
			X: float64(14*TileSize + TileSize/2), Y: float64(11*TileSize + TileSize/2),
			SpawnX: 14, SpawnY: 11,
			ScatterX: 25, ScatterY: 0,
			Dir: DirLeft, InHouse: false, ExitTimer: 0,
		},
		{
			ID: Pinky, Speed: 1.3,
			X: float64(12*TileSize + TileSize/2), Y: float64(14*TileSize + TileSize/2),
			SpawnX: 12, SpawnY: 14,
			ScatterX: 2, ScatterY: 0,
			Dir: DirDown, InHouse: true, ExitTimer: 0, // exits immediately
		},
		{
			ID: Inky, Speed: 1.3,
			X: float64(14*TileSize + TileSize/2), Y: float64(14*TileSize + TileSize/2),
			SpawnX: 14, SpawnY: 14,
			ScatterX: 27, ScatterY: 30,
			Dir: DirUp, InHouse: true, ExitTimer: 300, // exits after 5s
		},
		{
			ID: Clyde, Speed: 1.3,
			X: float64(16*TileSize + TileSize/2), Y: float64(14*TileSize + TileSize/2),
			SpawnX: 16, SpawnY: 14,
			ScatterX: 0, ScatterY: 30,
			Dir: DirUp, InHouse: true, ExitTimer: 600, // exits after 10s
		},
	}
}

// TileX returns the ghost's current tile column.
func (g *Ghost) TileX() int { return int(g.X) / TileSize }

// TileY returns the ghost's current tile row.
func (g *Ghost) TileY() int { return int(g.Y) / TileSize }

// ResetToSpawn returns the ghost to its initial spawn position.
func (g *Ghost) ResetToSpawn() {
	g.X = float64(g.SpawnX*TileSize + TileSize/2)
	g.Y = float64(g.SpawnY*TileSize + TileSize/2)
	if g.ID == Blinky {
		g.InHouse = false
		g.Dir = DirLeft
	} else {
		g.InHouse = true
		g.Dir = DirUp
	}
	g.Mode = GhostScatter
}
