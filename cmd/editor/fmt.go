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
// github:SymbolNotFound/hexoban/cmd/editor/fmt.go

package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/SymbolNotFound/hexoban/puzzle"
)

const PREFIX TokenType = '>'

// Pretty-prints the puzzle as a line of lines of text, utilizing a transform
// to the double-height offset coordinates equivalent.  Uses glyphs from .HSB
// inputs, i.e. {' ', '#', '.', '$', '*', '@', '+'} and sequencing a line of odd
// columns with a line of even columns to give a spaced-out hexgrid-like text.
//
// e.g.: (5 terrain tiles, 1 goal/crate, solution {bR})
// .     # # # #
// .    #   $ . #
// .	   # @   #
// .	    # # #
func MapString(p puzzle.Puzzle) (string, error) {
	if len(p.Terrain) == 0 {
		return "", errors.New("[terrain undefined]")
	}
	errors := multierror{make([]string, 0)}

	i, j := p.Terrain[0].I(), p.Terrain[0].J()
	minRow, minCol := i, (j<<1)-i
	maxRow, maxCol := i, (j<<1)-i
	// Compute the max/min range of coordinates.
	for _, coord := range p.Terrain {
		i, j = coord.I(), coord.J()
		row, col := i, (j<<1)-i
		if row < minRow {
			minRow = row
		} else if row > maxRow {
			maxRow = row
		}
		if col < minCol {
			minCol = col
		} else if col > maxCol {
			maxCol = col
		}
	}
	// After finding range for floors, add padding for surrounding walls.
	minRow -= 1
	maxRow += 1
	minCol -= 2
	maxCol += 2

	// Build left-aligned rectangular grid of FLOOR from terrain.
	rectgrid := NewRectGrid(uint(maxRow-minRow+1), uint(maxCol-minCol+1))
	for _, coord := range p.Terrain {
		floorxy := HexToRect(coord, -minRow, -minCol)
		rectgrid.Assign(floorxy.row, floorxy.col, TOKEN_FLOOR)
	}

	// wrap floors' perimeter with walls
	// We do this after collecting all the floors above
	// just in case they were not in a sorted order.

	// first line has walls adjacent to any floor tiles on row 1.
	for x, token := range rectgrid.glyphs[1] {
		if token == TOKEN_FLOOR {
			// we can prove that x-1 is still > 1 for all FLOOR tiles.
			rectgrid.Assign(0, uint(x-1), TOKEN_WALL)
			rectgrid.Assign(0, uint(x+1), TOKEN_WALL)
		}
	}
	// Index of the last line of the grid.
	last := uint(len(rectgrid.glyphs) - 1)

	// for internal lines: add before, inside and after.
	for i := uint(1); i < last; i++ {
		first := true
		for x, token := range rectgrid.glyphs[i] {
			switch token {
			case PREFIX:
				if !first {
					rectgrid.Assign(i, uint(x), TOKEN_WALL)
				}
			case TOKEN_FLOOR:
				if first {
					rectgrid.Assign(i, uint(x-2), TOKEN_WALL)
					first = false
				}
			}
		}
		rectgrid.glyphs[i] = append(rectgrid.glyphs[i], TOKEN_WALL)
	}

	// last line has walls adjacent to floor tiles on penultimate row.
	for x, token := range rectgrid.glyphs[last-1] {
		if token == TOKEN_FLOOR {
			rectgrid.Assign(last, uint(x-1), TOKEN_WALL)
			rectgrid.Assign(last, uint(x+1), TOKEN_WALL)
		}
	}

	// add goals on floors and build presence table of goal locations
	goals := make(map[RectCoord]bool)
	for _, goal := range p.Init.Goals {
		goalxy := HexToRect(goal, -minRow, -minCol)
		goals[goalxy] = true
		if !rectgrid.ValidFloor(goalxy.row, goalxy.col) {
			errors.add(fmt.Sprintf("goal not in-bounds: line %d col %d",
				goalxy.row, goalxy.col))
			continue
		}
		rectgrid.Assign(goalxy.row, goalxy.col, TOKEN_GOAL)
	}

	// add crates (and crates on goals)
	for _, crate := range p.Init.Crates {
		cratexy := HexToRect(crate, -minRow, -minCol)
		presentTile := rectgrid.glyphs[cratexy.row][cratexy.col]
		if presentTile == TOKEN_FLOOR {
			rectgrid.Assign(cratexy.row, cratexy.col, TOKEN_CRATE)
		} else if presentTile == TOKEN_GOAL {
			rectgrid.Assign(cratexy.row, cratexy.col, TOKEN_CRATE_GOAL)
		} else {
			errors.add(fmt.Sprintf("crate coordinate invalid: line %d col %d", cratexy.row, cratexy.col))
		}
	}

	// add ichiban position (@, or + if on goal)
	workerxy := HexToRect(p.Init.Ichiban, -minRow, -minCol)
	if goals[workerxy] {
		rectgrid.Assign(workerxy.row, workerxy.col, TOKEN_AT_GOAL)
	} else if rectgrid.ValidFloor(workerxy.row, workerxy.col) {
		rectgrid.Assign(workerxy.row, workerxy.col, TOKEN_AT)
	} else {
		errors.add(fmt.Sprintf("ichiban coordinate invalid: line %d col %d", workerxy.row, workerxy.col))
	}

	// finally, render the grid of glyphs to newline-separated strings.
	if len(errors.errs) == 0 {
		return rectgrid.Stringify(), nil
	}
	return "", errors
}

type RectGrid struct {
	glyphs [][]TokenType
}

func NewRectGrid(rows uint, cols uint) *RectGrid {
	rectgrid := RectGrid{
		make([][]TokenType, rows)}
	for i := range rectgrid.glyphs {
		rectgrid.glyphs[i] = make([]TokenType, 0, cols)
	}
	return &rectgrid
}

func (grid *RectGrid) Assign(line uint, col uint, glyph TokenType) {
	for int(col) >= len(grid.glyphs[line]) {
		grid.glyphs[line] = append(grid.glyphs[line], PREFIX)
	}
	grid.glyphs[line][col] = glyph
}

func (grid *RectGrid) ValidFloor(line uint, col uint) bool {
	return line < uint(len(grid.glyphs)) &&
		col < uint(len(grid.glyphs[line])) &&
		grid.glyphs[line][col] == TOKEN_FLOOR
}

func (grid *RectGrid) Stringify() string {
	lines := make([]string, len(grid.glyphs))
	for y, tokens := range grid.glyphs {
		line := make([]byte, len(tokens))
		for x, token := range tokens {
			switch token {
			case PREFIX, TOKEN_FLOOR:
				line[x] = byte(' ')
			default:
				line[x] = byte(token)
			}
		}
		lines[y] = string(line)
	}
	return strings.Join(lines, "\n")
}

// Simple type for tracking multiple errors, can be expanded for greater detail.
type multierror struct {
	errs []string
}

func (merr multierror) Error() string {
	return (strings.Join(merr.errs, "\n"))
}

func (merr *multierror) add(err string) {
	merr.errs = append(merr.errs, err)
}

// Returns a string with basic information about the puzzle.
func Info(puzzle puzzle.Puzzle) string {
	// TODO: what to put here?  sizeof Terrain?  #crates|goals?  Complexity?
	return fmt.Sprintf("%s\nby %s\n\n", puzzle.Title, puzzle.Author)
}
