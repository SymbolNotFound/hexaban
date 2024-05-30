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
// github:SymbolNotFound/hexaban/webapp/src/hexgrid/topology.ts

// HexCoord: fundamental coordinate type.
//
// The principle coordinate system is axial with `i` indicating movement
// in the Downward direction, `j` indicating the Rightward direction.
// The WorldCoord representation depends on this as a consistent representation,
// and the viewport has its own ClientCoord representation for render-visible
// manipulations that depend on WorldCoord.
//
// These HexCoord provide a consistent value between these layers
// and when communicating with the server about moves for save/load state.
export class HexCoord {
  readonly i: number
  readonly j: number

  constructor (i = 0, j = 0) {
    this.i = i
    this.j = j
  }

  // Neighbor coordinates from this HexCoord (i, j) via six cardinal directions.
  // May not necessarily exist in the puzzle's HexGrid;
  // must still be checked with HexGrid.index() or HexGrid.add(), or all neighbors
  // index & presence can be retrieved together with HexGrid.neighbors(coord).

  up (): HexCoord { return new HexCoord(this.i - 1, this.j) }
  down (): HexCoord { return new HexCoord(this.i + 1, this.j) }
  left (): HexCoord { return new HexCoord(this.i, this.j - 1) }
  right (): HexCoord { return new HexCoord(this.i, this.j + 1) }
  forward (): HexCoord { return new HexCoord(this.i + 1, this.j + 1) }
  backward (): HexCoord { return new HexCoord(this.i - 1, this.j - 1) }
}

// Neighbors' relative directions by (right, forward, down, left, backward, up).
type HexNeighbors = [number, number, number, number, number, number]

// A type alias for the densely-packed index of hex coordinates.
export type HexCoordIndex = number

// Represents a collection of hex coordinates and valid traversals through them.
export class HexGrid {
  constructor () {
    this.map = new Map()
    this.revmap = new Map()
    this.currentIndex = 1
  }

  // Returns an unsigned integer representing the HexCoord index.
  // If the coordinate does not exist in the HexGrid, zero (0) is returned.
  index (coord: HexCoord): HexCoordIndex {
    return this.map.get(this._hash(coord)) || 0
  }

  coord (index: HexCoordIndex): HexCoord | undefined {
    return this.revmap.get(index)
  }

  // Adds the hex coordinate, if it doesn't already exist, and return its index.
  // If it already existed, its index is returned and no new index is created.
  add (coord: HexCoord): HexCoordIndex {
    let index = this.index(coord)
    if (index === 0) {
      index = this.currentIndex++
      this.map.set(this._hash(coord), index)
    }
    return index
  }

  // Removes the coordinate if it existed in this collection.
  remove (coord: HexCoord): HexCoordIndex {
    const index = this.index(coord)
    if (index === 0) {
      return 0
    }
    this.revmap.delete(index)
    this.map.delete(this._hash(coord))
    return index
  }

  // Returns the indices for all six neighbors.
  // Zero values take the place of any neighbors that are not in the HexGrid.
  neighbors (coord: HexCoord): HexNeighbors {
    const neighbors: HexNeighbors = [
      this.index(coord.right()),
      this.index(coord.forward()),
      this.index(coord.down()),
      this.index(coord.left()),
      this.index(coord.backward()),
      this.index(coord.up())
    ]
    return neighbors
  }

  // Forward- and reverse-index for coordinates to id.
  map: Map<string, HexCoordIndex>
  revmap: Map<HexCoordIndex, HexCoord>

  // Tracks the value to next use as an index.
  currentIndex: number

  _hash (coord: HexCoord): string {
    return coord.i + ',' + coord.j
  }
}
