package main

import (
	"reflect"
	"testing"

	"github.com/SymbolNotFound/hexaban"
)

func Test_convertEgevad(t *testing.T) {
	var text string = `; Hexobans by Sven Egevad

"svenx001"
    # # # 
   #     # 
  #       # 
   # . . # 
  #   .   # 
 #   $ $ $ # 
  # # # *   # 
     # @   # 
      # # #

"svenx002"
          # # # 
   # # # #     # 
  #     $     # # #
 #   $ $   $ # #   #
  #   # #   #   . . # 
 #   $ # $     . # . # 
#   #     # #   . . # 
 # @   # # # # # # # 
  # # # #
`

	// The details here are abritrary for the purpose of testing, but
	// in recognition of the designer's effort this is accurate.
	collection := hexaban.Collection{
		Source: "http://web.telia.com/~u40915103/welcome.htm",
		Author: "Sven Egevad",
	}

	// for convenience, alias the coordinate constructor
	at := hexaban.NewHexCoord

	expected := []hexaban.Puzzle{
		{
			Identity: "SvenHex/001",
			Name:     "sven x 001",
			Author:   collection.Author,
			Source:   collection.Source,
			Terrain: []hexaban.HexCoord{
				/*		# # #     */
				/*   #     #    */ at(3, -2), at(4, -3),
				/*  #       #   */ at(3, -1), at(4, -2), at(5, -3),
				/*   # . . #    */ at(4, -1), at(5, -2),
				/*  #   .   #   */ at(4, 0), at(5, -1), at(6, -2),
				/* #   $ $ $ #  */ at(4, 1), at(5, 0), at(6, -1), at(7, -2),
				/*  # # # *   # */ at(7, -1), at(8, -2),
				/*     # @   #  */ at(7, 0), at(8, -1),
				/*      # # #   */
			},
			Init: hexaban.Init{
				Goals: []hexaban.HexCoord{
					at(4, -1), at(5, -2), at(5, -1), at(7, -1),
				},
				Crates: []hexaban.HexCoord{
					at(5, 0), at(6, -1), at(7, -2), at(7, -1),
				},
				Ichiban: at(7, 0),
			},
		},
		{
			Identity: "SvenHex/002",
			Name:     "sven x 002",
			Author:   collection.Author,
			Source:   collection.Source,
			Terrain: []hexaban.HexCoord{
				/*           # # #        */
				/*    # # # #     #       */ at(6, -5), at(7, -6),
				/*   #     $     # # #    */ at(3, -1), at(4, -2), at(5, -3), at(6, -4), at(7, -5),
				/*  #   $ $   $ # #   #   */ at(3, 0), at(4, -1), at(5, -2), at(6, -3), at(7, -4), at(10, -7),
				/*   #   # #   #   . . #  */ at(4, 0), at(7, -3), at(9, -5), at(10, -6), at(11, -7),
				/*  #   $ # $     . # . # */ at(4, 1), at(5, 0), at(7, -2), at(8, -3), at(9, -4), at(10, -5), at(12, -7),
				/* #   #     # #   . . #  */ at(4, 2), at(6, 0), at(7, -1), at(10, -4), at(11, -5), at(12, -6),
				/*  # @   # # # # # # #   */ at(5, 2), at(6, 1),
				/*   # # # #              */
			},
			Init: hexaban.Init{
				Goals: []hexaban.HexCoord{
					at(10, -6), at(11, -7), at(10, -5), at(12, -7), at(11, -5), at(12, -6),
				},
				Crates: []hexaban.HexCoord{
					at(5, -3), at(4, -1), at(5, -2), at(7, -4), at(5, 0), at(7, -2),
				},
				Ichiban: at(5, 2),
			},
		},
	}

	t.Run("Convert SvenHex puzzles to axial coordinates", func(t *testing.T) {
		got, err := convertEgevad([]byte(text), collection)
		if err != nil {
			t.Errorf("convertEgevad() error = %v", err)
			return
		}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("convertEgevad() = %v, expected %v", got, expected)
		}
	})
}
