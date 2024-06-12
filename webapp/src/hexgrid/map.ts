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
// github:SymbolNotFound/hexoban/webapp/src/hexgrid/layout.ts

import { HexGrid, HexCoordIndex } from './topology'

// Odd values are flat-topped, even values are pointy-topped.
// Orientation follow clockwise in 30-degree increments for 12 total rotations.
export type Orientation = 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11

export class GridLayout {
  grid: HexGrid
  transform: [[number, number], [number, number]]

  constructor (grid: HexGrid, orientation: Orientation = 0) {
    this.grid = grid
    if (orientation === 0) {
      this.transform = [ // TODO derive transform from orientation.
        // identity matrix; no scaling or rotation (orientation=0).
        [1, 0],
        [0, 1]]
    } else {
      // TODO
      this.transform = [[1, 0], [0, 1]]
    }
  }
}

// Valid directions are only one of these six strings,
// corresponding to up | backward | left | down | forward | right.
export type Direction = 'U'|'B'|'L'|'D'|'F'|'R'

export class HexMap<V> {
  map: Map<HexCoordIndex, V>

  constructor () {
    this.map = new Map()
  }
}
