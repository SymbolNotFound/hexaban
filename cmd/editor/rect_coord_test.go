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

package main

import (
	"reflect"
	"testing"

	"github.com/SymbolNotFound/hexoban"
)

func TestRectCoord_ToHex(t *testing.T) {
	coord := hexoban.NewHexCoord
	tests := []struct {
		name  string
		coord RectCoord
		want  hexoban.HexCoord
	}{
		{"zero", RectCoord{0, 0}, coord(0, 0)},
		{"simple", RectCoord{1, 1}, coord(1, 1)},
		{"forty-two", RectCoord{0, 4}, coord(4, 2)},
		{"hex 5 6", RectCoord{7, 5}, coord(5, 6)},
		{"hex 2 3", RectCoord{4, 2}, coord(2, 3)},
		{"hex 7 9", RectCoord{11, 7}, coord(7, 9)},
		{"hex 7 5", RectCoord{3, 7}, coord(7, 5)},
		{"hex 7 7", RectCoord{7, 7}, coord(7, 7)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.coord.ToHex(); got.I() != tt.want.I() || got.J() != tt.want.J() {
				t.Errorf("RectCoord.ToHex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHexToRect(t *testing.T) {
	coord := hexoban.NewHexCoord
	tests := []struct {
		name   string
		coord  hexoban.HexCoord
		expect RectCoord
	}{
		{"zero", coord(0, 0), RectCoord{0, 0}},
		{"simple", coord(1, 1), RectCoord{1, 1}},
		{"forty-two", coord(4, 2), RectCoord{0, 4}},
		{"hex 5 6", coord(5, 6), RectCoord{7, 5}},
		{"hex 2 3", coord(2, 3), RectCoord{4, 2}},
		{"hex 7 9", coord(7, 9), RectCoord{11, 7}},
		{"hex 7 5", coord(7, 5), RectCoord{3, 7}},
		{"hex 7 7", coord(7, 7), RectCoord{7, 7}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HexToRect(tt.coord, 0, 0); !reflect.DeepEqual(got, tt.expect) {
				t.Errorf("HexToRect() = %v, expected %v", got, tt.expect)
			}
		})
	}
}
