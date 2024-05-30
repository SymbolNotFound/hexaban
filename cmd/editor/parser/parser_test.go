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
// github:SymbolNotFound/hexoban/cmd/editor/parser/parser_test.go

package parser

import (
	"bufio"
	"reflect"
	"strings"
	"testing"

	"github.com/SymbolNotFound/hexoban/puzzle"
)

func TestParsePuzzleDefinition(t *testing.T) {
	at := puzzle.NewHexCoord
	tests := []struct {
		name     string
		text     string
		expected puzzle.Puzzle
	}{
		{
			"small room",
			"  # #\n #   # \n  # #",
			puzzle.Puzzle{
				Terrain: []puzzle.HexCoord{at(2, -1)},
			},
		},
		{
			"odd small room",
			" # # \n#   #\n # #\n",
			puzzle.Puzzle{
				Terrain: []puzzle.HexCoord{at(2, 0)},
			},
		},
		{
			"dws001",
			`
    # # #
   #     #
    #     #
   #   . #
  #   .   #
 #   $ $   #
  # # # *   #
     # @   #
      # # #

; credit for this puzzle goes to David W. Skinner`,
			puzzle.Puzzle{
				Terrain: []puzzle.HexCoord{
					/*    # # #     */
					/*   #     #    */ at(3, -2), at(4, -3),
					/*    #     #   */ at(4, -2), at(5, -3),
					/*   #   . #    */ at(4, -1), at(5, -2),
					/*  #   .   #   */ at(4, 0), at(5, -1), at(6, -2),
					/* #   $ $   #  */ at(4, 1), at(5, 0), at(6, -1), at(7, -2),
					/*  # # # *   # */ at(7, -1), at(8, -2),
					/*     # @   #  */ at(7, 0), at(8, -1),
					/*      # # #   */
				},
				Init: puzzle.Init{
					Goals: []puzzle.HexCoord{
						at(5, -2), at(5, -1), at(7, -1),
					},
					Crates: []puzzle.HexCoord{
						at(5, 0), at(6, -1), at(7, -1),
					},
					Ichiban: at(7, 0),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.text))
			puzzle := puzzle.Puzzle{}
			if err := ParsePuzzleDefinition(reader, &puzzle); err != nil {
				t.Errorf("ParsePuzzleDefinition() error = %v", err)
				return
			}
			if !reflect.DeepEqual(puzzle, tt.expected) {
				t.Errorf("ParsePuzzleDefinition() = %v, expected %v", puzzle, tt.expected)
			}
		})
	}
}
