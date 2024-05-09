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
// github:SymbolNotFound/cmd/inspector/main.go

package main

// Main entry point for inspector.exe
//
// Inspects a puzzle, validates that it is well-formed, reports statistics.

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/SymbolNotFound/hexaban"
)

func main() {
	for puzzlePath := range allPuzzles("../levels/") {
		filedata, err := os.ReadFile(puzzlePath)
		if err != nil {
			fmt.Println(err)
			return
		}

		var puzzle hexaban.Puzzle
		if err = json.Unmarshal(filedata, &puzzle); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("**%s** *by %s*\n", puzzle.Name, puzzle.Author)
		fmt.Printf("retrieved from %s\n", puzzle.Source)
		fmt.Printf("%d traversible tiles.\n", len(puzzle.Terrain))

		errorlog := validatePuzzle(puzzle)
		if errorlog == nil {
			fmt.Println("looks good!")
		} else {
			for _, err := range errorlog {
				fmt.Println(err)
			}
		}
		fmt.Println()
	}
}

func allPuzzles(rootDir string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		contents, _ := os.ReadDir(rootDir)
		for _, name := range contents {
			if name.IsDir() {
				puzzles, _ := os.ReadDir(path.Join(rootDir, name.Name()))
				for _, filename := range puzzles {
					if strings.HasSuffix(filename.Name(), ".json") {
						out <- path.Join(rootDir, name.Name(), filename.Name())
					}
				}
			}
		}
	}()

	return out
}

// Checks various properties of the puzzle for well-formedness and consistency.
// Returns nil if there were no detected issues, or a string for each offense.
func validatePuzzle(puzzle hexaban.Puzzle) []string {
	errs := make([]string, 0)

	if len(puzzle.Terrain) == 0 {
		errs = append(errs, "no terrain coordinates are defined")
		return errs
	}

	if len(puzzle.Init.Goals) != len(puzzle.Init.Crates) {
		errs = append(errs,
			fmt.Sprintf("# goals (%d) differnt from # crates (%d)\n", len(puzzle.Init.Goals), len(puzzle.Init.Crates)))
	}

	coordinates := make(map[hexaban.HexCoord]bool, len(puzzle.Terrain))
	for _, coord := range puzzle.Terrain {
		coordinates[coord] = true
	}

	for _, goal := range puzzle.Init.Goals {
		if !coordinates[goal] {
			errs = append(errs,
				fmt.Sprintf("Found a goal on a non-coordinate %v", goal))
		}
	}

	for _, crate := range puzzle.Init.Crates {
		if !coordinates[crate] {
			errs = append(errs,
				fmt.Sprintf("Found a crate on a non-coordinate %v", crate))
		}
	}

	// TODO minimum matching bipartite graph connecting the crates and goals,
	// (assignment problem, Hungarian algorithm) validate that matching is total.

	// Done validating.
	if len(errs) == 0 {
		return nil
	}
	return errs
}
