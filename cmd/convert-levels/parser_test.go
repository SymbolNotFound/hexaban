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
// github:SymbolNotFound/hexaban/cmd/convert-levels/parser_test.go

package main

import (
	"reflect"
	"testing"

	"github.com/SymbolNotFound/hexaban"
)

func TestPuzzleParser_ParseTextGrid(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		expect []Tile
	}{
		{
			"simple", "  # #\n #   #\n  # #\n",
			[]Tile{
				hexaban.Wall{HexCoord: hexaban.NewHexCoord(1, -1)},
				hexaban.Wall{HexCoord: hexaban.NewHexCoord(2, -2)},
				hexaban.Wall{HexCoord: hexaban.NewHexCoord(1, 0)},
				hexaban.Floor{HexCoord: hexaban.NewHexCoord(2, -1)},
				hexaban.Wall{HexCoord: hexaban.NewHexCoord(3, -2)},
				hexaban.Wall{HexCoord: hexaban.NewHexCoord(2, 0)},
				hexaban.Wall{HexCoord: hexaban.NewHexCoord(3, -1)},
			},
		},
		{
			"simple shifted", "   # #\n  #   #\n   # #",
			// Small box, only one traversible (floor) tile, shifted ij<1, 0>.
			[]Tile{
				hexaban.Wall{HexCoord: hexaban.NewHexCoord(2, -1)},
				hexaban.Wall{HexCoord: hexaban.NewHexCoord(3, -2)},
				hexaban.Wall{HexCoord: hexaban.NewHexCoord(2, 0)},
				hexaban.Floor{HexCoord: hexaban.NewHexCoord(3, -1)},
				hexaban.Wall{HexCoord: hexaban.NewHexCoord(4, -2)},
				hexaban.Wall{HexCoord: hexaban.NewHexCoord(3, 0)},
				hexaban.Wall{HexCoord: hexaban.NewHexCoord(4, -1)},
			},
		},
		{
			"full ending doublenewline", "  # #\n #   #\n  # #\n\n",
			[]Tile{
				hexaban.Wall{HexCoord: hexaban.NewHexCoord(1, -1)},
				hexaban.Wall{HexCoord: hexaban.NewHexCoord(2, -2)},
				hexaban.Wall{HexCoord: hexaban.NewHexCoord(1, 0)},
				hexaban.Floor{HexCoord: hexaban.NewHexCoord(2, -1)},
				hexaban.Wall{HexCoord: hexaban.NewHexCoord(3, -2)},
				hexaban.Wall{HexCoord: hexaban.NewHexCoord(2, 0)},
				hexaban.Wall{HexCoord: hexaban.NewHexCoord(3, -1)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser([]byte(tt.text))
			var tiles []Tile
			var err error
			if tiles, err = parser.ParseTextGrid(); err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(tiles, tt.expect) {
				t.Errorf("PuzzleParser.ParseTextGrid() = %v, want %v", tiles, tt.expect)
			}
		})
	}
}

func TestPuzzleParser_ParseProperty(t *testing.T) {
	puzzle := []byte(`
          # # # 
         #     # 
        # $ .   # 
   # # #   *   # 
  # #   $ . $ # 
 # @ . * . .   # 
  #   $   $ * $ # 
   # # # #   .   # 
          #     # 
           # # #
Author: David Holland
Source: morehex.hsb
`[1:])
	parser := NewParser(puzzle)
	authorKey := "Author"
	authorExpect := "David Holland"

	sourceKey := "Source"
	sourceExpect := "morehex.hsb"

	t.Run("", func(t *testing.T) {
		if got := parser.ParseProperty(authorKey); got != authorExpect {
			t.Errorf("PuzzleParser.ParseProperty() = %v, want %v", got, authorExpect)
		}

		if string(parser.filedata) != `
          # # # 
         #     # 
        # $ .   # 
   # # #   *   # 
  # #   $ . $ # 
 # @ . * . .   # 
  #   $   $ * $ # 
   # # # #   .   # 
          #     # 
           # # #
Source: morehex.hsb
`[1:] {
			t.Errorf("resulting puzzle definition should have its property removed. got:\n%s", parser.filedata[parser.cursor:])
		}

		if got := parser.ParseProperty(sourceKey); got != sourceExpect {
			t.Errorf("PuzzleParser.ParseProperty() = %v, want %v", got, sourceExpect)
		}
	})
}
