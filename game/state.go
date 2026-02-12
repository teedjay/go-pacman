package game

// GameState represents the current phase of the game.
type GameState int

const (
	StateTitle GameState = iota
	StateReady
	StatePlaying
	StateDeath
	StateLevelClear
	StateGameOver
)

// validTransitions defines which state transitions are allowed.
var validTransitions = map[GameState][]GameState{
	StateTitle:      {StateReady},
	StateReady:      {StatePlaying},
	StatePlaying:    {StateDeath, StateLevelClear},
	StateDeath:      {StatePlaying, StateGameOver},
	StateLevelClear: {StateReady},
	StateGameOver:   {StateTitle},
}

// isValidTransition returns true if the transition from -> to is allowed.
func isValidTransition(from, to GameState) bool {
	targets, ok := validTransitions[from]
	if !ok {
		return false
	}
	for _, t := range targets {
		if t == to {
			return true
		}
	}
	return false
}
