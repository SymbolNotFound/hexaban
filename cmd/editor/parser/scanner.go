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
// github:SymbolNotFound/hexaban/cmd/editor/parser/scanner.go

package parser

import (
	"bufio"
	"fmt"
	"io"
)

// Produces a stream of puzzle-specific tokens with their (line, row)
// coordinates.  Detects double-newline for end of file and maintains
// odd/even consistency by use of TOKEN_ERROR.  A newline token is
// produced at the end of each line, but it isn't required for callers
// pulling from this channel due to the column property of the token.
func (parser *parserState) StreamTokens() <-chan token {
	tokenchan := make(chan token)
	go func() {
		defer close(tokenchan)
		parser.cursor.Reset()
		parser.readNextToken = scanLineStart
		for parser.readNextToken != nil {
			tokenchan <- parser.readNextToken(parser)
		}
	}()
	return tokenchan
}

// The parser's private state representation, only used within this package.
type parserState struct {
	scanner       *bufio.Reader
	cursor        cursor
	readNextToken func(*parserState) token
}

// Reads a byte from the parser's input reader and updates cursor appropriately,
// returning the byte found and handling \r and \r\n line terminators.
func (parser *parserState) NextByte() byte {
	glyph, err := parser.scanner.ReadByte()
	parser.cursor.NextByte()
	if err != nil {
		if err == io.EOF {
			return byte(TOKEN_EOF)
		}
		return byte(TOKEN_ERROR_WHAT)
	}

	// Move the cursor depending on the read value.
	if glyph == '\r' {
		// The nested ReadByte() will advance the cursor to the next line.
		subsequent := parser.NextByte()
		if subsequent != '\n' && subsequent != 0 {
			return byte(TOKEN_ERROR_WHAT)
		}
		glyph = subsequent
	}
	if glyph == '\n' {
		parser.cursor.NextLine()
		parser.readNextToken = scanLineStart
	}

	return glyph
}

// Adjusts the cursor position based on the next character read.
// This assumes a context of having recently read a tile glyph,
// so the only acceptable characters for the next byte are
//
//	' ' (expected)
//	EOF (done)
//	'\n' (end of line)
//
// Any tile or non-tile character would be an error.  The '\r'
// will be swallowed by NextByte() so only handling '\n' is
// sufficient for OS-independent line handling.
func (parser *parserState) NextTilePosition() error {
	glyph := parser.NextByte()
	if glyph == ' ' || glyph == 0 {
		return nil
	}
	if glyph == '\n' {
		if parser.cursor.lines&1 == 1 {
			space_or_nl := parser.NextByte()
			if space_or_nl == '\n' {
				parser.readNextToken = nil
				return nil
			}
			if space_or_nl != ' ' {
				return fmt.Errorf("odd line with non-space character %c at initial position", space_or_nl)
			}
		}
		return nil
	}

	return fmt.Errorf("unexpected character %c at line %d, col %d",
		glyph, parser.cursor.lines, parser.cursor.cols)
}

// Location of the reader.  Distinct from the token's position properties
// so that cursor transformations aren't exposed on token positioning.
type cursor struct {
	lines uint
	cols  int
}

// Some basic cursor-movement functionality.

func (c *cursor) Reset()    { c.lines, c.cols = 0, -1 }
func (c *cursor) NextByte() { c.cols += 1 }
func (c *cursor) NextLine() { c.lines, c.cols = c.lines+1, -1 }

// Retrieves the next valid token, expecting one of the tile glyphs
// (see TokenType.IsTile()).  Handles newlines and EOF appropriately, too.
func scanTiles(parser *parserState) token {
	glyph := TokenType(parser.NextByte())
	if glyph == TOKEN_EOF {
		parser.readNextToken = nil
		return Token(TOKEN_EOF, parser.cursor)
	}
	if glyph == TOKEN_NEWLINE {
		parser.readNextToken = scanLineStart
		return Token(TOKEN_NEWLINE, parser.cursor)
	}
	if glyph.IsTile() {
		if parser.cursor.cols&1 != int(parser.cursor.lines)&1 {
			return Token(TOKEN_ERROR_ALIGN, parser.cursor)
		}
		token := Token(glyph, parser.cursor)
		parser.NextTilePosition()
		return token
	}

	parser.readNextToken = nil
	return Token(TOKEN_ERROR_WHAT, parser.cursor)
}

// Retrieves the next valid token, ignoring empty lines while accounting for
// the leading spaces before the first wall (#) token.  No other tile tokens
// will be accepted as a first non-space glyph, thus resulting in TOKEN_ERROR.
func scanLineStart(parser *parserState) token {
	glyph := TokenType(parser.NextByte())
	for glyph != TOKEN_EOF {
		if glyph == TOKEN_NEWLINE {
			// For initial-line parsing, reset the position until tiles are seen.
			if parser.cursor.lines == 1 {
				parser.cursor.Reset()
				glyph = TokenType(parser.NextByte())
				continue
			} else {
				// otherwise, this is a double-newline after some tiles have been read.
				break
			}
		}
		if glyph == ' ' {
			glyph = TokenType(parser.NextByte())
			continue
		}
		if glyph == TOKEN_WALL {
			// Outer wall encountered, begin scanning tiles.
			parser.readNextToken = scanTiles
			if parser.cursor.lines == 0 {
				// Fix the odd/even alignment of the initial line,
				// since the first line should always be correct
				// but subsequent lines will be checked for alignment.
				parser.cursor.lines = uint(parser.cursor.cols) & 1
			}
			if parser.cursor.cols&1 != int(parser.cursor.lines)&1 {
				return Token(TOKEN_ERROR_ALIGN, parser.cursor)
			}

			token := Token(glyph, parser.cursor)
			if err := parser.NextTilePosition(); err != nil {
				return Token(TOKEN_ERROR_WHAT, parser.cursor)
			}
			return token
		}

		// any other glyph would be considered an error.
		parser.readNextToken = nil
		return Token(TOKEN_ERROR_WHAT, parser.cursor)
	}

	// Reader has reached the end of the stream or a double-newline, terminate.
	parser.readNextToken = nil
	return Token(TOKEN_EOF, parser.cursor)
}
