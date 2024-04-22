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

// This assumes that the rectangular coordinates collapse each pair of rows
// into a common row coordinate, aligning the columns, and that hexes have
// their orientation such that up/down movement is possible (the other four
// directions are left/right and north/south).
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
