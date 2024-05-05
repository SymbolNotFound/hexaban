// Copyright (c) 2024 Symbol Not Found L.L.C.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// github:SymbolNotFound/hexaban/hexaban.go

package hexaban

// Hexaban (some know it as Hexoban to match the 'o' in Sokoban).
//
// Structures for representing the initial state and sequence of moves
// for the solitaire game that is a hexagonal variant of the classic
// puzzle known as Sokoban ("warehouse worker," the crate pusher).

import (
	"encoding/json"
)

type Collection struct {
	Puzzles []Puzzle `json:"puzzles"`
	Source  string   `json:"source"`
	Author  string   `json:"author,omitempty"`
}

type Puzzle struct {
	Identity   string     `json:"id"`
	Name       string     `json:"name"`
	Author     string     `json:"author"`
	Source     string     `json:"source"`
	Terrain    []HexCoord `json:"terrain"`
	Init       Init       `json:"init"`
	Difficulty int        `json:"difficulty,omitempty"`
}

type Init struct {
	// All goals and crates need to be on coordinates where terrain exists.
	Goals  []HexCoord `json:"goals"`
	Crates []HexCoord `json:"crates"`
	Player HexCoord   `json:"ichiban,omitempty"` // Default is (0, 0)
}

type Tile interface {
	I() int
	J() int
	Coord() HexCoord
	Type() TileType
}

// This assumes that the rectangular coordinates collapse each pair of rows
// into a common row coordinate, aligning the columns, and that hexes have
// their orientation such that up/down movement is possible (the other four
// directions are left/right and north/south).
//
// .           up           /===================== rect{col, row}
// .    left _---_ north    | 0,0     2,0     4,0
// .        <_   _>         |     1,0     3,0
// .   south  ```  right    | 0,1     2,1     4,1
// .          down          |     1,1     3,1
// .                        | 0,2     2,2     4,2
// rect{6,1} = hex{4,2}     \     ...     ...
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

func NewRectCoord(col, row uint) RectCoord {
	return RectCoord{col, row}
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
func (coord HexCoord) I() int          { return coord.i }
func (coord HexCoord) J() int          { return coord.j }
func (coord HexCoord) Coord() HexCoord { return coord }

func NewHexCoord(i, j int) HexCoord {
	return HexCoord{i, j}
}

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

// Facilitates the different types of tiles that may appear at a map coordinate.
// While the text-file representation overlaps multiple tile types with the same
// glyph, with this representation each crate is separate from each goal, having
// distinct object representations for each (and the player).  Meanwhile, walls
// are temporary (they become implicit in the JSON representation) and likewise,
// floors are there to assert which coordinates are valid for goal/crate/player.
type TileType string

const (
	TILE_FLOOR  TileType = "FLOOR"
	TILE_WALL   TileType = "WALL"
	TILE_GOAL   TileType = "GOAL"
	TILE_CRATE  TileType = "CRATE"
	TILE_PLAYER TileType = "PLAYER"
)

type Floor struct{ HexCoord }
type Wall struct{ HexCoord }
type Goal struct{ HexCoord }
type Crate struct{ HexCoord }
type Player struct{ HexCoord }

func (tile Floor) Type() TileType  { return TILE_FLOOR }
func (tile Wall) Type() TileType   { return TILE_WALL }
func (tile Goal) Type() TileType   { return TILE_GOAL }
func (tile Crate) Type() TileType  { return TILE_CRATE }
func (tile Player) Type() TileType { return TILE_PLAYER }

// Adds the provided tiles to this puzzle, some validation is performed.
func (puzzle *Puzzle) AddTiles(tiles []Tile) error {
	for _, tile := range tiles {
		switch tile.Type() {
		case TILE_FLOOR:
			puzzle.Terrain = append(puzzle.Terrain, tile.Coord())
		case TILE_WALL:
			// Walls can be inferred by positions not existing in the terrain,
			// they do not even need to be rendered, but may be inferred from
			// any tile adjacent to a Terrain coordinate that is not in Terrain.
			continue
		case TILE_GOAL:
			puzzle.Init.Goals = append(puzzle.Init.Goals, tile.Coord())
		case TILE_CRATE:
			puzzle.Init.Crates = append(puzzle.Init.Crates, tile.Coord())
		case TILE_PLAYER:
			puzzle.Init.Player = tile.Coord()
		}
	}
	return nil
}
