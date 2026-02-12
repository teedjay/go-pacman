package game

import "testing"

func TestCollisionDetection(t *testing.T) {
	p := NewPacMan()
	gh := &Ghost{X: p.X, Y: p.Y} // same position
	if !CheckCollision(p, gh) {
		t.Error("overlapping positions should collide")
	}
	gh.X = p.X + 20 // far away
	if CheckCollision(p, gh) {
		t.Error("distant positions should not collide")
	}
}

func TestCollisionThreshold(t *testing.T) {
	p := NewPacMan()
	gh := &Ghost{X: p.X + 5, Y: p.Y} // close but within threshold
	if !CheckCollision(p, gh) {
		t.Error("positions within threshold (6px) should collide")
	}
	gh.X = p.X + 7 // just outside threshold
	if CheckCollision(p, gh) {
		t.Error("positions outside threshold should not collide")
	}
}

func TestGhostEatingScore(t *testing.T) {
	g := New()
	g.ghostsEatenCombo = 0
	score := g.ghostEatScore()
	if score != 200 {
		t.Errorf("first ghost should give 200, got %d", score)
	}
	g.ghostsEatenCombo = 1
	if g.ghostEatScore() != 400 {
		t.Errorf("second ghost should give 400, got %d", g.ghostEatScore())
	}
	g.ghostsEatenCombo = 2
	if g.ghostEatScore() != 800 {
		t.Errorf("third ghost should give 800, got %d", g.ghostEatScore())
	}
	g.ghostsEatenCombo = 3
	if g.ghostEatScore() != 1600 {
		t.Errorf("fourth ghost should give 1600, got %d", g.ghostEatScore())
	}
}
