package game

import "testing"

func TestStateTransitions(t *testing.T) {
	tests := []struct {
		from, to GameState
		event    string
	}{
		{StateTitle, StateReady, "start"},
		{StateReady, StatePlaying, "ready_timeout"},
		{StatePlaying, StateDeath, "pacman_dies"},
		{StateDeath, StatePlaying, "respawn"},
		{StateDeath, StateGameOver, "no_lives"},
		{StatePlaying, StateLevelClear, "all_dots_eaten"},
		{StateLevelClear, StateReady, "next_level"},
		{StateGameOver, StateTitle, "continue"},
	}
	for _, tt := range tests {
		t.Run(tt.event, func(t *testing.T) {
			if !isValidTransition(tt.from, tt.to) {
				t.Errorf("transition %d -> %d should be valid", tt.from, tt.to)
			}
		})
	}
}

func TestInvalidTransitions(t *testing.T) {
	if isValidTransition(StateTitle, StatePlaying) {
		t.Error("Title -> Playing should be invalid (must go through Ready)")
	}
	if isValidTransition(StatePlaying, StateTitle) {
		t.Error("Playing -> Title should be invalid")
	}
}
