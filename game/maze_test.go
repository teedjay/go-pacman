package game

import "testing"

func TestNewMaze(t *testing.T) {
	m := NewMaze()
	if m.Width != MazeCols {
		t.Errorf("width: got %d, want %d", m.Width, MazeCols)
	}
	if m.Height != MazeRows {
		t.Errorf("height: got %d, want %d", m.Height, MazeRows)
	}
}

func TestMazeTileTypes(t *testing.T) {
	m := NewMaze()
	if m.TileAt(0, 0) != TileWall {
		t.Error("(0,0) should be wall")
	}
	if m.TileAt(1, 1) != TileDot {
		t.Error("(1,1) should be dot")
	}
	if m.TileAt(1, 3) != TilePowerPellet {
		t.Error("(1,3) should be power pellet")
	}
}

func TestMazePassable(t *testing.T) {
	m := NewMaze()
	if !m.IsPassable(1, 1) {
		t.Error("dot tile should be passable")
	}
	if m.IsPassable(0, 0) {
		t.Error("wall should not be passable")
	}
}

func TestMazeDots(t *testing.T) {
	m := NewMaze()
	initial := m.RemainingDots()
	if initial == 0 {
		t.Fatal("expected dots in maze")
	}
	if !m.ConsumeDot(1, 1) {
		t.Error("should consume dot at (1,1)")
	}
	if m.RemainingDots() != initial-1 {
		t.Error("dot count should decrease")
	}
	if m.ConsumeDot(1, 1) {
		t.Error("should not consume dot twice")
	}
}

func TestMazeReset(t *testing.T) {
	m := NewMaze()
	initial := m.RemainingDots()
	m.ConsumeDot(1, 1)
	m.Reset()
	if m.RemainingDots() != initial {
		t.Error("reset should restore all dots")
	}
}
