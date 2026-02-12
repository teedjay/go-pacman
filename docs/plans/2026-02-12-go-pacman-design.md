# Go Pac-Man Design

A graphical Pac-Man game in Go using Ebitengine v2. Simplified classic gameplay with pixel art retro visuals, programmatic sprites and sound, keyboard controls, and multi-level difficulty scaling.

## Architecture

Standalone Go binary built with Ebitengine v2. Game logic is separated from rendering.

```
go-pacman/
  main.go              # Entry point, window setup
  game/
    game.go            # Top-level game state machine (menu, playing, game-over)
    maze.go            # Maze definition, tile types, dot tracking
    pacman.go          # Pac-Man entity: position, direction, animation
    ghost.go           # Ghost entity: position, AI state machine
    sprite.go          # Programmatic pixel sprite generation
    sound.go           # Programmatic sound effect generation
    input.go           # Keyboard input handling
    hud.go             # Score, lives, level display
```

Core loop uses Ebitengine's `Update()` / `Draw()` cycle at 60 TPS. Game logic runs on a tile-based grid (28x31 tiles). Movement is pixel-smooth but snaps decisions to tile boundaries.

## Maze & Tiles

The maze is a 2D array of tile types: wall, dot, power pellet, empty, and ghost house. The classic layout is hardcoded as a string grid parsed at startup:

```go
var mazeData = []string{
    "############################",
    "#............##............#",
    "#.####.#####.##.#####.####.#",
    "#o####.#####.##.#####.####o#",
    // ...
}
```

- `#` = wall, `.` = dot, `o` = power pellet, `-` = ghost house door, space = empty
- Walls drawn as blue pixel-art tiles with rounded inner edges
- Dots: 2x2 pixel squares centered in tile
- Power pellets: 6x6 pixels, blink on a timer
- Tunnel wrapping at maze edges teleports to opposite side
- Level clears when all dots and power pellets are consumed

## Pac-Man Movement & Input

Pac-Man has a current direction and a queued direction. Arrow keys and WASD are supported. At each tile center:

1. Can Pac-Man turn in the queued direction? If yes, switch to it.
2. Otherwise, continue in the current direction if possible.
3. If neither works, Pac-Man stops.

This pre-turn buffering makes movement feel responsive.

Movement speed is in pixels per tick, starting at ~1.5 px/tick at level 1, increasing with level.

Animation cycles through 3 frames (mouth closed, half open, fully open), rotated based on current direction. Frames advance only while moving.

Collision detection:
- Dots: tile-based, consumed when Pac-Man's center enters the tile
- Ghosts: pixel-based with overlap threshold for fairer feel

## Ghost AI

Four ghosts share a three-mode state machine: chase, scatter, frightened.

**State transitions** follow a global timer alternating scatter and chase. Level 1 pattern: scatter 7s, chase 20s, scatter 7s, chase 20s, scatter 5s, chase forever. Higher levels shorten scatter and lengthen chase.

**Chase mode:** Target Pac-Man's current tile using BFS pathfinding. Each ghost gets a small random offset to its target to prevent clumping and create emergent flanking.

**Scatter mode:** Each ghost retreats to an assigned corner of the maze, giving the player breathing room.

**Frightened mode:** Triggered by power pellet. Ghosts reverse direction, turn blue, move randomly at intersections. Pac-Man can eat them for escalating points (200, 400, 800, 1600). Duration decreases with level (6s at level 1, down to 1s at high levels).

**Eaten ghosts:** Return to ghost house as floating eyes at double speed via direct pathfinding. Respawn and resume normal behavior.

## Sound Design

All sounds generated programmatically using Ebitengine's audio API. Short PCM waveforms generated in memory at startup.

- **Dot chomp:** Square wave alternating ~260Hz/~390Hz, ~60ms. Creates "waka waka" effect.
- **Power pellet:** Descending square wave sweep ~800Hz to ~200Hz, ~300ms.
- **Ghost eaten:** Ascending chirp ~200Hz to ~1200Hz, ~150ms.
- **Death:** Descending sine wave spiral ~800Hz to ~100Hz, ~1.5s with frequency modulation warble.
- **Level clear:** Ascending arpeggio, 3-4 sine tones stepping up a major scale, ~500ms.

Sounds are pre-generated `[]byte` buffers played via a small pool of audio players for overlap support.

## HUD, Scoring & Difficulty

**HUD** drawn outside the maze area:
- Top: score and high score
- Bottom: remaining lives (small Pac-Man sprites) and level number
- Text rendered using pixel font glyphs defined in code (digits 0-9, "SCORE", "HIGH SCORE", "LEVEL", "READY!", "GAME OVER")

**Scoring:**
- Dot: 10 points
- Power pellet: 50 points
- Ghosts: 200, 400, 800, 1600 (resets each power pellet)
- Extra life at 10,000 points (one time)

**Lives:** Start with 3. On death, play death animation, respawn Pac-Man, reset ghosts to ghost house. Dots remain. Game over when all lives lost.

**Difficulty scaling per level:**

| Parameter | Level 1 | Level 5 | Level 10+ |
|---|---|---|---|
| Pac-Man speed | 1.5 px/tick | 1.7 | 1.8 |
| Ghost speed | 1.3 | 1.6 | 1.8 |
| Frightened duration | 6s | 3s | 1s |
| Scatter phase length | 7s | 5s | 3s |

Values interpolated between anchors for smooth difficulty curve.

**Game state machine:** Title screen -> "READY!" countdown -> playing -> (death -> respawn | level clear -> next level) -> game over -> title screen.
