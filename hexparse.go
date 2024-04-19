package main

import (
	"fmt"
	"os"
	"path"
)

type Factory struct {
	Fn   func() (Collection, error)
	Path string
}

func main() {
	for _, factory := range []Factory{
		{readDWSHEX, "levels/dws"},
	} {
		hex_collection, err := factory.Fn()
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, hex_puzzle := range hex_collection.Puzzles {
			if err = os.WriteFile(path.Join(factory.Path, hex_puzzle.Identity), hex_puzzle.ToJson(), 0644); err != nil {
				fmt.Println(err)
				return
			}
		}
	}

}

// Parses the 16 puzzles from David W. Skinner's collection.
func readDWSHEX() (Collection, error) {
	source := "data/dwshex.hsb.txt"
	filedata, err := os.ReadFile(source)
	if err != nil {
		return Collection{}, err
	}

	collection := Collection{
		Puzzles: make([]HexabanPuzzle, 0),
		Source:  source,
	}

	// TODO parse contents of the collection; will have reusable parts but differ slightly on available metadata.
	fmt.Println(string(filedata))

	return collection, nil
}
