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
// github:SymbolNotFound/hexoban/puzzle/puzzle.go

package puzzle

import (
	"encoding/json"
	"fmt"
)

type Puzzle struct {
	Identity   string     `json:"id"`
	Title      string     `json:"title"`
	Author     string     `json:"author"`
	Source     string     `json:"source"`
	Difficulty int        `json:"difficulty,omitempty"`
	Terrain    []HexCoord `json:"terrain"`
	Init       Init       `json:"init"`
}

type Init struct {
	// All goals and crates need to be on coordinates where terrain exists.
	Goals  []HexCoord `json:"goals"`
	Crates []HexCoord `json:"crates"`

	// The player's initial position may be anywhere reachable from here.
	Ichiban HexCoord `json:"ichiban,omitempty"` // Default is (0, 0)
}

// Runs some basic validation over the puzzle definition.
func (puzzle Puzzle) Validate() error {
	// TODO
	return nil
}

// Satisfies the fmt.Stringer interface, presently just prints the json-marshaled version.
func (puzzle Puzzle) String() string {
	output, err := json.MarshalIndent(puzzle, "  ", "")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(output)
}

// Provides additional info and statistics about the puzzle.
func (puzzle Puzzle) Info() string {
	info := fmt.Sprintf("%s by %s (%s)\nlen terrain, num objectives (crates|goals), ...",
		puzzle.Title, puzzle.Author, puzzle.Identity)
	// TODO
	return info
}
