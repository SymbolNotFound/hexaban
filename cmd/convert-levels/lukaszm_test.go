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
// github:SymbolNotFound/hexaban/cmd/convert-levels/dws_test.go

package main

import (
	"reflect"
	"testing"

	"github.com/SymbolNotFound/hexaban"
)

func Test_convertLukaszM(t *testing.T) {
	// Partial representation of DWS collection.
	var text string = `Author: LukaszM

; Level 01
 # # # #
# .   @ #
 # # $   #
    #   #
     # #
Difficulty: 1

; Level 02
 # # # # #
# .   $   #
 # @ $   . #
  # # # # #
Difficulty: 1

; Level 30
    # # # # # #
   # . . . .   #
  #       $   #
 #   # # #   #
#     $   $ #
 # #     $ @ #
    # # #   #
         # #
Difficulty: 8`

	// The details here are abritrary for the purpose of testing, but
	// in recognition of the designer's effort this is accurate.
	collection := hexaban.Collection{
		Source: "https://play.fancade.com/5FA6BCFD16EB8B3B",
		Author: "LukaszM",
	}

	// for convenience, alias the coordinate constructor
	at := hexaban.NewHexCoord

	expected := []hexaban.Puzzle{
		{
			Identity: "LukaszM/01",
			Name:     "Level 01",
			Author:   collection.Author,
			Source:   collection.Source,
			Terrain: []hexaban.HexCoord{
				/*  # # # #   */
				/* # .   @ #  */ at(2, 0), at(3, -1), at(4, -2),
				/*  # # $   # */ at(4, -1), at(5, -2),
				/*     #   #  */ at(5, -1),
				/*      # #   */
			},
			Init: hexaban.Init{
				Goals: []hexaban.HexCoord{
					at(2, 0),
				},
				Crates: []hexaban.HexCoord{
					at(4, -1),
				},
				Ichiban: at(4, -2),
			},
			Difficulty: 1,
		},
		{
			Identity: "LukaszM/02",
			Name:     "Level 02",
			Author:   collection.Author,
			Source:   collection.Source,
			Terrain: []hexaban.HexCoord{
				/*  # # # # #   */
				/* # .   $   #  */ at(2, 0), at(3, -1), at(4, -2), at(5, -3),
				/*  # @ $   . # */ at(3, 0), at(4, -1), at(5, -2), at(6, -3),
				/*   # # # # #  */
			},
			Init: hexaban.Init{
				Goals: []hexaban.HexCoord{
					at(2, 0), at(6, -3),
				},
				Crates: []hexaban.HexCoord{
					at(4, -2), at(4, -1),
				},
				Ichiban: at(3, 0),
			},
			Difficulty: 1,
		},
		{
			Identity: "LukaszM/30",
			Name:     "Level 30",
			Author:   collection.Author,
			Source:   collection.Source,
			Terrain: []hexaban.HexCoord{
				/*     # # # # # #  */
				/*    # . . . .   # */ at(3, -2), at(4, -3), at(5, -4), at(6, -5), at(7, -6),
				/*   #       $   #  */ at(3, -1), at(4, -2), at(5, -3), at(6, -4), at(7, -5),
				/*  #   # # #   #   */ at(3, 0), at(7, -4),
				/* #     $   $ #    */ at(3, 1), at(4, 0), at(5, -1), at(6, -2), at(7, -3),
				/*  # #     $ @ #   */ at(5, 0), at(6, -1), at(7, -2), at(8, -3),
				/*     # # #   #    */ at(8, -2),
				/*          # #     */
			},
			Init: hexaban.Init{
				Goals: []hexaban.HexCoord{
					at(3, -2), at(4, -3), at(5, -4), at(6, -5),
				},
				Crates: []hexaban.HexCoord{
					at(6, -4), at(5, -1), at(7, -3), at(7, -2),
				},
				Ichiban: at(8, -3),
			},
			Difficulty: 8,
		},
	}

	t.Run("Convert Lukaszm puzzles to hex/json", func(t *testing.T) {
		got, err := convertLukaszM([]byte(text), collection)
		if err != nil {
			t.Errorf("convertLukaszM() error = %v", err)
			return
		}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("convertLukaszM() = %v, expected %v", got, expected)
		}
	})
}
