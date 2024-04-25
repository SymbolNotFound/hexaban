// Hexaban (some know it as Hexoban to match the 'o' in Sokoban).
//
// Structures for representing the initial state and sequence of moves
// for the solitaire game that is a hexagonal variant of the classic
// puzzle known as Sokoban ("warehouse worker," the crate pusher).

package hexaban

import (
	"encoding/json"
)

type Collection struct {
	Puzzles []Puzzle `json:"puzzles"`
	Source  string   `json:"source"`
	Author  string   `json:"author,omitempty"`
}

type Puzzle struct {
	Identity string     `json:"id"`
	Author   string     `json:"author"`
	Name     string     `json:"name"`
	Terrain  []HexCoord `json:"terrain"`
	Init     Init       `json:"init"`
}

type Init struct {
	Walls  []HexCoord `json:"walls,omitempty"`
	Goals  []HexCoord `json:"goals"`
	Crates []HexCoord `json:"crates"`
	Player HexCoord   `json:"ichiban"`
}

type Tile interface {
	I() int
	J() int
	Type() TileType
}

// This assumes that the rectangular coordinates collapse each pair of rows
// into a common row coordinate, aligning the columns, and that hexes have
// their orientation such that up/down movement is possible (the other four
// directions are left/right and north/south).
//
// .           up           |
// .    left _---_ north    | 0,0     2,0
// .        <_   _>         |     1,0     3,0
// .   south  ```  right    | 0,1     2,1
// .          down          |     1,1     3,1
//
// The selection of cardinal directions is arbitrary, this is good as any.  But
// notice that sometimes a column difference of |1| is a movement `right`, and
// sometimes it is a movement `north`.  Likewise, rows are always down, but
// their position when projected onto a canvas is higher or lower based on
// their column (despite the higher/lower dimension being a row-order feature!).
//
// For this reason, the rectilinear format is abandoned in favor of a coordinate
// system based on `right` and `south` directions, as `HexCoord`.  The RectCoord
// is still needed for parsing and converting maps defined elsewhere using this
// format.  There are many instances of this because the rectangular coordinates
// are easier to map to memory and can be printed to line-based consoles easily.
type RectCoord struct {
	col uint
	row uint
}

// Converts a rectangular coordinate into the equivalent hexagonal coordinate.
// Assumes that the (0, 0) center for both systems exists at top-left position.
func (coord RectCoord) ToHex() HexCoord {
	cHalf := coord.col >> 1
	cOdd := coord.col & 1

	// Each column is half i, half -j (with odd columns adding one more i).  Each
	// row is one i and one j.  We can reconstruct (i, j) from rows and columns
	// by accumulating these partial representations per-column and per-row.
	return HexCoord{
		i: int(cHalf + cOdd + coord.row),
		j: int(coord.row - cHalf),
	}
}

// Similar to the conversion offered by ToHex(), while translating the entire
// board to be centered at the indicated (hex) coordintae.
func (coord RectCoord) ToHexCenteredAt(center HexCoord) HexCoord {
	hex_coord := coord.ToHex()
	hex_coord.i -= center.i
	hex_coord.j -= center.j
	return hex_coord
}

// The HexCoord has a basis of two unit vectors (i, j) corresponding to
// right (downward-right in pixel coordinates) and south (or downward-left
// in pixel coords).  This results in a back-to-front ordering if sorted by
// (i + j), with ties broken by ascending `i`.  It is also a more suitable
// representation for translations, rotations and other affine transformations.
//
// The representation is serialized into JSON as a two-element array, [i, j].
type HexCoord struct {
	i int // `json:[0]`
	j int // `json:[1]`
}

// Accessors are read-only, the coordinate can only be changed by this package.
func (coord HexCoord) I() int { return coord.i }
func (coord HexCoord) J() int { return coord.j }

// Serialize the HexCoord value as a JSON formatted byte slice.
// NOTE: Does not need to be used directly, will be inferred by a parent
// struct's property being of type HexCoord, when the parent is serialized
// via json.Marshal(...).
func (coord HexCoord) MarshalJSON() ([]byte, error) {
	asArray := []int{coord.i, coord.j}
	return json.Marshal(asArray)
}

// Parse JSON into (already allocated) HexCoord struct.  Do not call this
// directly, prefer to call json.Unmarshal(...) directly.
func (coord *HexCoord) UnmarshalJSON(encoded []byte) error {
	array2D := make([]int, 2)
	if err := json.Unmarshal(encoded, &array2D); err != nil {
		return err
	}

	coord.i = array2D[0]
	coord.j = array2D[1]
	return nil
}

type TileType string

const (
	TILE_FLOOR  TileType = "FLOOR"
	TILE_WALL   TileType = "WALL"
	TILE_GOAL   TileType = "GOAL"
	TILE_CRATE  TileType = "CRATE"
	TILE_PLAYER TileType = "PLAYER"
)

type Wall struct {
	HexCoord
}

func (tile Wall) Type() TileType { return TILE_WALL }

type Floor struct {
	HexCoord
}

func (tile Floor) Type() TileType { return TILE_FLOOR }

type Crate struct {
	HexCoord
}

func (tile Crate) Type() TileType { return TILE_CRATE }

type Goal struct {
	HexCoord
}

func (tile Goal) Type() TileType { return TILE_GOAL }

type Player struct {
	HexCoord
}

func (tile Player) Type() TileType { return TILE_PLAYER }
