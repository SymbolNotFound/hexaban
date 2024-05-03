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
	Author     string
	Source     string
	InputPath  string
	OutputPath string
	ParseFn    ParseFn
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
			json, err := json.MarshalIndent(hex_puzzle, "", "  ")
			if err != nil {
				fmt.Println(err)
				continue
			}
			json = regexp.MustCompile(`\[\n\s+(-?\d+)`).ReplaceAll(json, []byte("[$1"))
			json = regexp.MustCompile(`\[(-?\d+),\n\s+`).ReplaceAll(json, []byte("[$1, "))
			json = regexp.MustCompile(`(-?\d+)\n\s+\]`).ReplaceAll(json, []byte("$1]"))
			json = regexp.MustCompile(`(-?\d+\],?)\n\s+`).ReplaceAll(json, []byte("$1 "))
			json = regexp.MustCompile(`\] ([\]}])`).ReplaceAll(json, []byte("]\n  $1"))

			err = os.WriteFile(
				fmt.Sprintf("%s.json", path.Join(factory.OutputPath, hex_puzzle.Identity)),
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
			"levels/DWS/",
			convertDWS,
		},
		{
			"François Marques",
			"http://hexoban.online.fr/",
			"data/heloban.hsb",
			"levels/heloban/",
			convertMarques,
		},
		{
			"François Marques",
			"http://hexoban.online.fr/",
			"data/heroban.hsb",
			"levels/heroban/",
			convertMarques,
		},
		{
			"Erim SEVER",
			"www.erimsever.com/sokoban/Erim_Levels/E_Hexoban.zip",
			"data/all_E_Hex.hsb",
			"levels/ErimSEVER/",
			convertSEVER,
		},
		{
			"Aymeric du Peloux",
			"http://membres.lycos.fr/nabokos/",
			"data/hexocet.hsb",
			"levels/hexocet/",
			convertPeloux,
		},
		{
			"LukaszM",
			"https://play.fancade.com/5FA6BCFD16EB8B3B",
			"data/lukaszm.hsb",
			"levels/LukaszM/",
			convertLukaszM,
		},
		{
			"", // Mixture of authors.
			"http://users.bentonrea.com/~sasquatch/sokoban/morehex.hsb",
			"data/morehex.hsb",
			"levels/more/",
			convertSingles,
		},
	}
}

func convertSingles(text []byte, collection hexaban.Collection) ([]hexaban.Puzzle, error) {
	puzzles := make([]hexaban.Puzzle, 0)
	parser := NewParser(text)

	errors := errorGroup{make([]string, 0)}
	for !parser.EOF() {
		puzzle := hexaban.Puzzle{}
		puzzle.Source = collection.Source

		if !parser.NextToken("; ") {
			errors.AddError("expected comment for puzzle identifier")
		}
		puzzle.Identity = parser.NextQuotedString()
		if puzzle.Identity == "" {
			errors.AddError(
				fmt.Sprintf("expected a quoted string for the name of level %d", len(puzzles)))
		} else {
			// strip quotes from Identity and also assign to Name
			puzzle.Identity = withoutQuotes(puzzle.Identity)
			puzzle.Name = puzzle.Identity
		}
		if !parser.NextLine() {
			errors.AddError("expected newline after puzzle id")
			parser.NextSection()
			continue
		}

		grid_parser := NewParser(parser.NextSection())
		if !grid_parser.BytesAvailable(2) {
			errors.AddError("not enough data for a puzzle definition")
			break
		}
		puzzle.Author = grid_parser.ParseProperty("Author")
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

func withoutQuotes(quoted string) string {
	return quoted[1 : len(quoted)-1]
}
