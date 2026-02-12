package game

// Tile type constants.
const (
	TileWall = iota
	TileDot
	TilePowerPellet
	TileEmpty
	TileGhostHouse
	TileGhostDoor
)

// Spawn position constants.
const (
	PacmanSpawnX, PacmanSpawnY           = 14, 23
	GhostHouseCenterX, GhostHouseCenterY = 14, 14
)

// mazeLayout defines the classic Pac-Man maze as a 28x31 character grid.
// Characters: '#'=wall, '.'=dot, 'o'=power pellet, '-'=ghost door, 'G'=ghost house, ' '=empty
var mazeLayout = []string{
	"############################", // row 0
	"#............##............#", // row 1
	"#.####.#####.##.#####.####.#", // row 2
	"#o####.#####.##.#####.####o#", // row 3
	"#.####.#####.##.#####.####.#", // row 4
	"#..........................#", // row 5
	"#.####.##.########.##.####.#", // row 6
	"#.####.##.########.##.####.#", // row 7
	"#......##....##....##......#", // row 8
	"######.##### ## #####.######", // row 9
	"     #.##### ## #####.#     ", // row 10
	"     #.##          ##.#     ", // row 11
	"     #.## ###--### ##.#     ", // row 12
	"######.## #GGGGGG# ##.######", // row 13
	"      .   #GGGGGG#   .      ", // row 14
	"######.## #GGGGGG# ##.######", // row 15
	"     #.## ######## ##.#     ", // row 16
	"     #.##          ##.#     ", // row 17
	"     #.## ######## ##.#     ", // row 18
	"######.## ######## ##.######", // row 19
	"#............##............#", // row 20
	"#.####.#####.##.#####.####.#", // row 21
	"#.####.#####.##.#####.####.#", // row 22
	"#o..##.......  .......##..o#", // row 23
	"###.##.##.########.##.##.###", // row 24
	"###.##.##.########.##.##.###", // row 25
	"#......##....##....##......#", // row 26
	"#.##########.##.##########.#", // row 27
	"#.##########.##.##########.#", // row 28
	"#..........................#", // row 29
	"############################", // row 30
}

// Maze represents the game maze with tile tracking.
type Maze struct {
	Width         int
	Height        int
	tiles         [][]int
	remainingDots int
}

// NewMaze creates a new Maze by parsing the layout.
func NewMaze() *Maze {
	m := &Maze{
		Width:  MazeCols,
		Height: MazeRows,
	}
	m.parse()
	return m
}

// parse reads the mazeLayout and populates the tiles grid and dot count.
func (m *Maze) parse() {
	m.tiles = make([][]int, m.Height)
	m.remainingDots = 0
	for y := 0; y < m.Height; y++ {
		m.tiles[y] = make([]int, m.Width)
		row := mazeLayout[y]
		for x := 0; x < m.Width; x++ {
			var ch byte
			if x < len(row) {
				ch = row[x]
			} else {
				ch = ' '
			}
			switch ch {
			case '#':
				m.tiles[y][x] = TileWall
			case '.':
				m.tiles[y][x] = TileDot
				m.remainingDots++
			case 'o':
				m.tiles[y][x] = TilePowerPellet
				m.remainingDots++
			case '-':
				m.tiles[y][x] = TileGhostDoor
			case 'G':
				m.tiles[y][x] = TileGhostHouse
			case ' ':
				m.tiles[y][x] = TileEmpty
			default:
				m.tiles[y][x] = TileEmpty
			}
		}
	}
}

// TileAt returns the tile type at the given grid position.
// Returns TileWall for out-of-bounds coordinates.
func (m *Maze) TileAt(x, y int) int {
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return TileWall
	}
	return m.tiles[y][x]
}

// IsPassable returns true if Pac-Man can move through the tile at (x, y).
// Dots, power pellets, and empty tiles are passable. Walls, ghost house, and ghost doors are not.
func (m *Maze) IsPassable(x, y int) bool {
	t := m.TileAt(x, y)
	return t == TileDot || t == TilePowerPellet || t == TileEmpty
}

// IsPassableForGhost returns true if a ghost can move through the tile at (x, y).
// Ghosts can pass through everything except walls.
func (m *Maze) IsPassableForGhost(x, y int) bool {
	t := m.TileAt(x, y)
	return t != TileWall
}

// ConsumeDot consumes a dot or power pellet at (x, y).
// Returns true if a dot/pellet was consumed, false otherwise.
func (m *Maze) ConsumeDot(x, y int) bool {
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return false
	}
	t := m.tiles[y][x]
	if t == TileDot || t == TilePowerPellet {
		m.tiles[y][x] = TileEmpty
		m.remainingDots--
		return true
	}
	return false
}

// RemainingDots returns the number of unconsumed dots and power pellets.
func (m *Maze) RemainingDots() int {
	return m.remainingDots
}

// Reset re-parses the layout to restore all dots and power pellets.
func (m *Maze) Reset() {
	m.parse()
}
