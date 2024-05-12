package main

import (
	"errors"
	"fmt"

	"github.com/SymbolNotFound/hexaban"
)

func convertEgevad(text []byte, collection hexaban.Collection) ([]hexaban.Puzzle, error) {
	puzzles := make([]hexaban.Puzzle, 0)
	parser := NewParser(text)
	if !parser.NextToken("; Hexobans by Sven Egevad\n") {
		return nil, errors.New("first line mismatch for SvenHex, missing attribution comment")
	}
	if !parser.NextLine() {
		return puzzles, errors.New("expected a double-newline after the file header")
	}

	errors := errorGroup{make([]string, 0)}
	for !parser.EOF() {
		// Each puzzle in this collection is just a quoted string (puzzle id) and the grid.
		puzzle := hexaban.Puzzle{}
		puzzle.Author = collection.Author
		puzzle.Source = collection.Source

		puzzle.Identity = parser.NextQuotedString()
		if puzzle.Identity == "" {
			errors.AddError(
				fmt.Sprintf("expected a quoted string for the name of level %d", len(puzzles)))
		} else {
			puzzle.Identity = puzzle.Identity[1 : len(puzzle.Identity)-1]
			if puzzle.Identity[:5] == "svenx" {
				puzzle.Name = puzzle.Identity[5:]
				puzzle.Identity = fmt.Sprintf("SvenHex/%s", puzzle.Name)
				puzzle.Name = "sven x " + puzzle.Name
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
			break
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
