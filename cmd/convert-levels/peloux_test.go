package main

import (
	"reflect"
	"testing"

	"github.com/SymbolNotFound/hexaban"
)

func Test_convertPeloux(t *testing.T) {
	var text string = `Author: Aymeric du Peloux

; HEXOCET 01
       # #
      # @ #
     # . $ #
    # $ #   #
 # # .   . $ # #
#     #   #     #
 #   # #   #   #
#               #
 # # # # #   # #
          # #
Title: A
Difficulty: 6
Created: Jan 2002
Hint:  First, you should put 2 boxes in the center of the maze (one on the right target). Then push down the 3rd on the left down side, before to center it too.

; HEXOCET 02
      # #
   # # @ #
  #       # #
 #   # $ $   #
#     * . #   #
 #   # * #   #
  #     .   #
   # # # # #
Title: Perfume
Difficulty: 5
Created: Jan 2002
Hint: The difficulty is to center the box which will be on top-right. To do it, you should separate the 3 other boxes in a triangular position.
`

	collection := hexaban.Collection{}

	// For convenience, an alias to the coordinate constructor.
	at := hexaban.NewHexCoord

	expected := []hexaban.Puzzle{
		{
			Identity: "hexocet/01",
			Name:     "A",
			Author:   collection.Author,
			Source:   collection.Source,
			Terrain: []hexaban.HexCoord{
				/*        # #        */
				/*       # @ #       */ at(5, -3),
				/*      # . $ #      */ at(5, -2), at(6, -3),
				/*     # $ #   #     */ at(5, -1), at(7, -3),
				/*  # # .   . $ # #  */ at(5, 0), at(6, -1), at(7, -2), at(8, -3),
				/* #     #   #     # */ at(4, 2), at(5, 1), at(7, -1), at(9, -3), at(10, -4),
				/*  #   # #   #   #  */ at(5, 2), at(8, -1), at(10, -3),
				/* #               # */ at(5, 3), at(6, 2), at(7, 1), at(8, 0), at(9, -1), at(10, -2), at(11, -3),
				/*  # # # # #   # #  */ at(10, -1),
				/*           # #     */
			},
			Init: hexaban.Init{
				Goals: []hexaban.HexCoord{
					at(5, -2), at(5, 0), at(7, -2),
				},
				Crates: []hexaban.HexCoord{
					at(6, -3), at(5, -1), at(8, -3),
				},
				Ichiban: at(5, -3),
			},
			Difficulty: 6,
		},
		{
			Identity: "hexocet/02",
			Name:     "Perfume",
			Author:   collection.Author,
			Source:   collection.Source,
			Terrain: []hexaban.HexCoord{
				/*       # #       */
				/*    # # @ #      */ at(4, -3),
				/*   #       # #   */ at(3, -1), at(4, -2), at(5, -3),
				/*  #   # $ $   #  */ at(3, 0), at(5, -2), at(6, -3), at(7, -4),
				/* #     * . #   # */ at(3, 1), at(4, 0), at(5, -1), at(6, -2), at(8, -4),
				/*  #   # * #   #  */ at(4, 1), at(6, -1), at(8, -3),
				/*   #     .   #   */ at(5, 1), at(6, 0), at(7, -1), at(8, -2),
				/*    # # # # #    */
			},
			Init: hexaban.Init{
				Goals: []hexaban.HexCoord{
					at(5, -1), at(6, -2), at(6, -1), at(7, -1),
				},
				Crates: []hexaban.HexCoord{
					at(5, -2), at(6, -3), at(5, -1), at(6, -1),
				},
				Ichiban: at(4, -3),
			},
			Difficulty: 5,
		},
	}

	t.Run("convert hexocet levels", func(t *testing.T) {
		got, err := convertPeloux([]byte(text), collection)
		if err != nil {
			t.Errorf("convertPeloux() error = %v", err)
			return
		}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("convertPeloux() = %v, expected %v", got, expected)
		}
	})
}
