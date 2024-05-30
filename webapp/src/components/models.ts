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
// github:SymbolNotFound/hexoban/webapp/src/components/models.ts

import { HexCoord } from '../hexgrid/topology'

// This file contains shared model types for serialization, model representation
// and user interface coordination.  Common hexgrid and puzzle logic are in the
// directories webapp/src/hexgrid and webapp/src/puzzle, respectively.  These
// definitions in this file are primarily for the (de-)serialization of puzzles.

// A derivative type of HexCoord that is serialized as a 2D array (2-tuple) and
// which can be converted into a HexCoord struct type.
export type HexCoordTuple = [number, number]

export function TupleToCoord (arr: HexCoordTuple): HexCoord {
  return new HexCoord(arr[0], arr[1])
}

export function CoordToTuple (coord: HexCoord): HexCoordTuple {
  return [coord.i, coord.j]
}

// Specifies the starting positions for the goals and crates of the current
// puzzle, as defined in the serialized (json) representation.
export interface PuzzleInitJSON {
  goals: HexCoordTuple[]
  crates: HexCoordTuple[]
}

// Serialized representation of a puzzle's boundaries and initial conditions,
// along with some metadata about the puzzle.
export interface PuzzleJSON {
  readonly id: string
  name: string
  author: string
  source: string
  terrain: HexCoordTuple[]
  init: PuzzleInitJSON
}

type HexNeighbors = [number, number, number, number, number, number]

// Represents a collection of hex coordinates and valid traversals through them.
export interface HexGrid {

  // Returns an unsigned integer representing the HexCoord index.
  // If the coordinate does not exist in the HexGrid, zero (0) is returned.
  index: (coord: HexCoord) => number

  // Adds the hex coordinate, if it doesn't already exist, and return its index.
  // If it already existed, its index is returned and no new index is created.
  add: (coord: HexCoord) => number

  // Removes the coordinate if it existed in this collection.
  remove: (coord: HexCoord) => void

  // Returns the indices for all six neighbors.
  // Zero values take the place of any neighbors that are not in the HexGrid.
  neighbors: (coord: HexCoord) => HexNeighbors
}
