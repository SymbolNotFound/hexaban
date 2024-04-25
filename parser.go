package hexaban

import (
	"fmt"
	"regexp"
)

type PuzzleParser struct {
	filedata []byte
	cursor   uint
}

func NewParser(text []byte) PuzzleParser {
	return PuzzleParser{text, 0}
}

func (parser *PuzzleParser) EOF() bool {
	return parser.cursor >= uint(len(parser.filedata))
}

// Reads an expected string and advances the cursor to its end.
func (parser *PuzzleParser) Expect(image string) error {
	limit := parser.cursor + uint(len(image))
	if string(parser.filedata[parser.cursor:limit]) != image {
		return fmt.Errorf("expected [%s] found [%s]",
			image, parser.filedata[parser.cursor:limit])
	}

	parser.cursor = limit
	return nil
}

// Advances the parser past any newlines (\r\n or \n).
// Set count to a negative number to parse all subsequent newlines.
func (parser *PuzzleParser) ExpectLines(count int) bool {
	found := count
	for count < 0 || found > 0 {
		if parser.filedata[parser.cursor] == '\r' {
			continue
		}
		// Handle next character if it's a newline.
		if parser.filedata[parser.cursor] == '\n' {
			found--
			parser.cursor++
			continue
		}
		// Otherwise we found a non-newline character, break out of the loop.
		break
	}

	// Did we find all the lines we were looking for?
	return found <= 0
}

func (parser *PuzzleParser) ExpectPattern(regex string) []string {
	result := make([]string, 0)
	// Ensure the pattern is anchored to the beginning of the string.
	if regex[0] != '^' {
		regex = "^" + regex
	}
	// It's inefficient but this is a batch job, don't worry about it.
	matcher := regexp.MustCompile(regex)
	for _, match := range matcher.FindSubmatch(parser.filedata[parser.cursor:]) {
		result = append(result, string(match))
	}
	return result
}

// Expect the next parsable token to be a string literal.
func (parser *PuzzleParser) ExpectQuotedString() string {
	match := parser.ExpectPattern("^\"[^\"]*\"")
	return match[0]
}

func (parser *PuzzleParser) ParseLevelDefinition(puzzle *Puzzle) error {
	row := 0
	for {
		// Parse puzzle line by line.
		line := parser.PuzzleLine(row)

		// If we've reached the double-newline, this puzzle definition is done.
		if parser.ExpectLines(1) {
			break
		}
	}

	// No errors.
	return nil
}

func (parser *PuzzleParser) PuzzleLine(row int) []Tile {
	result := make([]Tile, 0)
	column := 0
	// Read until newline,
	for parser.filedata[parser.cursor] != '\n' {
		// storing every other character...
		switch parser.filedata[parser.cursor] {

		}
		// and expecting spaces between cell data.
	}

	// Read past the
	parser.ExpectLines(1)
	return result
}
