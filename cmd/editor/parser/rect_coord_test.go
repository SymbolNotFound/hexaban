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
// github:SymbolNotFound/hexoban/cmd/editor/rectcoord_test.go

package parser

import (
	"testing"

	"github.com/SymbolNotFound/hexoban/puzzle"
)

func TestRectCoord_ToHex(t *testing.T) {
	tests := []struct {
		name  string
		coord RectCoord
		want  puzzle.HexCoord
	}{
		{"zero", RectCoord{0, 0}, puzzle.NewHexCoord(0, 0)},
		{"i", RectCoord{1, 0}, puzzle.NewHexCoord(1, 0)},
		{"i+j", RectCoord{0, 1}, puzzle.NewHexCoord(1, 1)},
		{"col 4", RectCoord{4, 0}, puzzle.NewHexCoord(2, -2)},
		{"3i - j", RectCoord{3, 1}, puzzle.NewHexCoord(3, 0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.coord.ToHex(); got.I() != tt.want.I() || got.J() != tt.want.J() {
				t.Errorf("RectCoord.ToHex() = %v, want %v", got, tt.want)
			}
		})
	}
}
