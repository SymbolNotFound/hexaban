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
// github:SymbolNotFound/hexoban/cmd/editor/main.go

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	builder := PuzzleBuilder()

	// Flags/prompts for common and required properties.
	builder.RequiredString("author", "", "The puzzle's author", nil, []string{"a"})
	builder.RequiredString("title", "", "The puzzle's title", nil, []string{"t", "name", "n"})
	builder.IntegerInput("difficulty", 0, "the puzzle's difficulty (zero for unknown)",
		func(in int) bool { return in >= 0 }, []string{"d"})
	builder.OptionalFields("difficulty", "source", "id", "map")

	// check for filename input
	flag.Parse()
	stdin := bufio.NewReader(os.Stdin)
	outpath := filepath.Clean(flag.Arg(0))

	// if input file (first argument) exists
	// and its extension indicates it is HSB,
	// attempt to parse its contents as a puzzle.
	if strings.HasSuffix(outpath, ".hsb") {
		f, err := os.Open(outpath)
		if err == nil {
			reader := bufio.NewReader(f)
			err = builder.ParsePuzzleMap(reader)
			if err == nil {
				fmt.Println(MapString(builder.BuildPuzzle()))
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}
	} else if !strings.HasSuffix(outpath, ".json") && outpath != "" {
		// Append the extension .json if it isn't found.
		// TODO perhaps if we find existing .json file we are updating a puzzle.
		// TODO this complicates the `id` representation later, prefer extensionless (here, and for .hsb above).
		outpath += ".json"
	}
	builder.RequiredString("id", outpath, "The puzzle's ID, and file name",
		func(s string) bool { return true }, []string{})

	builder.ParseMeta(stdin)
	for len(builder.BuildPuzzle().Terrain) == 0 {
		// read puzzle contents from stdin.
		fmt.Println("Enter puzzle as a doubled-height offset grid here:")
		err := builder.ParsePuzzleMap(stdin)
		if err == nil {
			fmt.Println(MapString(builder.BuildPuzzle()))
		} else {
			fmt.Println(err)
			fmt.Println("\nerror parsing puzzle, try again:")
		}
	}

	puzzle := builder.BuildPuzzle()
	fmt.Println("Successfully parsed metadata & puzzle definition.")
	fmt.Println(Info(puzzle))

	if yesnoPrompt(stdin, "Write to JSON file?", true) {
		// Write the puzzle to JSON-serialized format.
		bytes, err := json.Marshal(puzzle)
		if err == nil {
			err = os.WriteFile(outpath, bytes, 0644)
			if err == nil {
				fmt.Printf("Wrote %d bytes to %s\n", len(bytes), outpath)
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}
	}
}
