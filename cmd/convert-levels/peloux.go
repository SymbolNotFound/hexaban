package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/SymbolNotFound/hexaban"
)

func convertPeloux(text []byte, collection hexaban.Collection) ([]hexaban.Puzzle, error) {
	puzzles := make([]hexaban.Puzzle, 0)
	parser := NewParser(text)
	errors := errorGroup{make([]string, 0)}

	author := parser.ParseProperty("Author")
	if author == "" {
		errors.AddError("expected Author property, none found")
	} else if author != "Aymeric du Peloux" {
		errors.AddError("expected correct author (" + collection.Author + ") but found " + author)
	}
	parser.NextLine()

	for !parser.EOF() {
		puzzle := hexaban.Puzzle{}
		puzzle.Source = collection.Source
		puzzle.Author = collection.Author

		id_matcher := regexp.MustCompile(`^; HEXOCET (\d+)`)
		puzid := id_matcher.FindSubmatch(parser.filedata[parser.cursor:])
		if puzid == nil {
			errors.AddError(
				fmt.Sprintf("expected a comment-line with level ID, got %s", parser.filedata[parser.cursor:]))
			break
		}
		puzzle.Identity = fmt.Sprintf("hexocet/%s", string(puzid[1]))
		parser.cursor += uint(len(puzid[0]))
		parser.NextLine()

		grid_parser := NewParser(parser.NextSection())
		if !grid_parser.BytesAvailable(2) {
			errors.AddError("not enough data for a puzzle definition")
			break
		}
		puzzle.Name = grid_parser.ParseProperty("Title")
		_ = grid_parser.ParseProperty("Created")
		_ = grid_parser.ParseProperty("Hint")
		_ = grid_parser.ParseProperty("Comment")
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
