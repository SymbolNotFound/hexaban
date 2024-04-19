// Hexaban (some know it as Hexoban to match the 'o' in Sokoban).
//
// Structures for representing the initial state and sequence of moves
// for the solitaire game that is a hexagonal representation of the classic
// known as Sokoban.

package main

type Collection struct {
	Puzzles []HexabanPuzzle `json:"puzzles"`
	Source  string          `json:"source"`
}

type HexabanPuzzle struct {
	Lines    []string        `json:"-,omitempty"`
	Identity string          `json:"id"`
	Author   string          `json:"author"`
	Name     string          `json:"name"`
	Terrain  [][]HexabanTile `json:"terrain"`
}

func (puzzle HexabanPuzzle) ToJson() []byte {

	// TODO
	return nil
}

type HexabanTile struct {
}

const ERIM_HEX string = "www.erimsever.com/sokoban/Erim_Levels/E_Hexoban.zip"
const HEXOCET string = "http://membres.lycos.fr/nabokos/"
const DWS_SRC string = "http://users.bentonrea.com/~sasquatch/sokoban/hex.html"
const LUCASZM_SRC string = "https://play.fancade.com/5FA6BCFD16EB8B3B"
const HEROBAN_SRC string = "http://hexoban.online.fr/"
