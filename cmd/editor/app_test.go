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
// github:SymbolNotFound/hexaban/cmd/editor/app_test.go

package main

import (
	"bufio"
	"reflect"
	"strings"
	"testing"

	"github.com/SymbolNotFound/hexaban"
)

func Test_app(t *testing.T) {
	at := hexaban.NewHexCoord
	tests := []struct {
		name     string
		inputs   string
		expected hexaban.Puzzle
	}{
		{
			"quick small",
			`
Anonymous
Prison Cell
Source
testdata

 # #
# @ #
 # #
`,
			hexaban.Puzzle{
				Author:  "Anonymous",
				Name:    "Prison Cell",
				Source:  "testdata",
				Terrain: []hexaban.HexCoord{at(2, 0)},
				Init: hexaban.Init{
					Ichiban: at(2, 0),
				},
			},
		},
		{
			"novel entry",
			`Kevin Damm
One-Way Mirror
Source
https://hexaban.com/levels/damm/011.json

       # #
    # #   # #
 # #     $   #
# .   $ #   @ #
 # #     #     #
#   * $ . $ *   #
 #     #     # #
  #   . # .   . #
   #   $     # #
    # #   # #
       # #
`,
			hexaban.Puzzle{
				Author: "Kevin Damm",
				Name:   "One-Way Mirror",
				Source: "https://hexaban.com/levels/damm/011.json",
				Terrain: []hexaban.HexCoord{
					/*        # #        */
					/*     # #   # #     */ at(5, -3),
					/*  # #     $   #    */ at(4, -1), at(5, -2), at(6, -3), at(7, -4),
					/* # .   $ #   @ #   */ at(3, 1), at(4, 0), at(5, -1), at(7, -3), at(8, -4),
					/*  # #     #     #  */ at(5, 0), at(6, -1), at(8, -3), at(9, -4),
					/* #   * $ . $ *   # */ at(4, 2), at(5, 1), at(6, 0), at(7, -1), at(8, -2), at(9, -3), at(10, -4),
					/*  #     #     # #  */ at(5, 2), at(6, 1), at(8, -1), at(9, -2),
					/*   #   . # .   . # */ at(6, 2), at(7, 1), at(9, -1), at(10, -2), at(11, -3),
					/*    #   $     # #  */ at(7, 2), at(8, 1), at(9, 0), at(10, -1),
					/*     # #   # #     */ at(9, 1),
					/*        # #        */
				},
				Init: hexaban.Init{
					Goals: []hexaban.HexCoord{
						at(3, 1),
						at(5, 1),
						at(7, -1),
						at(9, -3),
						at(7, 1),
						at(9, -1),
						at(11, -3),
					},
					Crates: []hexaban.HexCoord{
						at(6, -3),
						at(5, -1),
						at(5, 1),
						at(6, 0),
						at(8, -2),
						at(9, -3),
						at(8, 1),
					},
					Ichiban: at(8, -4),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reader *bufio.Reader = bufio.NewReader(strings.NewReader(tt.inputs))
			var puzzle hexaban.Puzzle = hexaban.Puzzle{}

			app(reader, &puzzle)
			if !reflect.DeepEqual(puzzle, tt.expected) {
				t.Errorf("puzzle after app():\n%v\n\n  expected:\n%v", puzzle, tt.expected)
			}
		})
	}
}
