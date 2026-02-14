package game

import (
	"math"
	"math/rand"
)

// BFS finds the shortest path from (startX, startY) to (targetX, targetY)
// using breadth-first search on the tile grid. Returns a slice of directions.
// Returns empty slice if no path found or start is not passable.
func BFS(m *Maze, startX, startY, targetX, targetY int) []Direction {
	if !m.IsPassableForGhost(startX, startY) {
		return nil
	}
	if startX == targetX && startY == targetY {
		return []Direction{}
	}

	type point struct{ x, y int }
	visited := make(map[point]bool)
	parent := make(map[point]point)
	dirMap := make(map[point]Direction)

	start := point{startX, startY}
	target := point{targetX, targetY}
	visited[start] = true

	queue := []point{start}
	dirs := []Direction{DirUp, DirLeft, DirDown, DirRight}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		for _, d := range dirs {
			nx, ny := nextTile(cur.x, cur.y, d)
			np := point{nx, ny}
			if visited[np] || !m.IsPassableForGhost(nx, ny) {
				continue
			}
			visited[np] = true
			parent[np] = cur
			dirMap[np] = d

			if np == target {
				// Reconstruct path
				var path []Direction
				for p := target; p != start; p = parent[p] {
					path = append([]Direction{dirMap[p]}, path...)
				}
				return path
			}
			queue = append(queue, np)
		}
	}
	return nil
}

// reverseDir returns the opposite direction.
func reverseDir(d Direction) Direction {
	switch d {
	case DirUp:
		return DirDown
	case DirDown:
		return DirUp
	case DirLeft:
		return DirRight
	case DirRight:
		return DirLeft
	}
	return DirNone
}

// ModeTimer manages the global ghost mode (scatter/chase) cycle.
type ModeTimer struct {
	phases       []modePhase
	currentPhase int
	ticksInPhase int
}

type modePhase struct {
	mode  GhostMode
	ticks int // -1 for infinite
}

// NewModeTimer creates a mode timer for the given level.
func NewModeTimer(level int) *ModeTimer {
	// Level 1 phases: scatter 7s, chase 20s, scatter 7s, chase 20s, scatter 5s, chase forever
	// Scale scatter durations down for higher levels
	scatterScale := 1.0
	if level > 1 {
		scatterScale = math.Max(0.4, 1.0-float64(level-1)*0.07)
	}

	return &ModeTimer{
		phases: []modePhase{
			{GhostScatter, int(420 * scatterScale)},
			{GhostChase, 1200},
			{GhostScatter, int(420 * scatterScale)},
			{GhostChase, 1200},
			{GhostScatter, int(300 * scatterScale)},
			{GhostChase, -1}, // forever
		},
	}
}

// Tick advances the mode timer by one tick.
func (mt *ModeTimer) Tick() {
	if mt.currentPhase >= len(mt.phases) {
		return
	}
	phase := mt.phases[mt.currentPhase]
	if phase.ticks == -1 {
		return // infinite phase, never advance
	}
	mt.ticksInPhase++
	if mt.ticksInPhase >= phase.ticks {
		mt.currentPhase++
		mt.ticksInPhase = 0
	}
}

// CurrentMode returns the current ghost mode.
func (mt *ModeTimer) CurrentMode() GhostMode {
	if mt.currentPhase >= len(mt.phases) {
		return GhostChase
	}
	return mt.phases[mt.currentPhase].mode
}

// Reset restarts the mode timer.
func (mt *ModeTimer) Reset() {
	mt.currentPhase = 0
	mt.ticksInPhase = 0
}

// ChooseDirection picks the best direction for a ghost at a tile center.
// It never reverses (unless forced). In frightened mode, picks randomly.
func (g *Ghost) ChooseDirection(m *Maze, targetX, targetY int) Direction {
	tx, ty := g.TileX(), g.TileY()
	reverse := reverseDir(g.Dir)

	// Priority order: up, left, down, right (classic Pac-Man priority)
	dirs := []Direction{DirUp, DirLeft, DirDown, DirRight}

	bestDir := DirNone
	bestDist := math.MaxFloat64

	for _, d := range dirs {
		if d == reverse {
			continue // never reverse
		}
		nx, ny := nextTile(tx, ty, d)
		if !m.IsPassableForGhost(nx, ny) {
			continue
		}
		dx := float64(nx - targetX)
		dy := float64(ny - targetY)
		dist := dx*dx + dy*dy // squared distance is fine for comparison
		if dist < bestDist {
			bestDist = dist
			bestDir = d
		}
	}
	return bestDir
}

// ChooseRandomDirection picks a random valid direction (not reverse).
func (g *Ghost) ChooseRandomDirection(m *Maze) Direction {
	tx, ty := g.TileX(), g.TileY()
	reverse := reverseDir(g.Dir)

	var valid []Direction
	for _, d := range []Direction{DirUp, DirLeft, DirDown, DirRight} {
		if d == reverse {
			continue
		}
		nx, ny := nextTile(tx, ty, d)
		if m.IsPassableForGhost(nx, ny) {
			valid = append(valid, d)
		}
	}
	if len(valid) == 0 {
		return DirNone
	}
	return valid[rand.Intn(len(valid))]
}

// isAtTileCenter checks if the ghost is within Speed pixels of the nearest tile center.
func (g *Ghost) isAtTileCenter() bool {
	centerX := float64(g.TileX()*TileSize + TileSize/2)
	centerY := float64(g.TileY()*TileSize + TileSize/2)
	return math.Abs(g.X-centerX) < g.Speed && math.Abs(g.Y-centerY) < g.Speed
}

// UpdateGhost updates a ghost's position and behavior for one tick.
func UpdateGhost(g *Ghost, m *Maze, pacman *PacMan, globalMode GhostMode) {
	// Handle ghost house exit
	if g.InHouse {
		if g.ExitTimer > 0 {
			g.ExitTimer--
			return
		}
		// Move toward ghost house exit (tile 14, 11 — just above the door)
		exitX := float64(14*TileSize + TileSize/2)
		exitY := float64(11*TileSize + TileSize/2)
		dx := exitX - g.X
		dy := exitY - g.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist < g.Speed {
			g.X = exitX
			g.Y = exitY
			g.InHouse = false
			g.Dir = DirLeft
		} else {
			g.X += (dx / dist) * g.Speed
			g.Y += (dy / dist) * g.Speed
		}
		return
	}

	// Determine effective mode
	mode := g.Mode
	if mode != GhostFrightened && mode != GhostEaten {
		mode = globalMode
	}

	if g.isAtTileCenter() {
		// Snap to tile center
		g.X = float64(g.TileX()*TileSize + TileSize/2)
		g.Y = float64(g.TileY()*TileSize + TileSize/2)

		// Choose direction based on mode
		switch mode {
		case GhostChase:
			// Target Pac-Man with small random offset to prevent clumping
			targetX := pacman.TileX() + (rand.Intn(5) - 2)
			targetY := pacman.TileY() + (rand.Intn(5) - 2)
			g.Dir = g.ChooseDirection(m, targetX, targetY)
		case GhostScatter:
			g.Dir = g.ChooseDirection(m, g.ScatterX, g.ScatterY)
		case GhostFrightened:
			g.Dir = g.ChooseRandomDirection(m)
		case GhostEaten:
			// Head back to ghost house entrance
			if g.TileX() == 14 && g.TileY() == 11 {
				// Arrived at house entrance, re-enter
				g.InHouse = true
				g.ExitTimer = 0 // exit immediately after respawn
				g.Mode = GhostScatter
				g.X = float64(14*TileSize + TileSize/2)
				g.Y = float64(14*TileSize + TileSize/2)
				return
			}
			g.Dir = g.ChooseDirection(m, 14, 11)
		}
	}

	// Move in current direction
	switch g.Dir {
	case DirUp:
		g.Y -= g.Speed
	case DirDown:
		g.Y += g.Speed
	case DirLeft:
		g.X -= g.Speed
	case DirRight:
		g.X += g.Speed
	}

	// Tunnel wrapping
	if g.X < 0 {
		g.X += float64(MazeCols * TileSize)
	} else if g.X >= float64(MazeCols*TileSize) {
		g.X -= float64(MazeCols * TileSize)
	}
}
