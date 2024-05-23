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
// github:SymbolNotFound/hexaban/puzzle/tile.go

package puzzle

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
