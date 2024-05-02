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
// github:SymbolNotFound/hexaban/hexaban_test.go

package hexaban

import (
	"testing"
)

func TestRectCoord_ToHex(t *testing.T) {
	tests := []struct {
		name  string
		coord RectCoord
		want  HexCoord
	}{
		{"zero", RectCoord{0, 0}, HexCoord{0, 0}},
		{"i", RectCoord{1, 0}, HexCoord{1, 0}},
		{"i+j", RectCoord{0, 1}, HexCoord{1, 1}},
		{"col 4", RectCoord{4, 0}, HexCoord{2, -2}},
		{"3i - j", RectCoord{3, 1}, HexCoord{3, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.coord.ToHex(); got.i != tt.want.i || got.j != tt.want.j {
				t.Errorf("RectCoord.ToHex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRectCoord_ToHexCenteredAt(t *testing.T) {
	tests := []struct {
		name   string
		coord  RectCoord
		center HexCoord
		want   HexCoord
	}{
		{"zero", RectCoord{0, 0}, HexCoord{0, 0}, HexCoord{0, 0}},
		{"i", RectCoord{1, 0}, HexCoord{0, 0}, HexCoord{1, 0}},
		{"i (recentered)", RectCoord{1, 0}, HexCoord{2, 2}, HexCoord{-1, -2}},
		{"i+j (recentered)", RectCoord{0, 1}, HexCoord{3, 2}, HexCoord{-2, -1}},
		{"3i - j at center", RectCoord{6, 1}, HexCoord{3, -1}, HexCoord{1, -1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.coord.ToHexCenteredAt(tt.center); got.i != tt.want.i || got.j != tt.want.j {
				t.Errorf("RectCoord.ToHex() = %v, want %v", got, tt.want)
			}
		})
	}
}
