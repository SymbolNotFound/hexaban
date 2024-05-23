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
// github:SymbolNotFound/hexaban/cmd/editor/main.go

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/SymbolNotFound/hexaban"
	"github.com/SymbolNotFound/hexaban/cmd/editor/parser"
)

func main() {
	// Some flags for common or required properties.
	authorFlag := flag.String("author", "", "Specifies the puzzle author's name")
	titleFlag := flag.String("title", "", "Specifies the puzzle's title")
	difficFlag := flag.Int("difficulty", 0, "Specifies the puzzle's difficulty (zero for unknown)")

	// also map short-form aliases for the above flags.
	var aliases = map[string]string{
		"a":    "author",
		"n":    "title",
		"name": "title",
		"d":    "difficulty",
	}
	for from, to := range aliases {
		flagSet := flag.Lookup(to)
		flag.Var(flagSet.Value, from, fmt.Sprintf("alias for %s", to))
	}
	flag.Parse()

	// Populate puzzle metadata
	puzzle := hexaban.Puzzle{
		Author:     *authorFlag,
		Name:       *titleFlag,
		Difficulty: *difficFlag,
	}

	reader := bufio.NewReader(os.Stdin)
	app(reader, &puzzle)

	// Display the puzzle and some stats about it.
	fmt.Println(puzzle)
	fmt.Println(puzzle.Info())

	// Determine the output path.  If user supplied an argument, use that
	outpath := filepath.Clean(flag.Arg(0))
getout:
	for outpath == "" {
		fmt.Print("Where to write the JSON output? ")
		if input, err := reader.ReadString('\n'); err != nil {
			outpath = filepath.Clean(input)
		}
	}
	if strings.HasSuffix(outpath, ".json") {
		// TODO
	} else {
		// If a simple extension other than .json is found, replace it.
		matcher := regexp.MustCompile(`[^.\s][a-z_]+$`)
		extension := matcher.FindStringIndex(outpath)
		if extension == nil {
			outpath += ".json"
		} else {
			outpath = outpath[0:extension[0]] + ".json"
		}
		if !yesnoPrompt(reader,
			fmt.Sprintf("\nFile extension unrecognized, did you mean '%s'? ", outpath),
			false) {
			goto getout
		}
	}
	// result of above implies the path ends with ".json",
	// extract the identity from the name without that extension.
	puzzle.Identity = filepath.Base(outpath)[:len(outpath)-5]

	// Write the puzzle to the indicated file.
	if err := os.WriteFile(outpath, []byte(puzzle.String()), 0644); err != nil {
		log.Fatal(err)
	}
}

func app(reader *bufio.Reader, puzzle *hexaban.Puzzle) {
	promptForMissingValues(puzzle, reader)

	// Prompt for puzzle in text format, then convert and validate the puzzle.
	fmt.Println("Enter the puzzle as glyphs in double-height offset coords:")
	err := parser.ParsePuzzleDefinition(reader, puzzle)
	if err != nil {
		log.Fatal(err)
	}
	err = puzzle.Validate()
	if err != nil {
		fmt.Println(puzzle)
		log.Fatal(err)
	}
}

func promptForMissingValues(puzzle *hexaban.Puzzle, reader *bufio.Reader) {
	puzzle.Author = readNonemptyLine(reader, "The puzzle's author? ")
	puzzle.Name = readNonemptyLine(reader, "The puzzle's name (or unique number)? ")

	for {
		propName := readLine(reader, "Additional (optional) properties? ")
		if strings.Trim(propName, " \t\r\n") == "" {
			break
		}
		switch strings.Trim(propName, " \t\r\n") {
		case "Name":
			puzzle.Name = readNonemptyLine(reader, "The puzzle's Name: ")
		case "Source":
			puzzle.Source = readNonemptyLine(reader, "The puzzle Source: ")
		case "Difficulty":
			difficulty := readInteger(reader, "The puzzle's difficulty (a positive integer): ", 0)
			puzzle.Difficulty = difficulty
		default:
			fmt.Println("That property name is not recognized; options are " +
				"(Name, Source, Difficulty)")
		}
	}
}

func readNonemptyLine(reader *bufio.Reader, prompt string) string {
	var input string = ""
	for input == "" {
		input = strings.Trim(
			readLine(reader, prompt),
			" \t\r\n\v\f\u0085\u00A0")
	}
	return input
}

func readLine(reader *bufio.Reader, prompt string) string {
	fmt.Printf("%s", prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return input
}

func readInteger(reader *bufio.Reader, prompt string, deft int) int {
	intValue, err := strconv.Atoi(readNonemptyLine(reader, prompt))
	if err != nil {
		fmt.Println(err)
		return deft
	}
	return intValue
}

// Wrapper around prompt which interprets the
func yesnoPrompt(reader *bufio.Reader, prompt string, otherwise bool) bool {
	if otherwise {
		prompt += " [Y]/n "
	} else {
		prompt += " y/[N] "
	}

	switch strings.ToLower(readLine(reader, prompt)) {
	case "yes":
	case "y":
	case "yeah":
	case "t":
	case "true":
		return true
	case "no":
	case "n":
	case "nah":
	case "f":
	case "false":
		return false
	}

	return otherwise
}
