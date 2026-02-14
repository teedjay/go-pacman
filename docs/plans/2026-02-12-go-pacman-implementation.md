# Go Pac-Man Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a graphical Pac-Man game in Go with Ebitengine v2, featuring pixel art sprites, programmatic sound, keyboard controls, and multi-level difficulty scaling.

**Architecture:** Standalone Go binary using Ebitengine v2 for rendering and audio. Game logic operates on a 28x31 tile grid (8x8 pixels per tile) with pixel-smooth movement. Game logic is separated from rendering for testability. Screen resolution is 224x288 (classic Pac-Man), scaled 3x to 672x864 window.

**Tech Stack:** Go 1.21+, Ebitengine v2 (`github.com/hajimehoshi/ebiten/v2`), standard library `math` for audio waveform generation.

**Reference:** See `docs/plans/2026-02-12-go-pacman-design.md` for full design rationale.

---

### Task 1: Project Scaffolding & Empty Window

**Files:**
- Create: `main.go`
- Create: `game/game.go`

**Step 1: Initialize Go module and add Ebitengine**

Run:
```bash
cd ./go-pacman
go mod init go-pacman
go get github.com/hajimehoshi/ebiten/v2@latest
```

**Step 2: Create `game/game.go`**

```go
package game

import "github.com/hajimehoshi/ebiten/v2"

const (
	TileSize     = 8
	MazeCols     = 28
	MazeRows     = 31
	HUDTopRows   = 3
	HUDBotRows   = 2
	ScreenWidth  = MazeCols * TileSize                        // 224
	ScreenHeight = (MazeRows + HUDTopRows + HUDBotRows) * TileSize // 288
	Scale        = 3
)

type Game struct{}

func New() *Game {
	return &Game{}
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
```

**Step 3: Create `main.go`**

```go
package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"go-pacman/game"
)

func main() {
	ebiten.SetWindowSize(game.ScreenWidth*game.Scale, game.ScreenHeight*game.Scale)
	ebiten.SetWindowTitle("Go Pac-Man")
	if err := ebiten.RunGame(game.New()); err != nil {
		log.Fatal(err)
	}
}
```

**Step 4: Build and run to verify**

Run: `go run main.go`
Expected: A black window opens titled "Go Pac-Man", 672x864 pixels. Close it manually.

**Step 5: Commit**

```bash
git add main.go game/game.go go.mod go.sum
git commit -m "feat: project scaffolding with empty Ebitengine window"
```

---

### Task 2: Maze Data Structure & Parsing

**Files:**
- Create: `game/maze.go`
- Create: `game/maze_test.go`

**Step 1: Write failing tests**

```go
// game/maze_test.go
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
```

**Step 2: Run tests to verify failure**

Run: `go test ./game/ -v`
Expected: FAIL — types not defined

**Step 3: Implement `game/maze.go`**

Define tile type constants, the maze layout as a `[]string`, and the `Maze` struct with methods:
- `NewMaze()` — parses the layout string into a 2D grid of tile types
- `TileAt(x, y)` — returns tile type at grid position
- `IsPassable(x, y)` — true for dot, power pellet, empty tiles; false for walls
- `ConsumeDot(x, y)` — if tile is dot or power pellet, change to empty, decrement counter, return true
- `RemainingDots()` — returns count of unconsumed dots + power pellets
- `Reset()` — re-parses the layout to restore all dots

Tile type constants:
```go
const (
	TileWall = iota
	TileDot
	TilePowerPellet
	TileEmpty
	TileGhostHouse
	TileGhostDoor
)
```

The maze layout should be a 28x31 character grid representing the classic Pac-Man maze. Use these characters:
- `#` = wall
- `.` = dot
- `o` = power pellet
- `-` = ghost house door
- `G` = ghost house interior (passable for ghosts only)
- ` ` (space) = empty/passable

The `Maze` struct stores:
- `Width, Height int`
- `tiles [][]int` — current state (mutable, dots get consumed)
- `initialTiles [][]int` — copy of initial state for reset

Ghost house tiles (`G`) should be treated as impassable for Pac-Man (via `IsPassable`) but passable for ghosts (add `IsPassableForGhost(x, y)` method).

Pac-Man spawn position: tile (14, 23) — center-bottom of maze.
Ghost spawn positions: tiles (12-15, 14) — inside ghost house.

Add exported constants for these spawn positions:
```go
const (
	PacmanSpawnX, PacmanSpawnY = 14, 23
	GhostHouseCenterX, GhostHouseCenterY = 14, 14
)
```

**Step 4: Run tests to verify pass**

Run: `go test ./game/ -v`
Expected: PASS

**Step 5: Commit**

```bash
git add game/maze.go game/maze_test.go
git commit -m "feat: maze data structure with parsing, dot tracking, and reset"
```

---

### Task 3: Maze Rendering

**Files:**
- Create: `game/sprite.go`
- Modify: `game/game.go` — add maze to Game struct, draw it

**Step 1: Create `game/sprite.go` with wall and dot rendering**

This file provides functions that return `*ebiten.Image` for each visual element. All sprites are generated programmatically at init time.

Implement:
- `GenerateWallTile()` — 8x8 blue tile. Draw a filled blue (#2121DE) square with 1px black border. This is the simple version; we'll refine wall shapes later.
- `GenerateDotSprite()` — 8x8 image with a 2x2 white square centered (at pixels 3-4, 3-4).
- `GeneratePowerPelletSprite()` — 8x8 image with a 6x6 white square centered (at pixels 1-6, 1-6).
- `GenerateEmptyTile()` — 8x8 fully black image.

Store generated sprites in package-level variables, initialized lazily or via an `InitSprites()` function called from `game.New()`.

**Step 2: Modify `game/game.go` to draw the maze**

Add `maze *Maze` field to Game struct. Initialize in `New()`. In `Draw()`, iterate over all tiles and draw the appropriate sprite at `(x*TileSize, y*TileSize + HUDTopRows*TileSize)`. The `HUDTopRows*TileSize` offset leaves room for the score display above the maze.

**Step 3: Run to verify**

Run: `go run main.go`
Expected: The Pac-Man maze appears — blue walls, white dots, larger white power pellets in the 4 corners. Black background.

**Step 4: Commit**

```bash
git add game/sprite.go game/game.go
git commit -m "feat: maze rendering with walls, dots, and power pellets"
```

---

### Task 4: Pac-Man Entity & Sprite

**Files:**
- Create: `game/pacman.go`
- Modify: `game/sprite.go` — add Pac-Man sprite generation
- Modify: `game/game.go` — add Pac-Man to Game, draw it

**Step 1: Create `game/pacman.go`**

```go
type Direction int

const (
	DirNone Direction = iota
	DirUp
	DirDown
	DirLeft
	DirRight
)

type PacMan struct {
	X, Y      float64   // pixel position (center of sprite)
	Dir       Direction // current movement direction
	NextDir   Direction // queued direction from input
	Speed     float64   // pixels per tick
	AnimFrame int       // 0, 1, 2 (closed, half, open)
	AnimTimer int       // ticks until next frame
	Alive     bool
}

func NewPacMan() *PacMan {
	return &PacMan{
		X:     float64(PacmanSpawnX*TileSize + TileSize/2),
		Y:     float64(PacmanSpawnY*TileSize + TileSize/2),
		Dir:   DirNone,
		Speed: 1.5,
		Alive: true,
	}
}

// TileX and TileY return current tile position
func (p *PacMan) TileX() int { return int(p.X) / TileSize }
func (p *PacMan) TileY() int { return int(p.Y) / TileSize }
```

**Step 2: Add Pac-Man sprites to `game/sprite.go`**

Generate 3 animation frames for Pac-Man as 13x13 pixel images, yellow (#FFFF00):
- Frame 0 (closed): filled circle
- Frame 1 (half open): circle with small triangular mouth (~30 degrees)
- Frame 2 (full open): circle with wide triangular mouth (~60 degrees)

The mouth faces right by default. When drawing, rotate the sprite based on `Dir`:
- Right: no rotation
- Left: flip horizontally
- Up: rotate 90° counter-clockwise
- Down: rotate 90° clockwise

Use `ebiten.GeoM` transforms for rotation when drawing.

**Step 3: Modify `game/game.go`**

Add `pacman *PacMan` field. Initialize in `New()`. In `Draw()`, draw Pac-Man's current animation frame at its pixel position, offset by the HUD area. Center the 13x13 sprite on the position.

**Step 4: Run to verify**

Run: `go run main.go`
Expected: Pac-Man (yellow circle with mouth) appears at spawn position in the maze.

**Step 5: Commit**

```bash
git add game/pacman.go game/sprite.go game/game.go
git commit -m "feat: Pac-Man entity with animated sprite"
```

---

### Task 5: Input Handling & Pac-Man Movement

**Files:**
- Create: `game/input.go`
- Create: `game/movement_test.go`
- Modify: `game/pacman.go` — add movement logic
- Modify: `game/game.go` — call input and movement in Update()

**Step 1: Write failing movement tests**

```go
// game/movement_test.go
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
	// Place at a tile center next to a wall, face into wall
	p.X = float64(1*TileSize + TileSize/2)
	p.Y = float64(0*TileSize + TileSize/2) // row 0 is all walls
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
	p.X = float64(1*TileSize + TileSize/2)
	p.Y = float64(1*TileSize + TileSize/2) // (1,1) is a dot
	p.Dir = DirRight
	p.NextDir = DirDown
	// At tile center, if down is passable, direction should switch
	p.Move(m)
	if m.IsPassable(1, 2) && p.Dir != DirDown {
		t.Error("should switch to queued direction when possible")
	}
}
```

**Step 2: Run tests to verify failure**

Run: `go test ./game/ -v -run TestPacMan`
Expected: FAIL

**Step 3: Implement movement in `game/pacman.go`**

Add methods:

`IsAtTileCenter() bool` — true if position is within `Speed` pixels of the nearest tile center.

`Move(m *Maze)` — the core movement logic:
1. If at tile center:
   a. Check if `NextDir` is passable → switch `Dir` to `NextDir`
   b. Check if current `Dir` is passable → continue
   c. Otherwise stop (set `Dir = DirNone`)
   d. Snap position to exact tile center
2. If `Dir != DirNone`, advance position by `Speed` pixels in `Dir`
3. Advance animation timer; cycle `AnimFrame` through 0→1→2→1→0...

Helper: `nextTile(x, y int, dir Direction) (int, int)` — returns the tile coordinates in the given direction from (x, y).

**Step 4: Create `game/input.go`**

```go
package game

import "github.com/hajimehoshi/ebiten/v2"

func ReadInput(p *PacMan) {
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		p.NextDir = DirUp
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		p.NextDir = DirDown
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		p.NextDir = DirLeft
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		p.NextDir = DirRight
	}
}
```

**Step 5: Wire up in `game/game.go`**

In `Update()`: call `ReadInput(g.pacman)` then `g.pacman.Move(g.maze)`.

**Step 6: Run tests**

Run: `go test ./game/ -v`
Expected: PASS

**Step 7: Run game to verify visually**

Run: `go run main.go`
Expected: Pac-Man moves with arrow keys/WASD. Stops at walls. Pre-turn buffering works.

**Step 8: Commit**

```bash
git add game/input.go game/movement_test.go game/pacman.go game/game.go
git commit -m "feat: keyboard input and Pac-Man movement with wall collision"
```

---

### Task 6: Dot Consumption & Scoring

**Files:**
- Create: `game/score_test.go`
- Modify: `game/game.go` — add score tracking, dot consumption in Update()

**Step 1: Write failing tests**

```go
// game/score_test.go
package game

import "testing"

func TestDotScoring(t *testing.T) {
	g := New()
	g.maze.ConsumeDot(1, 1) // simulate eating a dot
	// We need a score field and scoring logic
	// Test that consuming a dot at pacman's position awards 10 points
	g.pacman.X = float64(1*TileSize + TileSize/2)
	g.pacman.Y = float64(1*TileSize + TileSize/2)
	g.maze.Reset() // restore dots
	g.checkDotConsumption()
	if g.score != 10 {
		t.Errorf("expected score 10, got %d", g.score)
	}
}

func TestPowerPelletScoring(t *testing.T) {
	g := New()
	// Move pacman to a power pellet position
	g.pacman.X = float64(1*TileSize + TileSize/2)
	g.pacman.Y = float64(3*TileSize + TileSize/2)
	g.checkDotConsumption()
	if g.score != 50 {
		t.Errorf("expected score 50 for power pellet, got %d", g.score)
	}
}
```

**Step 2: Run tests to verify failure**

Run: `go test ./game/ -v -run TestDot`
Expected: FAIL — `score` field and `checkDotConsumption` not defined

**Step 3: Implement scoring in `game/game.go`**

Add fields to `Game`:
```go
score     int
highScore int
lives     int
level     int
```

Initialize `lives = 3`, `level = 1` in `New()`.

Add method `checkDotConsumption()`:
1. Get Pac-Man's current tile position
2. Check tile type at that position
3. If dot: `score += 10`, consume it
4. If power pellet: `score += 50`, consume it, trigger frightened mode (placeholder for now)
5. Check for extra life at 10,000 points

Call `checkDotConsumption()` in `Update()` after movement.

**Step 4: Run tests**

Run: `go test ./game/ -v`
Expected: PASS

**Step 5: Run game to verify**

Run: `go run main.go`
Expected: Dots disappear as Pac-Man moves over them.

**Step 6: Commit**

```bash
git add game/score_test.go game/game.go
git commit -m "feat: dot consumption with scoring (10 for dots, 50 for power pellets)"
```

---

### Task 7: Ghost Entities & Sprites

**Files:**
- Create: `game/ghost.go`
- Modify: `game/sprite.go` — add ghost sprite generation
- Modify: `game/game.go` — add ghosts to Game, draw them

**Step 1: Create `game/ghost.go`**

```go
type GhostMode int

const (
	GhostChase GhostMode = iota
	GhostScatter
	GhostFrightened
	GhostEaten
)

type GhostID int

const (
	Blinky GhostID = iota // red
	Pinky                 // pink
	Inky                  // cyan
	Clyde                 // orange
)

type Ghost struct {
	ID        GhostID
	X, Y      float64
	Dir       Direction
	Mode      GhostMode
	Speed     float64
	SpawnX    int // tile coords for scatter target
	SpawnY    int
	ScatterX  int // scatter corner target
	ScatterY  int
	InHouse   bool // still in ghost house
	ExitTimer int  // ticks until ghost leaves house
}
```

Define the 4 ghosts with their colors and scatter corners:
- Blinky (red): scatter to top-right (25, 0)
- Pinky (pink): scatter to top-left (2, 0)
- Inky (cyan): scatter to bottom-right (27, 30)
- Clyde (orange): scatter to bottom-left (0, 30)

`NewGhosts()` returns `[4]*Ghost` initialized at ghost house positions:
- Blinky starts outside house at (14, 11), others inside at (12, 14), (14, 14), (16, 14)
- Pinky exits after 0s, Inky after 5s, Clyde after 10s (stagger exits via `ExitTimer`)

**Step 2: Add ghost sprites to `game/sprite.go`**

Generate ghost sprite as a 13x13 image:
- Body: rounded top (semicircle), flat sides, wavy bottom edge (3 bumps)
- Eyes: two white circles with blue pupils pointing in movement direction
- Color parameter: red (#FF0000), pink (#FFB8FF), cyan (#00FFFF), orange (#FFB852)

Also generate:
- Frightened sprite: all blue (#2121FF) body, white squiggly mouth, no directional eyes
- Eaten sprite: just the eyes (white circles with blue pupils), no body

**Step 3: Modify `game/game.go`**

Add `ghosts [4]*Ghost` field. Initialize in `New()`. In `Draw()`, draw each ghost at its pixel position with the appropriate sprite based on `Mode`.

**Step 4: Run to verify**

Run: `go run main.go`
Expected: Four colored ghosts appear in/near the ghost house. They don't move yet.

**Step 5: Commit**

```bash
git add game/ghost.go game/sprite.go game/game.go
git commit -m "feat: ghost entities with colored pixel art sprites"
```

---

### Task 8: Ghost AI — State Machine & Pathfinding

**Files:**
- Create: `game/ghost_ai.go`
- Create: `game/ghost_ai_test.go`
- Modify: `game/game.go` — call ghost AI in Update()

**Step 1: Write failing tests**

```go
// game/ghost_ai_test.go
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
	// Path into a wall should return empty
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
```

**Step 2: Run tests to verify failure**

Run: `go test ./game/ -v -run TestBFS`
Expected: FAIL

**Step 3: Implement `game/ghost_ai.go`**

**BFS pathfinding:**
```go
func BFS(m *Maze, startX, startY, targetX, targetY int) []Direction
```
Standard BFS on the tile grid. Returns a slice of directions from start to target. Uses `m.IsPassable()` (or `IsPassableForGhost()` for ghost-house tiles). Returns empty slice if no path found.

**Mode timer:**
```go
type ModeTimer struct {
	phases []struct {
		mode   GhostMode
		ticks  int // -1 for infinite
	}
	currentPhase int
	ticksInPhase int
}

func NewModeTimer(level int) *ModeTimer
func (mt *ModeTimer) Tick()
func (mt *ModeTimer) CurrentMode() GhostMode
func (mt *ModeTimer) Reset()
```

Level 1 phases (at 60 TPS):
- Scatter 7s (420 ticks)
- Chase 20s (1200 ticks)
- Scatter 7s
- Chase 20s
- Scatter 5s (300 ticks)
- Chase forever (-1)

**Ghost movement — `ChooseDirection`:**

At each tile center, the ghost picks the direction that brings it closest to its target tile. Rules:
1. Never reverse direction (except when mode changes — handled separately)
2. For each possible direction (up, left, down, right — in that priority order), check if the next tile is passable
3. Pick the direction where the next tile has the smallest Euclidean distance to the target
4. In frightened mode, pick a random valid direction instead

**Ghost Update loop — `UpdateGhost(g *Ghost, m *Maze, pacman *PacMan, mode GhostMode)`:**
1. If `InHouse`: handle ghost house exit logic (decrement `ExitTimer`, move up through door when ready)
2. Set target based on mode:
   - Chase: Pac-Man's current tile (plus small random offset ±2 tiles to prevent clumping)
   - Scatter: ghost's assigned corner
   - Frightened: random direction at intersections
   - Eaten: ghost house entrance (14, 11)
3. If at tile center, call `ChooseDirection` to pick next direction
4. Move by `Speed` pixels in current direction

**Step 4: Wire up in `game/game.go`**

Add `modeTimer *ModeTimer` to Game. In `Update()`:
1. Tick the mode timer
2. For each ghost, call `UpdateGhost()` with current mode
3. Handle ghost house staggered exits

**Step 5: Run tests**

Run: `go test ./game/ -v`
Expected: PASS

**Step 6: Run game to verify**

Run: `go run main.go`
Expected: Ghosts leave the ghost house one by one and chase/scatter around the maze.

**Step 7: Commit**

```bash
git add game/ghost_ai.go game/ghost_ai_test.go game/game.go
git commit -m "feat: ghost AI with BFS pathfinding, chase/scatter mode timer"
```

---

### Task 9: Collision, Frightened Mode & Ghost Eating

**Files:**
- Create: `game/collision_test.go`
- Modify: `game/game.go` — collision detection, frightened mode trigger, death logic

**Step 1: Write failing tests**

```go
// game/collision_test.go
package game

import (
	"testing"
)

func TestCollisionDetection(t *testing.T) {
	p := NewPacMan()
	g := &Ghost{X: p.X, Y: p.Y} // same position
	if !CheckCollision(p, g) {
		t.Error("overlapping positions should collide")
	}
	g.X = p.X + 20 // far away
	if CheckCollision(p, g) {
		t.Error("distant positions should not collide")
	}
}

func TestCollisionThreshold(t *testing.T) {
	p := NewPacMan()
	g := &Ghost{X: p.X + 5, Y: p.Y} // close but not overlapping
	if !CheckCollision(p, g) {
		t.Error("positions within threshold (6px) should collide")
	}
	g.X = p.X + 7 // just outside threshold
	if CheckCollision(p, g) {
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
}
```

**Step 2: Run tests to verify failure**

Run: `go test ./game/ -v -run TestCollision`
Expected: FAIL

**Step 3: Implement collision system**

`CheckCollision(p *PacMan, g *Ghost) bool` — returns true if distance between centers is less than 6 pixels.

In `game.go`, add method `checkGhostCollisions()` called each tick:
1. For each ghost, check collision with Pac-Man
2. If ghost is in **chase or scatter** mode: Pac-Man dies
   - Set `pacman.Alive = false`
   - Decrement `lives`
   - Start death animation timer (placeholder — just respawn after 120 ticks)
   - Reset ghost positions
3. If ghost is in **frightened** mode: ghost is eaten
   - Award escalating points: 200 × 2^(ghostsEatenCombo)
   - Set ghost mode to `GhostEaten`
   - Increment `ghostsEatenCombo`
4. If ghost is **eaten**: no collision

Add frightened mode trigger in `checkDotConsumption()`:
- When power pellet consumed, set all non-eaten ghosts to `GhostFrightened`
- Reverse their direction
- Reset `ghostsEatenCombo` to 0
- Start frightened timer (based on level — see Task 13)

Add `frightenedTimer int` to Game. Decrement each tick. When it reaches 0, switch all frightened ghosts back to the current mode timer mode.

**Step 4: Run tests**

Run: `go test ./game/ -v`
Expected: PASS

**Step 5: Run game to verify**

Run: `go run main.go`
Expected: Eating a power pellet turns ghosts blue. Pac-Man can eat blue ghosts for points. Running into a normal ghost kills Pac-Man and respawns.

**Step 6: Commit**

```bash
git add game/collision_test.go game/game.go
git commit -m "feat: ghost collision, frightened mode, and ghost eating"
```

---

### Task 10: Sound Effects

**Files:**
- Create: `game/sound.go`
- Modify: `game/game.go` — initialize audio, play sounds on events

**Step 1: Implement `game/sound.go`**

Use Ebitengine's `audio` package. Initialize an `audio.Context` at 44100 Hz sample rate.

Generate PCM waveforms as `[]byte` buffers (16-bit signed little-endian stereo):

```go
type SoundManager struct {
	context    *audio.Context
	chomp1     []byte // ~260Hz square wave, 60ms
	chomp2     []byte // ~390Hz square wave, 60ms
	chompFlip  bool   // alternates between chomp1 and chomp2
	powerUp    []byte // descending sweep 800→200Hz, 300ms
	ghostEaten []byte // ascending chirp 200→1200Hz, 150ms
	death      []byte // descending warble 800→100Hz, 1.5s
	levelClear []byte // ascending arpeggio, 500ms
}

func NewSoundManager() *SoundManager
func (sm *SoundManager) PlayChomp()
func (sm *SoundManager) PlayPowerUp()
func (sm *SoundManager) PlayGhostEaten()
func (sm *SoundManager) PlayDeath()
func (sm *SoundManager) PlayLevelClear()
```

Waveform generation helper:
```go
func generateSquareWave(freq, duration float64, sampleRate int) []byte
func generateSweep(startFreq, endFreq, duration float64, sampleRate int) []byte
```

For the death sound, add a frequency modulation warble: `freq = baseFreq + math.Sin(t*warbleRate)*warbleDepth`

Use `audio.NewPlayerFromBytes()` to create one-shot players. Keep a pool of 4 players to allow overlapping sounds.

**Step 2: Wire up in `game/game.go`**

Add `sound *SoundManager` to Game. Initialize in `New()`.

Play sounds at the right moments:
- `PlayChomp()` — when dot consumed (alternates between chomp1/chomp2)
- `PlayPowerUp()` — when power pellet consumed
- `PlayGhostEaten()` — when ghost eaten
- `PlayDeath()` — when Pac-Man dies
- `PlayLevelClear()` — when all dots consumed

**Step 3: Run to verify**

Run: `go run main.go`
Expected: Hear "waka waka" when eating dots, descending sweep for power pellets, chirp for eating ghosts.

**Step 4: Commit**

```bash
git add game/sound.go game/game.go
git commit -m "feat: programmatic sound effects for all game events"
```

---

### Task 11: HUD — Score, Lives, Level Display

**Files:**
- Create: `game/hud.go`
- Modify: `game/sprite.go` — add pixel font glyphs
- Modify: `game/game.go` — call HUD drawing

**Step 1: Add pixel font to `game/sprite.go`**

Define a minimal pixel font as a map of rune to 5x7 pixel bitmap (standard small pixel font size). Only need: digits 0-9, letters A-Z (for "SCORE", "HIGH", "LEVEL", "READY", "GAME", "OVER"), space, and exclamation mark.

Each glyph is stored as a `[7]uint8` where each uint8's lower 5 bits represent one row of pixels.

```go
var fontGlyphs = map[rune][7]uint8{
	'0': {0x0E, 0x11, 0x13, 0x15, 0x19, 0x11, 0x0E},
	'1': {0x04, 0x0C, 0x04, 0x04, 0x04, 0x04, 0x0E},
	// ... etc
}
```

Function `DrawText(screen *ebiten.Image, text string, x, y int, color color.Color)` renders text using the font.

**Step 2: Create `game/hud.go`**

```go
func DrawHUD(screen *ebiten.Image, score, highScore, lives, level int)
```

Layout:
- Row 0 (y=0): "1UP" left-aligned, "HIGH SCORE" centered
- Row 1 (y=8): score left-aligned, high score centered
- Below maze (y = (HUDTopRows+MazeRows)*TileSize):
  - Left: small Pac-Man sprites × lives remaining
  - Right: "LEVEL" + level number

All text in white. Score shown without leading zeros except always at least 1 digit.

**Step 3: Wire up in `game/game.go`**

Call `DrawHUD()` in `Draw()` after drawing the maze and entities.

**Step 4: Run to verify**

Run: `go run main.go`
Expected: Score display at top updates as dots are eaten. Lives shown as small Pac-Man icons at bottom. Level number displayed.

**Step 5: Commit**

```bash
git add game/hud.go game/sprite.go game/game.go
git commit -m "feat: HUD with score, high score, lives, and level display"
```

---

### Task 12: Game State Machine

**Files:**
- Create: `game/state.go`
- Create: `game/state_test.go`
- Modify: `game/game.go` — integrate state machine into Update/Draw

**Step 1: Write failing tests**

```go
// game/state_test.go
package game

import "testing"

type GameState int

const (
	StateTitle GameState = iota
	StateReady
	StatePlaying
	StateDeath
	StateLevelClear
	StateGameOver
)

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
			// Verify each transition is valid
			if !isValidTransition(tt.from, tt.to) {
				t.Errorf("transition %d -> %d should be valid", tt.from, tt.to)
			}
		})
	}
}
```

**Step 2: Run tests to verify failure**

Run: `go test ./game/ -v -run TestState`
Expected: FAIL

**Step 3: Implement `game/state.go`**

Move `GameState` constants out of test file into `state.go`.

Add `state GameState` and `stateTimer int` to Game struct.

Each state controls what happens in `Update()` and `Draw()`:

- **StateTitle**: Draw "GO PAC-MAN" and "PRESS SPACE TO START". Wait for space key.
- **StateReady**: Show maze with "READY!" text overlay. Timer counts down 120 ticks (2 seconds), then transition to Playing.
- **StatePlaying**: Normal gameplay — input, movement, collision, ghost AI.
- **StateDeath**: Play death animation (Pac-Man shrinks/disappears over 90 ticks). Then either respawn (→ StateReady) or game over if lives == 0.
- **StateLevelClear**: Maze flashes (walls blink white/blue) for 120 ticks. Then increment level, reset maze, go to StateReady.
- **StateGameOver**: Show "GAME OVER" text for 180 ticks. Then return to title.

**Step 4: Run tests**

Run: `go test ./game/ -v`
Expected: PASS

**Step 5: Run to verify**

Run: `go run main.go`
Expected: Game starts at title screen. Space starts the game with "READY!" countdown. Death and game over flow correctly.

**Step 6: Commit**

```bash
git add game/state.go game/state_test.go game/game.go
git commit -m "feat: game state machine with title, ready, death, and game over states"
```

---

### Task 13: Level Progression & Difficulty Scaling

**Files:**
- Create: `game/difficulty.go`
- Create: `game/difficulty_test.go`
- Modify: `game/game.go` — use difficulty params

**Step 1: Write failing tests**

```go
// game/difficulty_test.go
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
	// Speed should increase monotonically
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
```

**Step 2: Run tests to verify failure**

Run: `go test ./game/ -v -run TestDifficulty`
Expected: FAIL

**Step 3: Implement `game/difficulty.go`**

```go
type DifficultyParams struct {
	PacManSpeed      float64
	GhostSpeed       float64
	FrightenedSpeed  float64
	FrightenedTicks  int // 60 TPS
	ScatterTicks     []int // scatter phase durations
	EatenSpeed       float64
}

func GetDifficulty(level int) DifficultyParams
```

Interpolate linearly between anchor points:

| Param | Level 1 | Level 5 | Level 10+ |
|---|---|---|---|
| PacManSpeed | 1.5 | 1.7 | 1.8 |
| GhostSpeed | 1.3 | 1.6 | 1.8 |
| FrightenedSpeed | 0.8 | 0.8 | 0.8 |
| FrightenedTicks | 360 (6s) | 180 (3s) | 60 (1s) |
| ScatterTicks[0] | 420 (7s) | 300 (5s) | 180 (3s) |
| EatenSpeed | 3.0 | 3.0 | 3.0 |

Clamp level to 10 for interpolation. Use `lerp(a, b, t)` where `t = (level-1) / 9.0` clamped to [0, 1].

**Step 4: Wire up in `game/game.go`**

When level changes or game starts, call `GetDifficulty(level)` and apply params:
- Set `pacman.Speed`
- Set each ghost's `Speed`
- Set `frightenedTimer` duration
- Pass scatter durations to `ModeTimer`

On level clear:
1. Increment `level`
2. Reset maze (restores all dots)
3. Reset Pac-Man to spawn
4. Reset ghosts to ghost house
5. Apply new difficulty params
6. Reset mode timer with new level

**Step 5: Run tests**

Run: `go test ./game/ -v`
Expected: PASS

**Step 6: Run game to verify**

Run: `go run main.go`
Expected: After clearing all dots, level increments, maze resets, game gets faster.

**Step 7: Commit**

```bash
git add game/difficulty.go game/difficulty_test.go game/game.go
git commit -m "feat: level progression with interpolated difficulty scaling"
```

---

### Task 14: Polish — Tunnel, Death Animation, Blinking Pellets, Extra Life

**Files:**
- Modify: `game/pacman.go` — tunnel wrapping, death animation
- Modify: `game/maze.go` — blinking power pellets
- Modify: `game/game.go` — extra life logic, tunnel support

**Step 1: Tunnel wrapping**

In `PacMan.Move()` and `UpdateGhost()`, after moving, check if position exits the maze bounds horizontally:
- If `X < 0`, set `X = MazeCols*TileSize` (wrap to right)
- If `X >= MazeCols*TileSize`, set `X = 0` (wrap to left)

The tunnel row in the maze layout should have passable tiles at both edges.

**Step 2: Death animation**

In `PacMan`, add `DeathFrame int` and `DeathTimer int`. When Pac-Man dies:
1. Freeze all movement for 60 ticks
2. Play death animation: Pac-Man's mouth opens wider and wider (like a pie chart going from full to empty) over 90 ticks — roughly 11 frames
3. Generate death animation frames in `sprite.go`: progressively larger mouth angle from 60° to 360°

In `Draw()`, when in `StateDeath`, draw the current death animation frame instead of the normal Pac-Man sprite.

**Step 3: Blinking power pellets**

In `Draw()`, toggle power pellet visibility every 15 ticks (quarter-second blink). Use a frame counter: `if (tickCount / 15) % 2 == 0 { draw pellet }`.

**Step 4: Extra life**

In `checkDotConsumption()` (or wherever score increases), check if score crossed 10,000 for the first time:
```go
if g.score >= 10000 && !g.extraLifeAwarded {
	g.lives++
	g.extraLifeAwarded = true
	// Optionally play a sound
}
```

**Step 5: Run to verify everything**

Run: `go run main.go`
Expected: Power pellets blink. Pac-Man wraps through tunnels. Death plays a shrinking animation. Extra life awarded at 10K.

**Step 6: Commit**

```bash
git add game/pacman.go game/maze.go game/sprite.go game/game.go
git commit -m "feat: tunnel wrapping, death animation, blinking pellets, extra life"
```

---

## Final Verification

After all 14 tasks, run the full test suite and play-test:

```bash
go test ./game/ -v -count=1
go run main.go
```

Verify:
- [ ] Game launches to title screen
- [ ] Space starts game with "READY!" countdown
- [ ] Pac-Man moves with arrow keys and WASD
- [ ] Pre-turn buffering feels responsive
- [ ] Dots and power pellets consumed with correct scores
- [ ] "Waka waka" sound plays when eating dots
- [ ] Ghosts leave house staggered, chase and scatter
- [ ] Power pellet turns ghosts blue, Pac-Man can eat them
- [ ] Eaten ghosts return to house as eyes
- [ ] Ghost collision kills Pac-Man with death animation
- [ ] Lives decrement, game over at 0 lives
- [ ] Clearing all dots advances level
- [ ] Difficulty increases with level (faster ghosts, shorter fright)
- [ ] Power pellets blink
- [ ] Tunnel wrapping works
- [ ] Extra life at 10,000 points
- [ ] High score tracks across games in session
