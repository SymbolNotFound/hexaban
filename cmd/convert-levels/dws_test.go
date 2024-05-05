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

func Test_convertDWS(t *testing.T) {
	// Partial representation of DWS collection.
	var text string = `; Hexobans by David W. Skinner

"dws001"
    # # #
   #     # 
    #     # 
   #   . # 
  #   .   # 
 #   $ $   # 
  # # # *   # 
     # @   # 
      # # #

"dws002"
              # # #
       # # # #     #
      #     $     #
     #   $ $   $ #
  # #   #   # $ #
 #     #   # #   # # # #
#     #   # #   #   . . #
 #   $             + . . #
  # # $ #   #   #   . . #
   # #     # # # # # # #
      # # #

"dws005"
      # # #
     #     #
      # @   #
     # #     #
    #   * *   #
   #   * $ *   #
  #     . *     #
 #   * * # . *   #
#   * $ . * $ *   #
 #   * *   * *   #
  #             #
   # # # # # # #
`
	// The details here are abritrary for the purpose of testing, but
	// in recognition of the designer's effort this is accurate.
	collection := hexaban.Collection{
		Source: "http://users.bentonrea.com/~sasquatch/sokoban/hex.html",
		Author: "David W. Skinner",
	}

	// for convenience, alias the coordinate constructor
	at := hexaban.NewHexCoord

	expected := []hexaban.Puzzle{
		{
			Identity: "DWS/001",
			Name:     "001",
			Author:   "David W. Skinner",
			Source:   collection.Source,
			Terrain: []hexaban.HexCoord{
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
			Init: hexaban.Init{
				Goals: []hexaban.HexCoord{
					at(5, -2), at(5, -1), at(7, -1),
				},
				Crates: []hexaban.HexCoord{
					at(5, 0), at(6, -1), at(7, -1),
				},
				Player: at(7, 0),
			},
		},
		{
			Identity: "DWS/002",
			Name:     "002",
			Author:   "David W. Skinner",
			Source:   collection.Source,
			Terrain: []hexaban.HexCoord{
				// Puzzle has been slightly modified to exercise the '+' glyph,
				// the complexity of the puzzle remains the same.
				// An important feature of this puzzle are several trailing spaces
				// that look like floor but invalid after last wall of each row/line.
				/*               # # #        */
				/*        # # # #     #       */ at(8, -7), at(9, -8),
				/*       #     $     #        */ at(5, -3), at(6, -4), at(7, -5), at(8, -6), at(9, -7),
				/*      #   $ $   $ #         */ at(5, -2), at(6, -3), at(7, -4), at(8, -5), at(9, -6),
				/*   # #   #   # $ #          */ at(5, -1), at(7, -3), at(9, -5),
				/*  #     #   # #   # # # #   */ at(4, 1), at(5, 0), at(7, -2), at(10, -5),
				/* #     #   # #   #   . . #  */ at(4, 2), at(5, 1), at(7, -1), at(10, -4), at(12, -6), at(13, -7), at(14, -8),
				/*  #   $             + . . # */ at(5, 2), at(6, 1), at(7, 0), at(8, -1), at(9, -2), at(10, -3), at(11, -4), at(12, -5), at(13, -6), at(14, -7), at(15, -8),
				/*   # # $ #   #   #   . . #  */ at(7, 1), at(9, -1), at(11, -3), at(13, -5), at(14, -6), at(15, -7),
				/*    # #     # # # # # # #   */ at(8, 1), at(9, 0),
				/*       # # #                */
			},
			Init: hexaban.Init{
				Goals: []hexaban.HexCoord{
					at(13, -7), at(14, -8), at(13, -6), at(14, -7), at(15, -8), at(14, -6), at(15, -7),
				},
				Crates: []hexaban.HexCoord{
					at(7, -5), at(6, -3), at(7, -4), at(9, -6), at(9, -5), at(6, 1), at(7, 1),
				},
				Player: at(13, -6),
			},
		},
		{
			Identity: "DWS/005",
			Name:     "005",
			Author:   "David W. Skinner",
			Source:   collection.Source,
			Terrain: []hexaban.HexCoord{
				/*       # # #         */
				/*      #     #        */ at(4, -3), at(5, -4),
				/*       # @   #       */ at(5, -3), at(6, -4),
				/*      # #     #      */ at(6, -3), at(7, -4),
				/*     #   * *   #     */ at(5, -1), at(6, -2), at(7, -3), at(8, -4),
				/*    #   * $ *   #    */ at(5, 0), at(6, -1), at(7, -2), at(8, -3), at(9, -4),
				/*   #     . *     #   */ at(5, 1), at(6, 0), at(7, -1), at(8, -2), at(9, -3), at(10, -4),
				/*  #   * * # . *   #  */ at(5, 2), at(6, 1), at(7, 0), at(9, -2), at(10, -3), at(11, -4),
				/* #   * $ . * $ *   # */ at(5, 3), at(6, 2), at(7, 1), at(8, 0), at(9, -1), at(10, -2), at(11, -3), at(12, -4),
				/*  #   * *   * *   #  */ at(6, 3), at(7, 2), at(8, 1), at(9, 0), at(10, -1), at(11, -2), at(12, -3),
				/*   #             #   */ at(7, 3), at(8, 2), at(9, 1), at(10, 0), at(11, -1), at(12, -2),
				/*    # # # # # # #    */
			},
			Init: hexaban.Init{
				Goals: []hexaban.HexCoord{
					at(6, -2), at(7, -3), at(6, -1), at(8, -3), at(7, -1), at(8, -2),
					at(6, 1), at(7, 0), at(9, -2), at(10, -3), at(6, 2), at(8, 0),
					at(9, -1), at(11, -3), at(7, 2), at(8, 1), at(10, -1), at(11, -2),
				},
				Crates: []hexaban.HexCoord{
					at(6, -2), at(7, -3), at(6, -1), at(7, -2), at(8, -3), at(8, -2),
					at(6, 1), at(7, 0), at(10, -3), at(6, 2), at(7, 1), at(9, -1),
					at(10, -2), at(11, -3), at(7, 2), at(8, 1), at(10, -1), at(11, -2),
				},
				Player: at(5, -3),
			},
		},
	}

	t.Run("Convert DWS rect to hex", func(t *testing.T) {
		got, err := convertDWS([]byte(text), collection)
		if err != nil {
			t.Errorf("convertDWS() error = %v", err)
			return
		}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("convertDWS() = %v, expected %v", got, expected)
		}
	})
}
