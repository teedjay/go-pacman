package game

import "testing"

func TestDotScoring(t *testing.T) {
	g := New()
	// Move pacman to a known dot position (1,1)
	g.pacman.X = float64(1*TileSize + TileSize/2)
	g.pacman.Y = float64(1*TileSize + TileSize/2)
	g.checkDotConsumption()
	if g.score != 10 {
		t.Errorf("expected score 10, got %d", g.score)
	}
}

func TestPowerPelletScoring(t *testing.T) {
	g := New()
	// Move pacman to a power pellet position (1,3)
	g.pacman.X = float64(1*TileSize + TileSize/2)
	g.pacman.Y = float64(3*TileSize + TileSize/2)
	g.checkDotConsumption()
	if g.score != 50 {
		t.Errorf("expected score 50 for power pellet, got %d", g.score)
	}
}

func TestDotConsumptionRemovesDot(t *testing.T) {
	g := New()
	initial := g.maze.RemainingDots()
	g.pacman.X = float64(1*TileSize + TileSize/2)
	g.pacman.Y = float64(1*TileSize + TileSize/2)
	g.checkDotConsumption()
	if g.maze.RemainingDots() != initial-1 {
		t.Error("dot should be consumed from maze")
	}
}

func TestNoDuplicateScoring(t *testing.T) {
	g := New()
	g.pacman.X = float64(1*TileSize + TileSize/2)
	g.pacman.Y = float64(1*TileSize + TileSize/2)
	g.checkDotConsumption()
	g.checkDotConsumption() // call again on same tile
	if g.score != 10 {
		t.Errorf("should not score twice, got %d", g.score)
	}
}
