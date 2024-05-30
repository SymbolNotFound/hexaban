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
// github:SymbolNotFound/hexoban/cmd/editor/parser/parser.go

package parser

import (
	"bufio"
	"fmt"

	hexoban "github.com/SymbolNotFound/hexoban"
)

// Parses a line-oriented sequence of (tile, gap) byte pairs representing the
// double-height offset coordinates layout of a hexagonal grid for a hexoban
// puzzle's initial conditions.  See TokenType for the glyph representations.
//
// The grammar has no recursion in its production rules so a hand-rolled parser
// is a slightly better representation than producing the parser from a grammar
// definition, as there is little need for nesting and a complex initial token.
func ParsePuzzleDefinition(reader *bufio.Reader, puzzle *hexoban.Puzzle) error {
	parser := parserState{scanner: reader}
	tiles := make([]hexoban.Tile, 0)

	for token := range parser.StreamTokens() {
		switch token.glyph {
		case TOKEN_ERROR_WHAT:
			return fmt.Errorf("invalid/unexpected token at line %d, col %d", token.line, token.col)
		case TOKEN_ERROR_ALIGN:
			return fmt.Errorf("odd/even mismatch at line %d, col %d", token.line, token.col)
		case TOKEN_EOF:
			puzzle.AddTiles(tiles)
			return nil

		case TOKEN_WALL, TOKEN_FLOOR, TOKEN_GOAL,
			TOKEN_CRATE, TOKEN_CRATE_GOAL,
			TOKEN_AT, TOKEN_AT_GOAL:
			// Don't add the tiles to the puzzle immediately,
			// an error may still be found during parsing.
			tiles = append(tiles, token.parseTile()...)
		}
	}

	// Add all the tiles together before returning.
	puzzle.AddTiles(tiles)
	return nil
}
