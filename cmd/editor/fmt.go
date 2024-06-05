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

	"github.com/SymbolNotFound/hexoban"
)

const (
	SPACE TokenType = ' '
	BLANK TokenType = '_'
)

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
func MapString(p hexoban.Puzzle) (string, error) {
	if len(p.Terrain) == 0 {
		return "", errors.New("[terrain undefined]")
	}
	errors := multierror{make([]string, 0)}

	i, j := p.Terrain[0].I(), p.Terrain[0].J()
	minRow, minCol := i, (j<<1)-i
	maxRow, maxCol := i, (j<<1)-i
	odd := minCol&1 == 0
	// Compute the min & max (range of) coordinates.
	for _, coord := range p.Terrain {
		i, j = coord.I(), coord.J()
		row, col := i, (j<<1)-i
		if row < minRow {
			minRow = row
			odd = col&1 == 0
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
	if odd { // plus extra padding for having an oddly-aligned first row
		minRow -= 1
		minCol -= 1
		maxRow += 1
	}

	// Build rectangular grid able to hold [minRow..maxRow] lines (inclusive)
	// and [minCol..maxCol]
	rectgrid := NewRectGrid(uint(maxRow-minRow)+1, uint(maxCol-minCol)+1)
	for _, coord := range p.Terrain {
		floorxy := HexToRect(coord, -minRow, -minCol)
		rectgrid.Assign(floorxy.row, floorxy.col, TOKEN_FLOOR)
	}

	// Wrap floors' perimeter with walls,
	// and replace internal voids with walls
	// (any neighbor of a floor that is not a floor).
	// We do this after collecting all the floors above
	// just in case they weren't in a sorted order.
	for _, coord := range p.Terrain {
		for _, neighbor := range neighbors(coord) {
			rectcoord := HexToRect(neighbor.Coord(), -minRow, -minCol)
			if rectgrid.Lookup(rectcoord.row, rectcoord.col) == BLANK {
				rectgrid.Assign(rectcoord.row, rectcoord.col, TOKEN_WALL)
			} // ignores invalid (out-of-bounds) and floor tiles.
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

func neighbors(coord hexoban.HexCoord) []hexoban.HexCoord {
	return []hexoban.HexCoord{
		hexoban.NewHexCoord(coord.I()-1, coord.J()-1),
		hexoban.NewHexCoord(coord.I()-1, coord.J()),
		hexoban.NewHexCoord(coord.I(), coord.J()-1),
		hexoban.NewHexCoord(coord.I(), coord.J()+1),
		hexoban.NewHexCoord(coord.I()+1, coord.J()),
		hexoban.NewHexCoord(coord.I()+1, coord.J()+1),
	}
}

type RectGrid struct {
	glyphs [][]TokenType
}

// Constructor function for RectGrid, a rectangular grid of glyphs
// that we incrementally draw the puzzle map onto.
func NewRectGrid(rows uint, cols uint) *RectGrid {
	rectgrid := RectGrid{
		make([][]TokenType, rows)}
	for i := range rectgrid.glyphs {
		rectgrid.glyphs[i] = make([]TokenType, 0, cols)
	}
	return &rectgrid
}

// Returns the type of tile at the given coordinate, if the coordinate exists.
// Returns SPACE if coordinate params are out of bounds or in between hexes.
// Returns BLANK if it is a valid coordinate and not a floor.
func (grid *RectGrid) Lookup(line uint, col uint) TokenType {
	if line >= uint(len(grid.glyphs)) {
		return SPACE
	}
	if col >= uint(len(grid.glyphs[line])) &&
		(line == 0 || col > uint(len(grid.glyphs[line-1]))) &&
		(line == uint(len(grid.glyphs)-1) || col > uint(len(grid.glyphs[line+1]))) {
		return SPACE
	}
	if col >= uint(len(grid.glyphs[line])) {
		return BLANK
	}
	return grid.glyphs[line][col]
}

func (grid *RectGrid) Assign(line uint, col uint, glyph TokenType) {
	for int(col) >= len(grid.glyphs[line]) {
		if int(col&1) != len(grid.glyphs[line])&1 {
			grid.glyphs[line] = append(grid.glyphs[line], SPACE)
		} else {
			grid.glyphs[line] = append(grid.glyphs[line], BLANK)
		}
	}
	grid.glyphs[line][col] = glyph
}

func (grid *RectGrid) ValidFloor(line uint, col uint) bool {
	return line < uint(len(grid.glyphs)) &&
		col < uint(len(grid.glyphs[line])) &&
		grid.glyphs[line][col] == TOKEN_FLOOR
}

func (grid *RectGrid) Stringify() string {
	skippable := 0
	for _, tokens := range grid.glyphs {
		if len(tokens) == 0 {
			skippable += 1
		} else {
			break
		}
	}

	lines := make([]string, len(grid.glyphs)-skippable)
	for y, tokens := range grid.glyphs[skippable:] {
		line := make([]byte, len(tokens))
		for x, token := range tokens {
			switch token {
			case BLANK, TOKEN_FLOOR:
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
func Info(puzzle hexoban.Puzzle) string {
	// TODO: what to put here?  sizeof Terrain?  #crates|goals?  Complexity?
	return fmt.Sprintf("%s\nby %s\n\n", puzzle.Title, puzzle.Author)
}
