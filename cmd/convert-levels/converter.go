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
// github:SymbolNotFound/hexaban/cmd/convert-levels/hexparse.go

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/SymbolNotFound/hexaban"
)

type ParseFn func([]byte, hexaban.Collection) ([]hexaban.Puzzle, error)

// Defines metadata for file being parsed to produce independent puzzle files.
type CollectionFactory struct {
	Author    string
	Source    string
	InputPath string
	ParseFn   ParseFn
}

// Composite of multiple errors, for parsing as many puzzles as possible.
type errorGroup struct {
	errors []string
}

func (errs errorGroup) Error() string           { return strings.Join(errs.errors, "\n") }
func (errs *errorGroup) AddError(errstr string) { errs.errors = append(errs.errors, errstr) }

func main() {
	for _, factory := range FileMetadata() {
		fmt.Printf("Parsing '%s'...\n", factory.InputPath)
		filedata, err := os.ReadFile(factory.InputPath)
		if err != nil {
			fmt.Println(err)
			continue
		}

		collection := hexaban.Collection{
			Puzzles: make([]hexaban.Puzzle, 0),
			Author:  factory.Author,
			Source:  factory.Source,
		}
		collection.Puzzles, err = factory.ParseFn(filedata, collection)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, hex_puzzle := range collection.Puzzles {
			hex_puzzle = translateCoordinates(hex_puzzle)
			json, err := prettyPrint(hex_puzzle)
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = os.WriteFile(
				fmt.Sprintf("%s.json", path.Join("levels", hex_puzzle.Identity)),
				json, 0644)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}

func FileMetadata() []CollectionFactory {
	return []CollectionFactory{
		{
			"David W. Skinner",
			"http://users.bentonrea.com/~sasquatch/sokoban/hex.html",
			"data/dwshex.hsb",
			convertDWS,
		},
		{
			"François Marques",
			"http://hexoban.online.fr/",
			"data/heloban.hsb",
			convertMarques,
		},
		{
			"François Marques",
			"http://hexoban.online.fr/",
			"data/heroban.hsb",
			convertMarques,
		},
		{
			"Erim SEVER",
			"www.erimsever.com/sokoban/Erim_Levels/E_Hexoban.zip",
			"data/all_E_Hex.hsb",
			convertSEVER,
		},
		{
			"Aymeric du Peloux",
			"http://membres.lycos.fr/nabokos/",
			"data/hexocet.hsb",
			convertPeloux,
		},
		{
			"LukaszM",
			"https://play.fancade.com/5FA6BCFD16EB8B3B",
			"data/lukaszm.hsb",
			convertLukaszM,
		},
		{
			"", // Mixture of authors.
			"http://users.bentonrea.com/~sasquatch/sokoban/morehex.hsb",
			"data/morehex.hsb",
			convertSingles,
		},
	}
}

func translateCoordinates(topleft hexaban.Puzzle) hexaban.Puzzle {
	centered := hexaban.Puzzle{
		Identity:   topleft.Identity,
		Name:       topleft.Name,
		Author:     topleft.Author,
		Source:     topleft.Source,
		Difficulty: topleft.Difficulty,
		Init:       hexaban.Init{},
	}
	player := topleft.Init.Ichiban

	centered.Terrain = make([]hexaban.HexCoord, len(topleft.Terrain))
	for index, terrain := range topleft.Terrain {
		centered.Terrain[index] = hexaban.NewHexCoord(terrain.I()-player.I(), terrain.J()-player.J())
	}

	centered.Init.Goals = make([]hexaban.HexCoord, len(topleft.Init.Goals))
	for index, goal := range topleft.Init.Goals {
		centered.Init.Goals[index] = hexaban.NewHexCoord(goal.I()-player.I(), goal.J()-player.J())
	}

	centered.Init.Crates = make([]hexaban.HexCoord, len(topleft.Init.Crates))
	for index, crate := range topleft.Init.Crates {
		centered.Init.Crates[index] = hexaban.NewHexCoord(crate.I()-player.I(), crate.J()-player.J())
	}

	return centered
}

// Special converter function for the `morehex.hsb` file.  Could be used for other
// solo contributions if any appear before the editor is ready.  Assumes the puzzle
// grid is followed by an "Author" property.
func convertSingles(text []byte, collection hexaban.Collection) ([]hexaban.Puzzle, error) {
	puzzles := make([]hexaban.Puzzle, 0)
	parser := NewParser(text)

	errors := errorGroup{make([]string, 0)}
	for !parser.EOF() {
		puzzle := hexaban.Puzzle{}
		puzzle.Source = collection.Source

		// puzzle ID, also represents puzzle's Name
		if !parser.NextToken("; ") {
			errors.AddError("expected comment for puzzle identifier")
		}
		puzzle.Identity = parser.NextQuotedString()
		if puzzle.Identity == "" {
			errors.AddError(
				fmt.Sprintf("expected a quoted string for the name of level %d", len(puzzles)))
		} else {
			// Valid identity found, strip quotes from Identity and also assign to Name
			puzzle.Identity = removeFirstAndLast(puzzle.Identity)
			puzzle.Name = puzzle.Identity
			puzzle.Identity = "more/" + puzzle.Identity
		}
		if !parser.NextLine() {
			errors.AddError("expected newline after puzzle id")
			parser.NextSection()
			continue
		}

		// puzzle grid
		grid_parser := NewParser(parser.NextSection())
		if !grid_parser.BytesAvailable(2) {
			errors.AddError("not enough data for a puzzle definition")
			break
		}
		// Author credit appears after grid section, pull it before parsing the grid
		puzzle.Author = grid_parser.ParseProperty("Author")
		puzzle.Name = grid_parser.ParseProperty("Title")
		tiles, err := grid_parser.ParseTextGrid()
		if err != nil {
			errors.AddError(
				fmt.Sprintf("failed to parse puzzle initial conditions: %v", err))
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

// Takes a quoted or parenthetical or bracketed character sequence and returns the
// same characters but with the first and last characters removed.
// Warning: does not check that the enclosing characters are at the outermost
// positions, nor that they match, or any such validation.
func removeFirstAndLast(quoted string) string {
	return quoted[1 : len(quoted)-1]
}

// Pretty-print the puzzle definition as human-friendly (multiline) JSON.
func prettyPrint(puzzle hexaban.Puzzle) ([]byte, error) {
	json, err := json.MarshalIndent(puzzle, "", "  ")
	if err != nil {
		return nil, err
	}

	// Some annoying newline-per-comma happens with MarshalIndent() -- too much indent!
	// And since json.MarshalIndent(_, prefix, separator) is all-on or all-off (even if
	// you explicitly flatten the []byte in a struct's custom marshaler).  We can fix it
	// back up a little with some regex-replace-all... some light regex-fu, if you will.
	// Note that the only thing appearing outside of (...) groups are newlines, spaces
	// and square brackets, and all are always replaced (modulo spacing) in the output.
	json = regexp.MustCompile(`\[\n\s+(-?\d+)`).ReplaceAll(json, []byte("[$1"))
	json = regexp.MustCompile(`\[(-?\d+),\n\s+`).ReplaceAll(json, []byte("[$1, "))
	json = regexp.MustCompile(`(-?\d+)\n\s+\]`).ReplaceAll(json, []byte("$1]"))
	json = regexp.MustCompile(`(-?\d+\],?)\n\s+`).ReplaceAll(json, []byte("$1 "))
	json = regexp.MustCompile(`\] (\],?)\n(\s+)"`).ReplaceAll(json, []byte("]\n$2$1\n$2\""))
	json = regexp.MustCompile(`\] (\],?)\n(\s+)}`).ReplaceAll(json, []byte("]\n$2  $1\n$2}"))

	return json, nil
}
