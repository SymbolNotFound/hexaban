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
// github:SymbolNotFound/hexaban/puzzle/hex_coord.go

package puzzle

import "encoding/json"

// The HexCoord has a basis of two unit vectors (i, j) corresponding to
// forward (downward-right in pixel coordinates) and down (or downward-left
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
