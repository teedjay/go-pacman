package game

import "testing"

func TestDifficultyLevel1(t *testing.T) {
	d := GetDifficulty(1)
	if d.PacManSpeed != 1.5 {
		t.Errorf("level 1 pac-man speed: got %f, want 1.5", d.PacManSpeed)
	}
	if d.GhostSpeed != 1.3 {
		t.Errorf("level 1 ghost speed: got %f, want 1.3", d.GhostSpeed)
	}
	if d.FrightenedTicks != 360 {
		t.Errorf("level 1 frightened ticks: got %d, want 360", d.FrightenedTicks)
	}
}

func TestDifficultyLevel10(t *testing.T) {
	d := GetDifficulty(10)
	if d.PacManSpeed != 1.8 {
		t.Errorf("level 10 pac-man speed: got %f, want 1.8", d.PacManSpeed)
	}
	if d.FrightenedTicks != 60 {
		t.Errorf("level 10 frightened ticks: got %d, want 60", d.FrightenedTicks)
	}
}

func TestDifficultyScalesLinearly(t *testing.T) {
	d1 := GetDifficulty(1)
	d3 := GetDifficulty(3)
	d5 := GetDifficulty(5)
	if d3.PacManSpeed <= d1.PacManSpeed {
		t.Error("speed should increase with level")
	}
	if d5.PacManSpeed <= d3.PacManSpeed {
		t.Error("speed should increase with level")
	}
}

func TestDifficultyCapsAtLevel10(t *testing.T) {
	d10 := GetDifficulty(10)
	d20 := GetDifficulty(20)
	if d10.PacManSpeed != d20.PacManSpeed {
		t.Error("difficulty should cap at level 10")
	}
}
