package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/SymbolNotFound/hexaban"
)

type ParseFn func([]byte, hexaban.Collection) ([]hexaban.Puzzle, error)

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

		collection := hexaban.Collection{
			Puzzles: make([]hexaban.Puzzle, 0),
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
			convertDWS,
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

// A file converter for the DWS puzzle set by David W. Skinner.
func convertDWS(text []byte, collection hexaban.Collection) ([]hexaban.Puzzle, error) {
	puzzles := make([]hexaban.Puzzle, 0)
	parser := hexaban.NewParser(text)
	parser.Expect("; Hexobans by David W. Skinner")
	if !parser.ExpectLines(2) {
		return puzzles, errors.New("expected a double-newline after the file header")
	}

	var err error = nil
	errors := make([]error, 0)

	for !parser.EOF() {
		puzzle := hexaban.Puzzle{}
		puzzle.Author = "David W. Skinner"
		puzzle.Name = parser.ExpectQuotedString()
		if puzzle.Name == "" {
			errors = append(errors,
				fmt.Errorf("expected a quoted string for the level name %d", len(puzzles)))
			continue
		}
		if err != nil {
			return puzzles, err
		}
		err = parser.ParseLevelDefinition(&puzzle)
		if err != nil {
			return puzzles, err
		}
		puzzles = append(puzzles, puzzle)
	}

	return puzzles, nil
}

func parseMore(text []byte, collection hexaban.Collection) ([]hexaban.Puzzle, error) {
	// TODO
	return nil, nil
}

func parseNYI(text []byte, collection hexaban.Collection) ([]hexaban.Puzzle, error) {
	fmt.Println(collection.Source)
	fmt.Println("UM, we haven't implemented a reader for this collection yet.")

	return nil, nil
}
