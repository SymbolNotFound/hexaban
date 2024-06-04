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
// github:SymbolNotFound/hexoban/cmd/editor/builder.go

package main

import (
	"bufio"
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/SymbolNotFound/hexoban/puzzle"
)

type Builder interface {
	// Builds a puzzle definition from flags and user input.
	ParseMeta(*bufio.Reader) error
	// Builds the Terrain and Init properties from a grid representation.
	ParsePuzzleMap(*bufio.Reader) error

	// Retrieve the puzzle definition
	// after BuildMeta() and/or ParsePuzzleMap()
	// have constructed it.
	BuildPuzzle() puzzle.Puzzle

	// Defines a flag and CLI combination for a string value from user input.
	RequiredString(
		named string,
		default_ string,
		usage string,
		validator func(string) bool,
		aliases []string)

	// Defines a flag and CLI combination for an integer value from user input.
	IntegerInput(
		named string,
		default_ int,
		usage string,
		validator func(int) bool,
		aliases []string)

	// List of field names (assumed to be string valued)
	// which may also be written to the json representation.
	// User will be prompted for additional fields, which must be from this list.
	OptionalFields(names ...string)
}

// Prepare a puzzle builder with its reader.  Sets up flag parsing specification.
func PuzzleBuilder() Builder {
	return &builderState{make([]IInput, 0), &puzzle.Puzzle{}}
}

// Maintains the list of inputs and values entered.
type builderState struct {
	inputs []IInput
	puzzle *puzzle.Puzzle
}

func (state *builderState) BuildPuzzle() puzzle.Puzzle {
	return *state.puzzle
}

func (state *builderState) ParseMeta(io *bufio.Reader) error {
	if !flag.Parsed() {
		flag.Parse()
	}

	for _, input := range state.inputs {
		// TODO
		input = input
		// prompt for each input that is required and undefined,
	}
	//for
	// prompt for additional inputs/corrections until empty input.
	//}

	return nil
}

func (state *builderState) ParsePuzzleMap(io *bufio.Reader) error {
	// TODO

	return nil
}

type IInput interface {
	Name() string
	Prompt() string
	Required() bool
	Validate(string) bool
	Value() any
}

// Base input type with common properties.
type input struct {
	flagName string
	prompt   string
	required bool
}

func (in *input) Name() string   { return in.flagName }
func (in *input) Prompt() string { return in.prompt }
func (in *input) Required() bool { return in.required }

// An input type for string-typed values.
type StringInput struct {
	input

	value     string
	validator func(string) bool
}

func (in *StringInput) Validate(input string) bool {
	return in.validator(input)
}

func (in *StringInput) Value() any {
	return in.value
}

// Builds a flag and CLI prompt, repeats prompt until validator returns true.
// The name, default and usage are passed directly to the flag definition
// and used in the prompt.  The validator may be nil in which case a basic
// nonempty-string check will be used.
func (state *builderState) RequiredString(
	name string,
	default_ string,
	usage string,
	validator func(string) bool,
	aliases []string) {

	if validator == nil {
		validator = func(s string) bool { return s != "" }
	}
	strInput := StringInput{
		input{
			name,
			usage,
			true,
		},
		default_,
		validator,
	}

	flag.StringVar(&strInput.value, name, default_, usage)
	if len(aliases) > 0 {
		flagSet := flag.Lookup(name)
		for _, alias := range aliases {
			flag.Var(flagSet.Value, alias, "alias for "+name)
		}
	}

	state.inputs = append(state.inputs, &strInput)
}

// An input type for integer-typed values.
type IntegerInput struct {
	input

	value     int
	validator func(int) bool
}

func (in *IntegerInput) Validate(input string) bool {
	strInput, err := strconv.Atoi(input)
	if err != nil {
		return false
	}
	return in.validator(strInput)
}

func (in *IntegerInput) Value() any {
	return in.value
}

// Builds a flag and CLI prompt, repeats prompt until validator returns true.
// The name, default and usage are passed directly to the flag definition
// and used in the prompt.  The validator may be nil in which case a basic
// nonempty-string check will be used.
func (state *builderState) IntegerInput(
	name string,
	default_ int,
	usage string,
	validator func(int) bool,
	aliases []string) {

	if validator == nil {
		validator = func(input int) bool { return true }
	}

	intInput := IntegerInput{
		input{
			name,
			usage,
			false,
		},
		default_,
		validator,
	}

	flag.IntVar(&intInput.value, name, default_, usage)
	if len(aliases) > 0 {
		flagSet := flag.Lookup(name)
		for _, alias := range aliases {
			flag.Var(flagSet.Value, alias, "alias for "+name)
		}
	}

	state.inputs = append(state.inputs, &intInput)
}

// Designated field names for string values that may optionally be defined.
// The names should be lowercase and match the JSON field names for Puzzle.
func (state *builderState) OptionalFields(names ...string) {
	for _, name := range names {
		state.inputs = append(state.inputs,
			&StringInput{
				input{
					name,
					"",
					false,
				},
				"",
				func(input string) bool { return len(input) > 0 },
			})
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
