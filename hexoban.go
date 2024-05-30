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
// github:SymbolNotFound/hexoban/hexoban.go

package hexoban

// Hexaban (some know it as Hexoban to match the 'o' in Sokoban).
//
// Structures for representing the initial state and sequence of moves
// for the solitaire game that is a hexagonal variant of the classic
// puzzle known as Sokoban ("warehouse worker," the crate pusher).

import "github.com/SymbolNotFound/hexoban/puzzle"

type Puzzle = puzzle.Puzzle
type Init = puzzle.Init
type HexCoord = puzzle.HexCoord
type Tile = puzzle.Tile

var NewHexCoord = puzzle.NewHexCoord
