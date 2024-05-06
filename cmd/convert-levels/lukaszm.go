package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/SymbolNotFound/hexaban"
)

func convertLukaszM(text []byte, collection hexaban.Collection) ([]hexaban.Puzzle, error) {
	puzzles := make([]hexaban.Puzzle, 0)
	parser := NewParser(text)
	errors := errorGroup{make([]string, 0)}

	if !parser.NextToken("Author: LukaszM") {
		errors.AddError("expected authorship line")
	}
	parser.NextLine()
	parser.NextLine()

	for !parser.EOF() {
		puzzle := hexaban.Puzzle{}
		puzzle.Source = collection.Source
		puzzle.Author = collection.Author

		id_matcher := regexp.MustCompile(`^; Level (\d+)`)
		puzid := id_matcher.FindSubmatch(parser.filedata[parser.cursor:])
		if puzid == nil {
			errors.AddError(
				fmt.Sprintf("expected a comment-line with level ID, got %s", parser.filedata[parser.cursor:]))
			break
		}
		puzzle.Identity = fmt.Sprintf("LukaszM/%s", string(puzid[1]))
		puzzle.Name = string(puzid[0][2:])
		parser.cursor += uint(len(puzid[0]))
		parser.NextLine()

		grid_parser := NewParser(parser.NextSection())
		if !grid_parser.BytesAvailable(2) {
			errors.AddError("not enough data for a puzzle definition")
			break
		}
		// Difficulty property appears with puzzle grid definition
		difficultyValue := grid_parser.ParseProperty("Difficulty")
		difficulty, err := strconv.Atoi(difficultyValue)
		if err != nil {
			errors.AddError("failed to parse difficulty " + difficultyValue)
		} else {
			puzzle.Difficulty = difficulty
		}

		tiles, err := grid_parser.ParseTextGrid()
		if err != nil {
			errors.AddError(
				fmt.Sprintf("failed to parse puzzle initial conditions: %v", err))
			fmt.Printf("...\n%s", grid_parser.filedata)
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
