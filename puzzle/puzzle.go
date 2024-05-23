package puzzle

import (
	"encoding/json"
	"fmt"
)

type Puzzle struct {
	Identity   string     `json:"id"`
	Name       string     `json:"name"`
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

// Satisfies the fmt.Stringer interface
func (puzzle Puzzle) String() string {
	//info := fmt.Sprintf("%s by %s (%s)\n", puzzle.Name, puzzle.Author, puzzle.Identity)
	// TODO
	output, err := json.MarshalIndent(puzzle, "  ", "")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(output)
}

// Provides additional info and statistics about the puzzle.
func (puzzle Puzzle) Info() string {
	// TODO
	return "len terrain, num objectives (crates|goals), ..."
}
