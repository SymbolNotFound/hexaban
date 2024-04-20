// Hexaban (some know it as Hexoban to match the 'o' in Sokoban).
//
// Structures for representing the initial state and sequence of moves
// for the solitaire game that is a hexagonal representation of the classic
// known as Sokoban.

package main

import (
	"encoding/json"
)

type Collection struct {
	Puzzles []HexabanPuzzle `json:"puzzles"`
	Source  string          `json:"source"`
	Author  string          `json:"author,omitempty"`
}

type HexabanPuzzle struct {
	Identity string      `json:"id"`
	Author   string      `json:"author"`
	Name     string      `json:"name"`
	Terrain  []HexCoord  `json:"terrain"`
	Init     HexabanInit `json:"init"`
}

type HexabanInit struct {
	Walls  []HexCoord `json:"walls,omitempty"`
	Goals  []HexCoord `json:"goals"`
	Crates []HexCoord `json:"crates"`
	Player HexCoord   `json:"ichiban"`
}

type RectCoord struct {
	col uint
	row uint
}

// This assumes that the rectangular coordinates collapse each pair of rows
// into a common row coordinate, aligning the columns, and that hexes have
// their orientation such that up/down movement is possible (the other four
// directions are left/right and north/south).
// The HexCoord being converted into has its unit vectors corresponding to
// right (downward-right in pixel coordinates) and south (or downward-left
// in pixel coords).
func (coord RectCoord) ToHex() HexCoord {
	rHalf := coord.row >> 1
	rOdd := coord.row & 1
	cHalf := coord.col >> 1
	cOdd := coord.col & 1

	// A bit of back-of-the-envelope math can derive these
	// from the staggered rectilinear coordinates:
	return HexCoord{
		i: int(cHalf + rHalf + rOdd),
		j: int(rHalf - cHalf - cOdd),
	}
}

type HexCoord struct {
	i int
	j int
}

func (coord HexCoord) I() int { return coord.i }
func (coord HexCoord) J() int { return coord.j }

func (coord HexCoord) MarshalJSON() ([]byte, error) {
	asArray := []int{coord.i, coord.j}
	return json.Marshal(asArray)
}
