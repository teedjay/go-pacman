package game

// DifficultyParams holds all level-dependent game parameters.
type DifficultyParams struct {
	PacManSpeed     float64
	GhostSpeed      float64
	FrightenedSpeed float64
	FrightenedTicks int // at 60 TPS
	EatenSpeed      float64
}

// GetDifficulty returns interpolated difficulty parameters for the given level.
// Parameters are linearly interpolated between level 1 and level 10, clamped at 10.
func GetDifficulty(level int) DifficultyParams {
	t := float64(level-1) / 9.0
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}

	return DifficultyParams{
		PacManSpeed:     lerp(1.5, 1.8, t),
		GhostSpeed:      lerp(1.3, 1.8, t),
		FrightenedSpeed: 0.8,
		FrightenedTicks: int(lerp(360, 60, t)),
		EatenSpeed:      3.0,
	}
}

func lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}
