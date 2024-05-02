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
// github:SymbolNotFound/hexaban/cmd/convert-levels/parser.go

package main

import (
	"fmt"
	"regexp"

	"github.com/SymbolNotFound/hexaban"
)

type PuzzleParser struct {
	filedata []byte
	cursor   uint
}

// Constructor function for PuzzleParser
func NewParser(text []byte) PuzzleParser {
	return PuzzleParser{text, 0}
}

// Has the parser reached End of File?
func (parser *PuzzleParser) EOF() bool {
	return parser.cursor >= uint(len(parser.filedata))
}

func (parser *PuzzleParser) EOL() bool {
	return parser.BytesAvailable(1) && parser.filedata[parser.cursor] == '\n'
}

// Does the file have at least this many bytes available for parsing?
func (parser *PuzzleParser) BytesAvailable(count int) bool {
	return parser.cursor+uint(count) <= uint(len(parser.filedata))
}

// Returns true if the next text after the cursor matches
func (parser *PuzzleParser) NextToken(image string) bool {
	if !parser.BytesAvailable(len(image)) {
		return false
	}

	limit := parser.cursor + uint(len(image))
	if string(parser.filedata[parser.cursor:limit]) != image {
		return false
	}

	parser.cursor += uint(len(image))
	return true
}

// Returns true if the next token is a newline \n, possibly \r\n,
// with any amount of leading spaces or tabs.  Returns false at EOF.
func (parser *PuzzleParser) NextLine() bool {
	if !parser.BytesAvailable(1) {
		return false
	}
	matcher := regexp.MustCompile(`^[\s\t]*\r?\n`)
	match := matcher.Find(parser.filedata[parser.cursor:])
	if match == nil {
		return false
	}

	parser.cursor += uint(len(match))
	return true
}

// Returns a quoted string (including the wrapping double-quotes) if it is the
// next token, or returns the empty string if matching fails.
func (parser *PuzzleParser) NextQuotedString() string {
	if !parser.BytesAvailable(2) {
		return ""
	}

	matcher := regexp.MustCompile(`^"(\["bfnrt/\]|\\u[a-fA-F0-9]{4}|[^"\n])*"`)
	match := matcher.Find(parser.filedata[parser.cursor:])
	if match == nil {
		return ""
	}

	parser.cursor += uint(len(match))
	return string(match)
}

// Reads the next section of filedata.
//
// Returns the contents of filedata until the next "\n\n" as []byte (suitable for
// obtaining the entire (or remaining) contents of a puzzle's definition).  Also
// advances the current parser's cursor past the double newline.
//
// If no double-newline is found, the remainder of filedata is returned.
func (parser *PuzzleParser) NextSection() []byte {
	if parser.BytesAvailable(2) {
		matcher := regexp.MustCompile(`^([^\n]|\n[^\n])*\n\n`)
		match := matcher.Find(parser.filedata[parser.cursor:])
		if match == nil {
			// return remainder of file if no \n\n is found.
			match = parser.filedata[parser.cursor:]
		}
		parser.cursor += uint(len(match))
		return match
	}

	// Only one or maybe zero bytes remain, return nothing.
	parser.cursor = uint(len(parser.filedata))
	return []byte{}
}

// The puzzles in text format are rectangular representations that approximately
// look like their hexagonal counterpart -- every other character is used, and
// rows are staggered across two text lines to represent their vertically offset
// positions in hexagonal grids.
//
// A tricky edge case is that some grids start with an even number of spaces and
// some grids start with an odd number of spaces, and this affects the coordinate
// system and properly converting it into the (i, j) unit vector representation.
// To handle this, we offset initial lines if they're odd, setting a cleared zero
// row to make even-indexed lines always even-formatted rows.  Later, all coord-
// inates are normalized to be centered around the player's initial position, so
// there is no confusion over alignment even for puzzles otherwise equivalent.
//
// The EBNF representation of this parser (approximately):
//
// puzzle_grid ::= [odd_line] (even_line odd_line)* even_line [odd_line] '\n'
//
// odd_line ::= ' ' tile+ '\n'
// even_line ::= tile+ '\n'
//
// tile ::= glyph sep
// glyph ::= '#' | ' ' | '.' | '$' | '*' | '@' | '+'
// sep ::= ' ' | '\n' | '\r' '\n' | <EOF>
func (parser *PuzzleParser) ParseTextGrid() ([]Tile, error) {
	tiles := make([]Tile, 0)
	row, column := 0, 0
	tile_matcher := regexp.MustCompile(`^([# .$*@+])( |\r?\n)?`)

	// Read until we reach end of file or until a double-newline
	for !parser.EOF() {
		for parser.NextToken(" ") {
			column += 1
		}
		if column%2 == 1 {
			if row == 0 {
				row = 1
			} else {
				if row%2 == 0 {
					return nil, fmt.Errorf("misalignment: odd column %d and even row %d", column, row)
				}
			}
		}

		for {
			match := tile_matcher.FindSubmatch(parser.filedata[parser.cursor:])
			if match == nil {
				if parser.EOF() || parser.EOL() {
					break
				}
				return tiles, fmt.Errorf("unrecognized puzzle tile %s",
					parser.filedata[parser.cursor:parser.cursor+1])
			}

			glyph, sep := match[1], ""
			if match[2] != nil {
				sep = string(match[2])
			}

			coord := hexaban.NewRectCoord(uint(column), uint(row>>1)).ToHex()
			tiles = append(tiles, parseTile(glyph, coord)...)
			parser.cursor += uint(len(match[0]))

			// separator indicates whether line continues or not
			if sep == " " {
				column += 2
				continue
			}
			if sep == "\n" || sep == "\r\n" || (sep == "" && parser.EOF()) {
				break
			}
		}

		for parser.NextLine() {
			// Keep advancing if there are multiple newlines (typically at the end).
		}
		row += 1
		column = 0
	}

	return tiles, nil
}

type (
	Tile   = hexaban.Tile
	Floor  = hexaban.Floor
	Wall   = hexaban.Wall
	Goal   = hexaban.Goal
	Crate  = hexaban.Crate
	Player = hexaban.Player

	HexCoord = hexaban.HexCoord
)

func parseTile(glyph []byte, coord HexCoord) []Tile {
	switch glyph[0] {
	case '#':
		return []Tile{Wall{HexCoord: coord}}
	case ' ':
		return []Tile{Floor{HexCoord: coord}}
	case '.':
		// goal without a crate on it
		return []Tile{
			Floor{HexCoord: coord},
			Goal{HexCoord: coord},
		}
	case '$':
		// crate on normal (non-goal) floor
		return []Tile{
			Floor{HexCoord: coord},
			Crate{HexCoord: coord},
		}
	case '*':
		// crate on goal
		return []Tile{
			Floor{HexCoord: coord},
			Goal{HexCoord: coord},
			Crate{HexCoord: coord},
		}
	case '@':
		// player on normal (non-goal) floor
		return []Tile{
			Floor{HexCoord: coord},
			Player{HexCoord: coord},
		}
	case '+':
		// player on goal
		return []Tile{
			Floor{HexCoord: coord},
			Player{HexCoord: coord},
			Goal{HexCoord: coord},
		}
	}
	return []Tile{}
}
