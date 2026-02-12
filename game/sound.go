package game

import (
	"bytes"
	"encoding/binary"
	"math"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const sampleRate = 44100

// SoundManager handles all game sound effects using programmatic audio.
type SoundManager struct {
	context    *audio.Context
	chomp1     []byte
	chomp2     []byte
	chompFlip  bool
	powerUp    []byte
	ghostEaten []byte
	death      []byte
	levelClear []byte
}

var audioContext *audio.Context

// NewSoundManager creates and initializes all sound effect buffers.
func NewSoundManager() *SoundManager {
	if audioContext == nil {
		audioContext = audio.NewContext(sampleRate)
	}
	sm := &SoundManager{
		context: audioContext,
	}
	sm.chomp1 = generateSquareWave(260, 0.06, sampleRate)
	sm.chomp2 = generateSquareWave(390, 0.06, sampleRate)
	sm.powerUp = generateSweep(800, 200, 0.3, sampleRate)
	sm.ghostEaten = generateSweep(200, 1200, 0.15, sampleRate)
	sm.death = generateDeathSound(sampleRate)
	sm.levelClear = generateArpeggio(sampleRate)
	return sm
}

// PlayChomp plays the alternating dot-eating sound.
func (sm *SoundManager) PlayChomp() {
	if sm.chompFlip {
		sm.playBuffer(sm.chomp1)
	} else {
		sm.playBuffer(sm.chomp2)
	}
	sm.chompFlip = !sm.chompFlip
}

// PlayPowerUp plays the power pellet sound.
func (sm *SoundManager) PlayPowerUp() {
	sm.playBuffer(sm.powerUp)
}

// PlayGhostEaten plays the ghost-eating chirp.
func (sm *SoundManager) PlayGhostEaten() {
	sm.playBuffer(sm.ghostEaten)
}

// PlayDeath plays the Pac-Man death sound.
func (sm *SoundManager) PlayDeath() {
	sm.playBuffer(sm.death)
}

// PlayLevelClear plays the level-clear arpeggio.
func (sm *SoundManager) PlayLevelClear() {
	sm.playBuffer(sm.levelClear)
}

func (sm *SoundManager) playBuffer(buf []byte) {
	player := sm.context.NewPlayerFromBytes(buf)
	player.Play()
}

// generateSquareWave creates a square wave at the given frequency and duration.
// Output is 16-bit signed little-endian stereo PCM.
func generateSquareWave(freq, duration float64, sr int) []byte {
	numSamples := int(duration * float64(sr))
	buf := new(bytes.Buffer)
	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(sr)
		var val int16
		if math.Sin(2*math.Pi*freq*t) >= 0 {
			val = 8000
		} else {
			val = -8000
		}
		binary.Write(buf, binary.LittleEndian, val) // left
		binary.Write(buf, binary.LittleEndian, val) // right
	}
	return buf.Bytes()
}

// generateSweep creates a frequency sweep from startFreq to endFreq.
func generateSweep(startFreq, endFreq, duration float64, sr int) []byte {
	numSamples := int(duration * float64(sr))
	buf := new(bytes.Buffer)
	phase := 0.0
	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(numSamples)
		freq := startFreq + (endFreq-startFreq)*t
		phase += 2 * math.Pi * freq / float64(sr)
		var val int16
		if math.Sin(phase) >= 0 {
			val = 8000
		} else {
			val = -8000
		}
		binary.Write(buf, binary.LittleEndian, val)
		binary.Write(buf, binary.LittleEndian, val)
	}
	return buf.Bytes()
}

// generateDeathSound creates the descending warble death sound.
func generateDeathSound(sr int) []byte {
	duration := 1.5
	numSamples := int(duration * float64(sr))
	buf := new(bytes.Buffer)
	phase := 0.0
	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(numSamples)
		baseFreq := 800 - 700*t // 800 -> 100
		warble := math.Sin(2*math.Pi*6*t) * 50
		freq := baseFreq + warble
		phase += 2 * math.Pi * freq / float64(sr)
		val := int16(math.Sin(phase) * 6000)
		binary.Write(buf, binary.LittleEndian, val)
		binary.Write(buf, binary.LittleEndian, val)
	}
	return buf.Bytes()
}

// generateArpeggio creates an ascending arpeggio for level clear.
func generateArpeggio(sr int) []byte {
	// C5, E5, G5, C6
	notes := []float64{523.25, 659.25, 783.99, 1046.50}
	noteDuration := 0.125
	buf := new(bytes.Buffer)
	for _, freq := range notes {
		numSamples := int(noteDuration * float64(sr))
		for i := 0; i < numSamples; i++ {
			t := float64(i) / float64(sr)
			val := int16(math.Sin(2*math.Pi*freq*t) * 6000)
			binary.Write(buf, binary.LittleEndian, val)
			binary.Write(buf, binary.LittleEndian, val)
		}
	}
	return buf.Bytes()
}
