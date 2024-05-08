package main

import (
	"reflect"
	"testing"

	"github.com/SymbolNotFound/hexaban"
)

func Test_convertSEVER(t *testing.T) {
	var text string = `; erim_hex1
   # # #
  #     #
 #   $   #
#   $   # # #
 # $ #       #
# #   $ @ #   #
 # $   # #   #
#     # #   #
 #     #   #
  # #   . #
     # .   #
    # . . . #
     # # # #
Author: Erim SEVER

; erim_hex4
    # # # # #
   #         #
  #   * * . $ # #
   #   *   *     #
  #   * * * *     #
 #   *   @   *   #
#     * * * *   #
 #     *   *   #
  # # $ . * *   #
     #         #
      # # # # #
Author: Erim SEVER

; erim_hex_nr5
        # # #
     # #     # #
    #   *       #
   #   $ * * $   #
  #   * .   . *   #
 #   *   * *   *   #
#   $ . * @ * . $   #
 #   *   * *   * * #
#     * .   . *     #
 #   * $ * * $     #
  # #           # #
     # # # # # #
Author: Erim SEVER`

	// The details here are abritrary for the purpose of testing, but
	// in recognition of the designer's effort this is accurate.
	collection := hexaban.Collection{
		Source: "www.erimsever.com/sokoban/Erim_Levels/E_Hexoban.zip",
		Author: "Erim SEVER",
	}
	// Convenience alias for the coordinate constructor.
	at := hexaban.NewHexCoord

	expected := []hexaban.Puzzle{
		{
			Identity: "ErimSEVER/01",
			Name:     "Hex 1",
			Author:   collection.Author,
			Source:   collection.Source,
			Terrain: []hexaban.HexCoord{
				/* 	  # # #        */
				/*   #     #       */ at(3, -1), at(4, -2),
				/*  #   $   #      */ at(3, 0), at(4, -1), at(5, -2),
				/* #   $   # # #   */ at(3, 1), at(4, 0), at(5, -1),
				/*  # $ #       #  */ at(4, 1), at(6, -1), at(7, -2), at(8, -3),
				/* # #   $ @ #   # */ at(5, 1), at(6, 0), at(7, -1), at(9, -3),
				/*  # $   # #   #  */ at(5, 2), at(6, 1), at(9, -2),
				/* #     # #   #   */ at(5, 3), at(6, 2), at(9, -1),
				/*  #     #   #    */ at(6, 3), at(7, 2), at(9, 0),
				/*   # #   . #     */ at(8, 2), at(9, 1),
				/*      # .   #    */ at(9, 2), at(10, 1),
				/*     # . . . #   */ at(9, 3), at(10, 2), at(11, 1),
				/*      # # # #    */
			},
			Init: hexaban.Init{
				Goals: []hexaban.HexCoord{
					at(9, 1), at(9, 2), at(9, 3), at(10, 2), at(11, 1),
				},
				Crates: []hexaban.HexCoord{
					at(4, -1), at(4, 0), at(4, 1), at(6, 0), at(5, 2),
				},
				Ichiban: at(7, -1),
			},
		},
		{
			Identity: "ErimSEVER/04",
			Name:     "Hex 4",
			Author:   collection.Author,
			Source:   collection.Source,
			Terrain: []hexaban.HexCoord{
				/*     # # # # #       */
				/*    #         #      */ at(3, -2), at(4, -3), at(5, -4), at(6, -5),
				/*   #   * * . $ # #   */ at(3, -1), at(4, -2), at(5, -3), at(6, -4), at(7, -5),
				/*    #   *   *     #  */ at(4, -1), at(5, -2), at(6, -3), at(7, -4), at(8, -5), at(9, -6),
				/*   #   * * * *     # */ at(4, 0), at(5, -1), at(6, -2), at(7, -3), at(8, -4), at(9, -5), at(10, -6),
				/*  #   *   @   *   #  */ at(4, 1), at(5, 0), at(6, -1), at(7, -2), at(8, -3), at(9, -4), at(10, -5),
				/* #     * * * *   #   */ at(4, 2), at(5, 1), at(6, 0), at(7, -1), at(8, -2), at(9, -3), at(10, -4),
				/*  #     *   *   #    */ at(5, 2), at(6, 1), at(7, 0), at(8, -1), at(9, -2), at(10, -3),
				/*   # # $ . * *   #   */ at(7, 1), at(8, 0), at(9, -1), at(10, -2), at(11, -3),
				/*      #         #    */ at(8, 1), at(9, 0), at(10, -1), at(11, -2),
				/*       # # # # #     */
			},
			Init: hexaban.Init{
				Goals: []hexaban.HexCoord{
					at(4, -2), at(5, -3), at(6, -4), at(5, -2), at(7, -4), at(5, -1), at(6, -2), at(7, -3), at(8, -4), at(5, 0), at(9, -4), at(6, 0), at(7, -1), at(8, -2), at(9, -3), at(7, 0), at(9, -2), at(8, 0), at(9, -1), at(10, -2),
				},
				Crates: []hexaban.HexCoord{
					at(4, -2), at(5, -3), at(7, -5), at(5, -2), at(7, -4), at(5, -1), at(6, -2), at(7, -3), at(8, -4), at(5, 0), at(9, -4), at(6, 0), at(7, -1), at(8, -2), at(9, -3), at(7, 0), at(9, -2), at(7, 1), at(9, -1), at(10, -2),
				},
				Ichiban: at(7, -2),
			},
		},
		{
			Identity: "ErimSEVER/54",
			Name:     "Hex 54",
			Author:   collection.Author,
			Source:   collection.Source,
			Terrain: []hexaban.HexCoord{
				/*         # # #         */
				/*      # #     # #      */ at(5, -4), at(6, -5),
				/*     #   *       #     */ at(4, -2), at(5, -3), at(6, -4), at(7, -5), at(8, -6),
				/*    #   $ * * $   #    */ at(4, -1), at(5, -2), at(6, -3), at(7, -4), at(8, -5), at(9, -6),
				/*   #   * .   . *   #   */ at(4, 0), at(5, -1), at(6, -2), at(7, -3), at(8, -4), at(9, -5), at(10, -6),
				/*  #   *   * *   *   #  */ at(4, 1), at(5, 0), at(6, -1), at(7, -2), at(8, -3), at(9, -4), at(10, -5), at(11, -6),
				/* #   $ . * @ * . $   # */ at(4, 2), at(5, 1), at(6, 0), at(7, -1), at(8, -2), at(9, -3), at(10, -4), at(11, -5), at(12, -6),
				/*  #   *   * *   * * #  */ at(5, 2), at(6, 1), at(7, 0), at(8, -1), at(9, -2), at(10, -3), at(11, -4), at(12, -5),
				/* #     * .   . *     # */ at(5, 3), at(6, 2), at(7, 1), at(8, 0), at(9, -1), at(10, -2), at(11, -3), at(12, -4), at(13, -5),
				/*  #   * $ * * $     #  */ at(6, 3), at(7, 2), at(8, 1), at(9, 0), at(10, -1), at(11, -2), at(12, -3), at(13, -4),
				/*   # #           # #   */ at(8, 2), at(9, 1), at(10, 0), at(11, -1), at(12, -2),
				/*      # # # # # #      */
			},
			Init: hexaban.Init{
				Goals: []hexaban.HexCoord{
					at(5, -3), at(6, -3), at(7, -4), at(5, -1), at(6, -2), at(8, -4), at(9, -5), at(5, 0), at(7, -2), at(8, -3), at(10, -5), at(6, 0), at(7, -1), at(9, -3), at(10, -4), at(6, 1), at(8, -1), at(9, -2), at(11, -4), at(12, -5), at(7, 1), at(8, 0), at(10, -2), at(11, -3), at(7, 2), at(9, 0), at(10, -1),
				},
				Crates: []hexaban.HexCoord{
					at(5, -3), at(5, -2), at(6, -3), at(7, -4), at(8, -5), at(5, -1), at(9, -5), at(5, 0), at(7, -2), at(8, -3), at(10, -5), at(5, 1), at(7, -1), at(9, -3), at(11, -5), at(6, 1), at(8, -1), at(9, -2), at(11, -4), at(12, -5), at(7, 1), at(11, -3), at(7, 2), at(8, 1), at(9, 0), at(10, -1), at(11, -2),
				},
				Ichiban: at(8, -2),
			},
		},
	}

	t.Run("convert Erim SEVER levels", func(t *testing.T) {
		got, err := convertSEVER([]byte(text), collection)
		if err != nil {
			t.Errorf("convertSEVER() error = %v", err)
			return
		}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("convertSEVER() = %v, expected %v", got, expected)
		}
	})
}
