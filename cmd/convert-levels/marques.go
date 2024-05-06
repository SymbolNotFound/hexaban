package main

import (
	"strconv"

	"github.com/SymbolNotFound/hexaban"
)

func convertMarques(text []byte, collection hexaban.Collection) ([]hexaban.Puzzle, error) {
	puzzles := make([]hexaban.Puzzle, 0)
	parser := NewParser(text)
	errors := errorGroup{make([]string, 0)}

	author := parser.ParseProperty("Author")
	if author == "" {
		errors.AddError("expected Author property, none found")
	} else if author != collection.Author {
		errors.AddError("expected correct author (" + collection.Author + ") but found " + author)
	}
	collectionName := parser.ParseProperty("Collection")
	if collectionName == "" {
		errors.AddError("expected Collection property, none found")
	}
	parser.NextLine()

	for !parser.EOF() {
		puzzle := hexaban.Puzzle{}
		puzzle.Source = collection.Source
		puzzle.Author = collection.Author

		grid_parser := NewParser(parser.NextSection())
		if !grid_parser.BytesAvailable(2) {
			errors.AddError("not enough data for a puzzle definition")
			break
		}
		title := grid_parser.ParseProperty("Title")
		puzzle.Name = title
		puzzle.Identity = collectionName + "/" + title
		_ = grid_parser.ParseProperty("Date")
		difficultyValue := grid_parser.ParseProperty("Difficulty")
		difficulty, err := strconv.Atoi(difficultyValue)
		if err == nil {
			puzzle.Difficulty = difficulty
		}

		tiles, err := grid_parser.ParseTextGrid()
		if err != nil {
			errors.AddError(err.Error())
			continue
		}
		puzzle.AddTiles(tiles)
		puzzles = append(puzzles, puzzle)
	}

	if len(errors.errors) == 0 {
		return puzzles, nil
	}
	return puzzles, errors
}
