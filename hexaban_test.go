// Hexaban (some know it as Hexoban to match the 'o' in Sokoban).

//

// Structures for representing the initial state and sequence of moves

// for the solitaire game that is a hexagonal representation of the classic

// known as Sokoban.

package hexaban

import (
	"testing"
)

func TestRectCoord_ToHex(t *testing.T) {
	tests := []struct {
		name  string
		coord RectCoord
		want  HexCoord
	}{
		{"zero", RectCoord{0, 0}, HexCoord{0, 0}},
		{"i", RectCoord{1, 0}, HexCoord{1, 0}},
		{"i+j", RectCoord{0, 1}, HexCoord{1, 1}},
		{"3i - j", RectCoord{4, 1}, HexCoord{3, -1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.coord.ToHex(); got.i != tt.want.i || got.j != tt.want.j {
				t.Errorf("RectCoord.ToHex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRectCoord_ToHexCenteredAt(t *testing.T) {
	tests := []struct {
		name   string
		coord  RectCoord
		center HexCoord
		want   HexCoord
	}{
		{"zero", RectCoord{0, 0}, HexCoord{0, 0}, HexCoord{0, 0}},
		{"i", RectCoord{1, 0}, HexCoord{0, 0}, HexCoord{1, 0}},
		{"i (recentered)", RectCoord{1, 0}, HexCoord{2, 2}, HexCoord{-1, -2}},
		{"i+j (recentered)", RectCoord{0, 1}, HexCoord{3, 2}, HexCoord{-2, 1}},
		{"3i - j at center", RectCoord{4, 1}, HexCoord{3, -1}, HexCoord{0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.coord.ToHexCenteredAt(tt.center); got.i != tt.want.i || got.j != tt.want.j {
				t.Errorf("RectCoord.ToHex() = %v, want %v", got, tt.want)
			}
		})
	}
}
