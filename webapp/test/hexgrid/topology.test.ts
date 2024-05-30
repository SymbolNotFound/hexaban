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
// github:SymbolNotFound/hexoban/webapp/test/hexgrid/geometry.ts

import { describe, it } from 'vitest'
import { strictEqual, deepStrictEqual } from 'assert'
import { HexCoord } from '../../src/hexgrid/topology'

describe('HexCoord', () => {
  it('getter properties', () => {
    strictEqual(new HexCoord(1, 3).i, 1)
    strictEqual(new HexCoord(1, 3).j, 3)
    strictEqual(new HexCoord(1).j, 0)
  })

  it('find neighbors', () => {
    deepStrictEqual(new HexCoord(4, 2).up(), new HexCoord(3, 2))
    deepStrictEqual(new HexCoord(4, 2).down(), new HexCoord(5, 2))
    deepStrictEqual(new HexCoord(4, 2).left(), new HexCoord(4, 1))
    deepStrictEqual(new HexCoord(4, 2).right(), new HexCoord(4, 3))
    deepStrictEqual(new HexCoord(4, 2).forward(), new HexCoord(5, 3))
    deepStrictEqual(new HexCoord(4, 2).backward(), new HexCoord(3, 1))
  })
})
