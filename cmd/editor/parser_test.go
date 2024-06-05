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
// github:SymbolNotFound/hexoban/cmd/editor/parser_test.go

package main

import (
	"bufio"
	"reflect"
	"strings"
	"testing"

	"github.com/SymbolNotFound/hexoban"
)

func TestParsePuzzleDefinition(t *testing.T) {
	at := hexoban.NewHexCoord
	tests := []struct {
		name     string
		text     string
		expected hexoban.Puzzle
	}{
		{
			"small room",
			"  # #\n #   # \n  # #",
			hexoban.Puzzle{
				Terrain: []hexoban.HexCoord{at(1, 2)},
			},
		},
		{
			"odd small room",
			" # # \n#   #\n # #\n",
			hexoban.Puzzle{ // TODO parser has a problem if this map is nudged any number of columns over
				Terrain: []hexoban.HexCoord{at(2, 2)},
			},
		},
		{
			"dws001",
			`   # # #
  #     #
   #     #
  #   . #
 #   .   #
#   $ $   #
 # # # *   #
    # @   #
     # # #

; credit for this puzzle goes to David W. Skinner`,
			hexoban.Puzzle{
				Terrain: []hexoban.HexCoord{
					/*    # # #     */
					/*   #     #    */ at(2, 3), at(2, 4),
					/*    #     #   */ at(3, 4), at(3, 5),
					/*   #   . #    */ at(4, 4), at(4, 5),
					/*  #   .   #   */ at(5, 4), at(5, 5), at(5, 6),
					/* #   $ $   #  */ at(6, 4), at(6, 5), at(6, 6), at(6, 7),
					/*  # # # *   # */ at(7, 7), at(7, 8),
					/*     # @   #  */ at(8, 7), at(8, 8),
					/*      # # #   */
				},
				Init: hexoban.Init{
					Goals: []hexoban.HexCoord{
						at(4, 5), at(5, 5), at(7, 7),
					},
					Crates: []hexoban.HexCoord{
						at(6, 5), at(6, 6), at(7, 7),
					},
					Ichiban: at(8, 7),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.text))
			puzzle := hexoban.Puzzle{}
			if err := ParsePuzzleDefinition(reader, &puzzle); err != nil {
				t.Errorf("ParsePuzzleDefinition() error = %v", err)
				return
			}
			if !reflect.DeepEqual(puzzle, tt.expected) {
				result, _ := MapString(puzzle)
				expected, _ := MapString(tt.expected)
				t.Errorf("ParsePuzzleDefinition()\n%s\n,,,\nexpected:\n%s", result, expected)
			}
		})
	}
}
