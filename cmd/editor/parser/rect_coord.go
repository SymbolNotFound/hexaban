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

package parser

import "github.com/SymbolNotFound/hexoban/puzzle"

// This RectCoord is only provided within the main package of a CLI tool because
// it should not be used in common with any other clients of the hexoban library.
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
// into a common row coordinate, aligning the columns, and that hexes have
// their orientation such that left/right movement is possible (the other four
// directions are up/down and back/forward) depicted here.
//
// .                                                 significant spaces shown as
// .                                                 '_' characters for clarity:
// .                       /==================== ...
// .   back /^\  up        | 0,0     2,0     4,0             # # #
// .       /   \           |     1,0     3,0     ...        # _ . #
// .  left |   | right     | 0,1     2,1     4,1           # _ $ @ #
// .       \   /           |     1,1     3,1     ...        # _ # #
// .   down \v/ forward    | 0,2     2,2     4,2             # #
// .                       |     ...     ...
// .                       \ ...                       each rect_coord{col, row}
// .                                                   adjacent to six neighbors
// RectCoord{6, 1} =~= HexCoord{4, -2}
//
// The selection of cardinal directions is arbitrary, this is good as any.  But
// notice that sometimes a column difference of |1| is a movement `forward`, and
// sometimes it is a movement `up`, but it depends on the column being moved
// from.  Their position when projected onto a canvas is higher or lower based
// on their column (despite higher/lower dimension being a row-order feature!).
//
// For this reason, the rectilinear format is abandoned in favor of a coordinate
// system based on `forward` and `down` directions, as (i, j) in `HexCoord`.  The
// RectCoord is used for parsing and converting maps defined elsewhere using this
// format.  There are many instances of this map format found in the wild because
// rectangular coordinates are easier to map to memory and can be printed to
// line-based consoles easily, and existing tools used the same for rectangular
// grids (except, of course, those could be completely aligned with x and y).
func (coord RectCoord) ToHex() puzzle.HexCoord {
	cHalf := coord.col >> 1
	cOdd := coord.col & 1

	// Each adjacent column (two column distance) is half i, half -j (with odd
	// columns adding one more i).  Each row is one i and one j.
	// We can reconstruct (i, j) from rows and columns by accumulating these
	// partial representations per-column and per-row.
	return puzzle.NewHexCoord(
		int(cHalf+cOdd+coord.row),
		int(coord.row-cHalf))
}
