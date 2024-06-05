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
// github:SymbolNotFound/hexoban/hexoban.go

package hexoban

import (
	"encoding/json"
	"fmt"
)

// Structures for representing the initial state and sequence of moves
// for the solitaire game that is a hexagonal variant of the classic
// puzzle known as Sokoban ("warehouse worker," the crate pusher).

// The puzzle's serialized structure.
type Puzzle struct {
	Identity   string     `json:"id"`
	Title      string     `json:"title"`
	Author     string     `json:"author"`
	Source     string     `json:"source"`
	Difficulty int        `json:"difficulty,omitempty"`
	Terrain    []HexCoord `json:"terrain"`
	Init       Init       `json:"init"`
}

type Init struct {
	// All goals and crates need to be on coordinates where terrain exists.
	Goals  []HexCoord `json:"goals"`
	Crates []HexCoord `json:"crates"`

	// The player's initial position may be anywhere reachable from here.
	Ichiban HexCoord `json:"ichiban,omitempty"` // Default is (0, 0)
}

// Runs some basic validation over the puzzle definition.
func (puzzle Puzzle) Validate() error {
	// TODO
	return nil
}

// Satisfies the fmt.Stringer interface, presently just prints the json-marshaled version.
func (puzzle Puzzle) String() string {
	output, err := json.MarshalIndent(puzzle, "  ", "")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(output)
}

// Provides additional info and statistics about the puzzle.
func (puzzle Puzzle) Info() string {
	info := fmt.Sprintf("%s by %s (%s)\nlen terrain, num objectives (crates|goals), ...",
		puzzle.Title, puzzle.Author, puzzle.Identity)
	// TODO
	return info
}

// The HexCoord has a basis of two unit vectors (i, j) corresponding to
// forward (downward-right in pixel coordinates) and down (or downward-left
// in pixel coords).  This results in a back-to-front ordering if sorted by
// (i + j), with ties broken by ascending `i`.  It is also a more suitable
// representation for translations, rotations and other affine transformations.
//
// The representation is serialized into JSON as a two-element array, [i, j].
type HexCoord struct {
	i int // `json:[-1]`
	j int // `json:[0]`
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
	array1D := make([]int, 2)
	if err := json.Unmarshal(encoded, &array1D); err != nil {
		return err
	}

	coord.i = array1D[0]
	coord.j = array1D[1]
	return nil
}

// An enum-like interface for a tile's type, bound to its position in a grid.
type Tile interface {
	I() int
	J() int
	Coord() HexCoord
	Type() TileType
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
			puzzle.Init.Ichiban = tile.Coord()
		}
	}
	return nil
}
