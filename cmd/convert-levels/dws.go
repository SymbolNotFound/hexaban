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
// github:SymbolNotFound/hexaban/cmd/convert-levels/dws.go

package main

import (
	"errors"
	"fmt"

	"github.com/SymbolNotFound/hexaban"
)

// A file converter for the DWS puzzle set by David W. Skinner.
func convertDWS(text []byte, collection hexaban.Collection) ([]hexaban.Puzzle, error) {
	puzzles := make([]hexaban.Puzzle, 0)
	parser := NewParser(text)
	if !parser.NextToken("; Hexobans by David W. Skinner\n") {
		return nil, errors.New("first line mismatch for DWS, missing attribution comment")
	}
	if !parser.NextLine() {
		return puzzles, errors.New("expected a double-newline after the file header")
	}

	errors := errorGroup{make([]string, 0)}
	for !parser.EOF() {
		puzzle := hexaban.Puzzle{}
		puzzle.Author = "David W. Skinner"
		puzzle.Source = collection.Source

		puzzle.Identity = parser.NextQuotedString()
		if puzzle.Identity == "" {
			errors.AddError(
				fmt.Sprintf("expected a quoted string for the name of level %d", len(puzzles)))
		} else {
			puzzle.Identity = puzzle.Identity[1 : len(puzzle.Identity)-1]
			if puzzle.Identity[:3] == "dws" {
				puzzle.Name = puzzle.Identity[3:]
			} else {
				errors.AddError(
					fmt.Sprintf("unexpected identity %s", puzzle.Identity))
			}
		}
		if !parser.NextLine() {
			errors.AddError("expected newline after puzzle id")
		}

		grid_parser := NewParser(parser.NextSection())
		if !grid_parser.BytesAvailable(2) {
			errors.AddError("not enough data for a puzzle definition")
			continue
		}
		tiles, err := grid_parser.ParseTextGrid()
		if err != nil {
			errors.AddError(
				fmt.Sprintf("failed to parse puzzle initial conditions: %v", err))
			break
		}
		puzzle.AddTiles(tiles)
		puzzles = append(puzzles, puzzle)
	}

	if len(errors.errors) == 0 {
		return puzzles, nil
	}
	return puzzles, errors
}
