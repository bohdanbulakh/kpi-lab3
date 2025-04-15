package lang

import (
	"github.com/bohdanbulakh/kpi-lab3/painter"
	"image/color"
	"strings"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	parser := Parser{}

	tests := []struct {
		name     string
		input    string
		expected painter.Operation
		wantErr  bool
	}{
		{"white fill", "white", painter.Fill{Color: color.White}, false},
		{"green fill", "green", painter.Fill{Color: color.RGBA{G: 255, A: 255}}, false},
		{"update", "update", painter.UpdateOp, false},
		{"bgrect", "bgrect 0.1 0.2 0.3 0.4", painter.Bgrect{X1: 0.1, Y1: 0.2, X2: 0.3, Y2: 0.4}, false},
		{"figure", "figure 0.5 0.5", painter.Figure{X: 0.5, Y: 0.5}, false},
		{"move", "move 0.2 0.2", painter.Move{X: 0.2, Y: 0.2}, false},
		{"reset", "reset", painter.ResetOp, false},
		{"invalid command", "unknown", nil, true},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			reader := strings.NewReader(testCase.input)
			result, err := parser.Parse(reader)

			if testCase.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if len(result) != 1 {
					t.Errorf("expected 1 operation, got %d", len(result))
					return
				}
				if result[0] != testCase.expected {
					t.Errorf("expected %v, got %v", testCase.expected, result[0])
				}
			}
		})
	}
}
