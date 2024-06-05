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
// github:SymbolNotFound/hexoban/cmd/editor/rect_coord.go

package main

import "github.com/SymbolNotFound/hexoban/puzzle"

// This RectCoord is only provided within the `main` package of a CLI tool.
// It should not be used in concert with other clients of the hexoban library.
//
// If using it, convert it into an axial hex coordinate using RectCoord.ToHex().
type RectCoord struct {
	col uint
	row uint
}

// Converts a rectangular coordinate into the equivalent hexagonal coordinate.
// Assumes that the (0, 0) center for both systems exists at top-left position.
//
// This uses a rectangular coordinate system where each pair of rows is combined
// into a common row coordinate, aligning the columns.  Each hexagonal tile has
// their orientation such that left/right movement is possible.  The other four
// directions are up/down and back/forth, depicted here.
//
// .                                                 significant spaces shown as
// .   back /^\  up                                  '_' characters for clarity:
// .       /   \           /==================== ...
// .  left |   | right     | 0,0     2,0     4,0           # # #
// .       \   /           |     1,0     3,0     ...      # _ . #     3,0 5,0
// .   down \v/ forth      | 0,1     2,1     4,1         # _ $ @ #  2,1 4,1 6,1
// .                       |     1,1     3,1     ...      # _ # #     3,1
// .                       | 0,2     2,2     4,2           # #             (c,r)
// .  RectCoord{4, 2}      |     ...     ...
// .      =~=              \ ...                       each rect_coord{col, row}
// .   HexCoord{4, 4}                                  adjacent to six neighbors
// .            /  -
//
// The selection of cardinal directions is arbitrary, this is good as any.  But
// notice that sometimes a column difference of |1| is a movement `forward`, and
// sometimes it is a movement `up`, but it depends on the column being moved
// from.  Their position when projected onto a canvas is higher or lower based
// on their column (despite higher/lower dimension being a row-order feature!).
//
// For this reason, the rectilinear format is abandoned in favor of a coordinate
// system based on `down` and `right` directions, as (i, j) in `HexCoord`.  The
// RectCoord is used for parsing and converting maps defined elsewhere using this
// format.  There are many instances of this map format found in the wild because
// rectangular coordinates are easier to map to memory and can be printed to
// line-based consoles easily, and existing tools used the same for rectangular
// grids (except, of course, those could be completely aligned with x and y).
//
// Down (i) and Right (j) become the unit vectors in this representation, and
// positioned (along with the naming of other coordinates) to align with the
// (i, j) from other game definitions, such that a transformation matrix is able
// to represent all rotations and reflections available for hexagonal grids.
// This can be used by frontend clients to provide alternate views of the same
// puzzle, as can be seen in layout.ts of the frontend (/webapp/src/hexgrid).
func (coord RectCoord) ToHex() puzzle.HexCoord {
	// We can reconstruct (i, j) from rows and columns by accumulating
	// its partial representations per-column and per-row.
	return puzzle.NewHexCoord(
		int(coord.row),
		int((coord.row+coord.col)>>1))
}

// Convenient inverse of the above function, to get line and column coordinates
// A translation from an origin (top-left position) is performed, all returned
// RectCoord coordinates will be added to fromCol, fromRow.
// All intervals on the same line are 2, all odd rows are odd-columned.
func HexToRect(hex puzzle.HexCoord, fromCol, fromRow int) RectCoord {
	return RectCoord{
		col: uint((hex.J() << 1) - hex.I() + fromCol),
		row: uint(hex.I() + fromRow),
	}
}
