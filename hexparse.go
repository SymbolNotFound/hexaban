package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

type ParseFn func([]byte, Collection) ([]HexabanPuzzle, error)

// Defines metadata for file being parsed to produce independent puzzle files.
type CollectionFactory struct {
	Author     string
	Source     string
	InputPath  string
	OutputPath string
	ParseFn    ParseFn
}

func main() {
	for _, factory := range FileMetadata() {
		fmt.Printf("Parsing '%s'...\n", factory.InputPath)
		filedata, err := os.ReadFile(factory.InputPath)
		if err != nil {
			fmt.Println(err)
			continue
		}

		collection := Collection{
			Puzzles: make([]HexabanPuzzle, 0),
			Author:  factory.Author,
			Source:  factory.Source,
		}
		collection.Puzzles, err = factory.ParseFn(filedata, collection)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, hex_puzzle := range collection.Puzzles {
			json, err := json.Marshal(hex_puzzle)
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = os.WriteFile(path.Join(factory.OutputPath, hex_puzzle.Identity), json, 0644)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}

func FileMetadata() []CollectionFactory {
	return []CollectionFactory{
		{
			"David W. Skinner",
			"http://users.bentonrea.com/~sasquatch/sokoban/hex.html",
			"data/dwshex.hsb",
			"levels/DWS/",
			parseNYI,
		},
		{
			"François Marques",
			"http://hexoban.online.fr/",
			"data/heloban.hsb",
			"levels/heloban/",
			parseNYI,
		},
		{
			"François Marques",
			"http://hexoban.online.fr/",
			"data/heroban.hsb",
			"levels/heroban/",
			parseNYI,
		},
		{
			"Erim SEVER",
			"www.erimsever.com/sokoban/Erim_Levels/E_Hexoban.zip",
			"data/all_E_Hex.hsb",
			"levels/ErimSEVER/",
			parseNYI,
		},
		{
			"Aymeric du Peloux",
			"http://membres.lycos.fr/nabokos/",
			"data/hexocet.hsb",
			"levels/hexocet/",
			parseNYI,
		},
		{
			"LukaszM",
			"https://play.fancade.com/5FA6BCFD16EB8B3B",
			"data/lukaszm.hsb",
			"levels/LukaszM/",
			parseNYI,
		},
		{
			"", // Mixture of authors.
			"http://users.bentonrea.com/~sasquatch/sokoban/morehex.hsb",
			"data/morehex.hsb",
			"levels/move/",
			parseMore,
		},
	}
}

func parseMore(text []byte, puzzles Collection) ([]HexabanPuzzle, error) {
	// TODO
	return nil, nil
}

func parseNYI(text []byte, puzzles Collection) ([]HexabanPuzzle, error) {
	fmt.Println(puzzles.Source)
	fmt.Println("UM, we haven't implemented that yet.")

	return nil, nil
}

type PuzzleParser struct {
	filedata []byte
	cursor   uint
}

func (parser *PuzzleParser) EOF() bool {
	return parser.cursor >= uint(len(parser.filedata))
}

func (parser *PuzzleParser) MaybeAuthor() (string, error) {
	return "", nil
}

func (parser *PuzzleParser) ExpectRoomDefn() ([]string, error) {
	lines := make([]string, 0)
	// TODO
	return lines, nil
}
