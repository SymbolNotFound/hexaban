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
// github:SymbolNotFound/hexoban/cmd/editor/fmt_test.go

package main

import (
	"testing"

	"github.com/SymbolNotFound/hexoban"
)

func TestMapString(t *testing.T) {
	at := hexoban.NewHexCoord
	tests := []struct {
		name   string
		p      hexoban.Puzzle
		expect string
	}{
		{
			"quick small",
			hexoban.Puzzle{
				Author: "Anonymous",
				Title:  "Treasure Room",
				Source: "testdata",
				Terrain: []hexoban.HexCoord{
					at(1, 2), at(1, 3), at(2, 2), at(2, 3), at(2, 4),
				},
				Init: hexoban.Init{
					Goals:   []hexoban.HexCoord{at(2, 4)},
					Crates:  []hexoban.HexCoord{at(2, 3)},
					Ichiban: at(2, 2),
				},
			},
			"  # # #\n #     #\n# @ $ . #\n # # # #",
		},
		{
			"more complete",
			hexoban.Puzzle{
				Author: "Kevin Damm",
				Title:  "One-Way Mirror",
				Source: "https://hexoban.com/levels/damm/011.json",
				Terrain: []hexoban.HexCoord{
					/*        # #        */
					/*     # #   # #     */ at(2, 5),
					/*  # #     $   #    */ at(3, 4), at(3, 5), at(3, 6), at(3, 7),
					/* # .   $ #   @ #   */ at(4, 3), at(4, 4), at(4, 5), at(4, 7), at(4, 8),
					/*  # #     #     #  */ at(5, 5), at(5, 6), at(5, 8), at(5, 9),
					/* #   * $ . $ *   # */ at(6, 4), at(6, 5), at(6, 6), at(6, 7), at(6, 8), at(6, 9), at(6, 10),
					/*  #     #     # #  */ at(7, 5), at(7, 6), at(7, 8), at(7, 9),
					/*   #   . # .   . # */ at(8, 6), at(8, 7), at(8, 9), at(8, 10), at(8, 11),
					/*    #   $     # #  */ at(9, 7), at(9, 8), at(9, 9), at(9, 10),
					/*     # #   # #     */ at(10, 9),
					/*        # #        */
				},
				Init: hexoban.Init{
					Goals: []hexoban.HexCoord{
						at(4, 3), at(6, 5), at(6, 7), at(6, 9), at(8, 7), at(8, 9), at(8, 11),
					},
					Crates: []hexoban.HexCoord{
						at(3, 6), at(4, 5), at(6, 5), at(6, 6), at(6, 8), at(6, 9), at(9, 8),
					},
					Ichiban: at(4, 8),
				},
			},
			`
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
       # #`[1:], // trim initial newline
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := MapString(tt.p); err != nil || got != tt.expect {
				if err != nil {
					t.Error(err)
				} else {
					t.Errorf("MapString()\n%v\nexpected:\n%v", got, tt.expect)
				}
			}
		})
	}
}
