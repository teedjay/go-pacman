package game

import "testing"

func TestBFS(t *testing.T) {
	m := NewMaze()
	// Find path from (1,1) to (1,5) — should find a path
	path := BFS(m, 1, 1, 1, 5)
	if len(path) == 0 {
		t.Error("BFS should find a path")
	}
}

func TestBFSBlocked(t *testing.T) {
	m := NewMaze()
	// Path from wall to passable should return empty
	path := BFS(m, 0, 0, 1, 1)
	if len(path) != 0 {
		t.Error("BFS from wall should return empty path")
	}
}

func TestGhostModeTimer(t *testing.T) {
	mt := NewModeTimer(1) // level 1
	if mt.CurrentMode() != GhostScatter {
		t.Error("should start in scatter mode")
	}
	// Advance past first scatter phase (7 seconds = 420 ticks at 60 TPS)
	for i := 0; i < 421; i++ {
		mt.Tick()
	}
	if mt.CurrentMode() != GhostChase {
		t.Error("should switch to chase after scatter")
	}
}

func TestGhostModeTimerFullCycle(t *testing.T) {
	mt := NewModeTimer(1)
	// Scatter 420 -> Chase 1200 -> Scatter 420 -> Chase 1200 -> Scatter 300 -> Chase forever
	// After all scatter+chase phases, should be in chase permanently
	total := 420 + 1200 + 420 + 1200 + 300 + 1
	for i := 0; i < total; i++ {
		mt.Tick()
	}
	if mt.CurrentMode() != GhostChase {
		t.Errorf("should be in permanent chase, got %d", mt.CurrentMode())
	}
}

func TestGhostChooseDirection(t *testing.T) {
	m := NewMaze()
	g := &Ghost{
		X:   float64(1*TileSize + TileSize/2),
		Y:   float64(5*TileSize + TileSize/2),
		Dir: DirRight,
	}
	// Ghost should choose a valid direction toward target
	dir := g.ChooseDirection(m, 26, 5)
	if dir == DirNone {
		t.Error("ghost should choose a direction")
	}
	// Ghost should not reverse direction
	if dir == DirLeft {
		t.Error("ghost should not reverse direction")
	}
}
