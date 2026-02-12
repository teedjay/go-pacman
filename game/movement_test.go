package game

import "testing"

func TestPacManAtTileCenter(t *testing.T) {
	p := NewPacMan()
	// Spawn position should be at tile center
	if !p.IsAtTileCenter() {
		t.Error("spawn position should be at tile center")
	}
}

func TestPacManMove(t *testing.T) {
	m := NewMaze()
	p := NewPacMan()
	p.Dir = DirLeft
	startX := p.X
	p.Move(m)
	if p.X >= startX {
		t.Error("moving left should decrease X")
	}
}

func TestPacManWallCollision(t *testing.T) {
	m := NewMaze()
	p := NewPacMan()
	// Place pacman at (1,1) facing up — row 0 is all walls
	p.X = float64(1*TileSize + TileSize/2)
	p.Y = float64(1*TileSize + TileSize/2)
	p.Dir = DirUp
	startY := p.Y
	p.Move(m)
	if p.Y != startY {
		t.Error("should not move into wall")
	}
}

func TestPacManQueuedDirection(t *testing.T) {
	m := NewMaze()
	p := NewPacMan()
	// Place at (1,1) which is a dot tile
	p.X = float64(1*TileSize + TileSize/2)
	p.Y = float64(1*TileSize + TileSize/2)
	p.Dir = DirRight
	p.NextDir = DirDown
	// If down leads to passable tile (1,2), should switch direction
	if m.IsPassable(1, 2) {
		p.Move(m)
		if p.Dir != DirDown {
			t.Error("should switch to queued direction when passable")
		}
	}
}
