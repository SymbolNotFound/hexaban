package main

import (
	"reflect"
	"testing"

	"github.com/SymbolNotFound/hexaban"
)

func Test_convertMarques(t *testing.T) {
	var text string = `Author: François Marques
Collection: heloban

     # #
    #   # # #
 # #   $ $   #
#     . .   #
 # $ . @ .   #
  # $ . . $   #
 #       $ # #
  # # #   #
       # #
Title: The circular saw
Date: 2002/02/07
Difficulty: 3

     # #
  # #   # # #
 #     $ $   #
#   # . . #   #
 # $ . @ .   #
  # $ . . $   #
 #     # $ # #
  # #     #
     # # #
Title: The coil toothed
Date: 2002/02/07

     # # #
  # # +   #
 #       #
  # *   #
 # * $ #
#       #
 #     #
  # # #
Title: Heloban 3
Date: 2002/02/10
Difficulty: 5
`
	// The details here are abritrary for the purpose of testing, but
	// in recognition of the designer's effort this is accurate.
	collection := hexaban.Collection{
		Source: "http://hexoban.online.fr/",
		Author: "François Marques",
	}

	// for convenience, alias the coordinate constructor
	at := hexaban.NewHexCoord

	expected := []hexaban.Puzzle{
		{
			Identity: "heloban/The circular saw",
			Name:     "The circular saw",
			Author:   collection.Author,
			Source:   collection.Source,
			Terrain: []hexaban.HexCoord{
				/*      # #        */
				/*     #   # # #   */ at(4, -2),
				/*  # #   $ $   #  */ at(4, -1), at(5, -2), at(6, -3), at(7, -4),
				/* #     . .   #   */ at(3, 1), at(4, 0), at(5, -1), at(6, -2), at(7, -3),
				/*  # $ . @ .   #  */ at(4, 1), at(5, 0), at(6, -1), at(7, -2), at(8, -3),
				/*   # $ . . $   # */ at(5, 1), at(6, 0), at(7, -1), at(8, -2), at(9, -3),
				/*  #       $ # #  */ at(5, 2), at(6, 1), at(7, 0), at(8, -1),
				/*   # # #   #     */ at(8, 0),
				/*        # #      */
			},
			Init: hexaban.Init{
				Goals: []hexaban.HexCoord{
					at(5, -1), at(6, -2), at(5, 0), at(7, -2), at(6, 0), at(7, -1),
				},
				Crates: []hexaban.HexCoord{
					at(5, -2), at(6, -3), at(4, 1), at(5, 1), at(8, -2), at(8, -1),
				},
				Player: at(6, -1),
			},
			Difficulty: 3,
		},
		{
			Identity: "heloban/The coil toothed",
			Name:     "The coil toothed",
			Author:   collection.Author,
			Source:   collection.Source,
			Terrain: []hexaban.HexCoord{
				/*      # #        */
				/*   # #   # # #   */ at(4, -2),
				/*  #     $ $   #  */ at(3, 0), at(4, -1), at(5, -2), at(6, -3), at(7, -4),
				/* #   # . . #   # */ at(3, 1), at(5, -1), at(6, -2), at(8, -4),
				/*  # $ . @ .   #  */ at(4, 1), at(5, 0), at(6, -1), at(7, -2), at(8, -3),
				/*   # $ . . $   # */ at(5, 1), at(6, 0), at(7, -1), at(8, -2), at(9, -3),
				/*  #     # $ # #  */ at(5, 2), at(6, 1), at(8, -1),
				/*   # #     #     */ at(7, 1), at(8, 0),
				/*      # # #      */
			},
			Init: hexaban.Init{
				Goals: []hexaban.HexCoord{
					at(5, -1), at(6, -2), at(5, 0), at(7, -2), at(6, 0), at(7, -1),
				},
				Crates: []hexaban.HexCoord{
					at(5, -2), at(6, -3), at(4, 1), at(5, 1), at(8, -2), at(8, -1),
				},
				Player: at(6, -1),
			},
		},
		{
			Identity: "heloban/Heloban 3",
			Name:     "Heloban 3",
			Author:   collection.Author,
			Source:   collection.Source,
			Terrain: []hexaban.HexCoord{
				/*      # # #  */
				/*   # # +   # */ at(4, -2), at(5, -3),
				/*  #       #  */ at(3, 0), at(4, -1), at(5, -2),
				/*   # *   #   */ at(4, 0), at(5, -1),
				/*  # * $ #    */ at(4, 1), at(5, 0),
				/* #       #   */ at(4, 2), at(5, 1), at(6, 0),
				/*  #     #    */ at(5, 2), at(6, 1),
				/*   # # #     */
			},
			Init: hexaban.Init{
				Goals: []hexaban.HexCoord{
					at(4, -2), at(4, 0), at(4, 1),
				},
				Crates: []hexaban.HexCoord{
					at(4, 0), at(4, 1), at(5, 0),
				},
				Player: at(4, -2),
			},
			Difficulty: 5,
		},
	}

	t.Run("Convert Heloban/Heroban levels", func(t *testing.T) {
		got, err := convertMarques([]byte(text), collection)
		if err != nil {
			t.Errorf("convertMarques() error = %v", err)
			return
		}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("convertMarques() = %v, expected %v", got, expected)
		}
	})
}
