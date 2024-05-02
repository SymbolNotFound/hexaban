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
			json, err := json.Marshal(hex_puzzle)
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = os.WriteFile(path.Join(factory.OutputPath, hex_puzzle.Identity), json, 0644)
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
			parseNYI,
		},
		{
			"François Marques",
			"http://hexoban.online.fr/",
			"data/heroban.hsb",
			"levels/heroban/",
			parseNYI,
		},
		{
			"Erim SEVER",
			"www.erimsever.com/sokoban/Erim_Levels/E_Hexoban.zip",
			"data/all_E_Hex.hsb",
			"levels/ErimSEVER/",
			parseNYI,
		},
		{
			"Aymeric du Peloux",
			"http://membres.lycos.fr/nabokos/",
			"data/hexocet.hsb",
			"levels/hexocet/",
			parseNYI,
		},
		{
			"LukaszM",
			"https://play.fancade.com/5FA6BCFD16EB8B3B",
			"data/lukaszm.hsb",
			"levels/LukaszM/",
			parseNYI,
		},
		{
			"", // Mixture of authors.
			"http://users.bentonrea.com/~sasquatch/sokoban/morehex.hsb",
			"data/morehex.hsb",
			"levels/move/",
			parseMore,
		},
	}
}

func parseMore(text []byte, collection hexaban.Collection) ([]hexaban.Puzzle, error) {
	puzzles := make([]hexaban.Puzzle, 0)
	// TODO
	return puzzles, nil
}

func parseNYI(text []byte, collection hexaban.Collection) ([]hexaban.Puzzle, error) {
	fmt.Println(collection.Source)
	fmt.Println("UM, we haven't implemented a reader for this collection yet.")

	return nil, nil
}
