import { HexCoord, HexGrid } from '../hexgrid/topology'
import { Direction } from '../hexgrid/layout'

interface PuzzleMove {
  position: HexCoord
  dir: Direction
}

export class PuzzleState {
  grid: HexGrid
  moves: PuzzleMove[]

  constructor (domain: HexGrid) {
    this.grid = domain
    this.moves = []
  }
}

export class Puzzle {
  state: PuzzleState

  constructor (hexgrid: HexGrid) {
    this.state = new PuzzleState(hexgrid)
  }
}
