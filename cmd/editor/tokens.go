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
// github:SymbolNotFound/hexoban/cmd/editor/tokens.go

package main

import "github.com/SymbolNotFound/hexoban/puzzle"

// A token is defined by its glyph (printable) representation and its position.
// The position is relative to the first non-empty line of the text, so it may
// differ from the absolute byte position as measured from the beginning.
type token struct {
	glyph TokenType
	line  uint
	col   uint
}

// Constructor function for a token instance.
func Token(glyph TokenType, cursor cursor) token {
	return token{glyph, cursor.lines, uint(cursor.cols)}
}

type TokenType byte

const (
	TOKEN_EOF         TokenType = 0
	TOKEN_UNKNOWN     TokenType = '?'
	TOKEN_ERROR_WHAT  TokenType = '!'
	TOKEN_ERROR_ALIGN TokenType = 'x'
	TOKEN_NEWLINE     TokenType = '\n'
	TOKEN_WALL        TokenType = '#'
	TOKEN_FLOOR       TokenType = ' '
	TOKEN_GOAL        TokenType = '.'
	TOKEN_CRATE       TokenType = '$'
	TOKEN_CRATE_GOAL  TokenType = '*'
	TOKEN_AT          TokenType = '@'
	TOKEN_AT_GOAL     TokenType = '+'
)

// Returns true if the glyph is a valid hexoban tile representation.
func (glyph TokenType) IsTile() bool {
	switch glyph {
	case TOKEN_WALL, TOKEN_FLOOR, TOKEN_GOAL,
		TOKEN_CRATE, TOKEN_CRATE_GOAL,
		TOKEN_AT, TOKEN_AT_GOAL:
		return true
	}
	return false
}

type (
	Tile   = puzzle.Tile
	Floor  = puzzle.Floor
	Wall   = puzzle.Wall
	Goal   = puzzle.Goal
	Crate  = puzzle.Crate
	Player = puzzle.Player

	HexCoord = puzzle.HexCoord
)

func (t token) parseTile() []Tile {
	coord := RectCoord{t.col, t.line}.ToHex()

	switch t.glyph {
	case TOKEN_WALL:
		// Walls are provided in the results but they can later be ignored.
		return []Tile{Wall{HexCoord: coord}}
	case TOKEN_FLOOR:
		// This assumes parseTile is only called after the first (wall) glyph.
		return []Tile{Floor{HexCoord: coord}}
	case TOKEN_GOAL:
		// A goal without a crate on it.
		return []Tile{
			Floor{HexCoord: coord},
			Goal{HexCoord: coord},
		}
	case TOKEN_CRATE:
		// A crate on normal (non-goal) floor.
		return []Tile{
			Floor{HexCoord: coord},
			Crate{HexCoord: coord},
		}
	case TOKEN_CRATE_GOAL:
		// A crate on a goal.
		return []Tile{
			Floor{HexCoord: coord},
			Goal{HexCoord: coord},
			Crate{HexCoord: coord},
		}
	case TOKEN_AT:
		// A player on a normal (non-goal) floor.
		return []Tile{
			Floor{HexCoord: coord},
			Player{HexCoord: coord},
		}
	case TOKEN_AT_GOAL:
		// A player on a goal.
		return []Tile{
			Floor{HexCoord: coord},
			Player{HexCoord: coord},
			Goal{HexCoord: coord},
		}
	}

	// Only unrecognized glyphs result in a nil response.
	return nil
}
