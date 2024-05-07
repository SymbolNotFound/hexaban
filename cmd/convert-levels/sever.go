package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/SymbolNotFound/hexaban"
)

func convertSEVER(text []byte, collection hexaban.Collection) ([]hexaban.Puzzle, error) {
	puzzles := make([]hexaban.Puzzle, 0)
	parser := NewParser(text)
	errors := errorGroup{make([]string, 0)}

	for !parser.EOF() {
		puzzle := hexaban.Puzzle{}
		puzzle.Source = collection.Source
		puzzle.Author = collection.Author

		id_matcher := regexp.MustCompile(`^; erim_hex(_nr)?(\d+)`)
		match := id_matcher.FindSubmatch(parser.filedata[parser.cursor:])
		if match == nil {
			errors.AddError(
				fmt.Sprintf("expected a comment-line with level ID, got %s", parser.filedata[parser.cursor:]))
			break
		}
		parser.cursor += uint(len(match[0]))
		puzNum, err := strconv.Atoi(string(match[2]))
		if err != nil {
			errors.AddError(err.Error())
			continue
		}
		if string(match[1]) != "" {
			puzNum += 49
		}
		puzzle.Identity = "ErimSEVER/" + fmt.Sprintf("%02d", puzNum)
		puzzle.Name = fmt.Sprintf("Hex %d", puzNum)

		parser.NextLine()

		grid_parser := NewParser(parser.NextSection())
		if !grid_parser.BytesAvailable(2) {
			errors.AddError("not enough data for a puzzle definition")
			break
		}
		author := grid_parser.ParseProperty("Author")
		if author == "" {
			errors.AddError("expected to find an Author property, none found")
			continue
		} else if author != collection.Author {
			errors.AddError("unexpected author value " + author)
			continue
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
