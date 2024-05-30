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
// github:SymbolNotFound/hexaban/webapp/src/hexgrid/layout.ts

import { HexGrid, HexCoordIndex } from './topology'

// Odd values are pointy-topped, even values are flat-topped.
// Rotations follow clockwise in 30-degree increments for twelve iid rotations.
export type Rotation = 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11

export class GridLayout {
  grid: HexGrid
  rotation: Rotation
  transform: [[number, number], [number, number]]

  constructor (grid: HexGrid, rotation: Rotation = 1) {
    this.grid = grid
    this.rotation = rotation
    this.transform = [[1, 0], [0, 1]] // identity
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
