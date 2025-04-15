package painter

import (
	"github.com/bohdanbulakh/kpi-lab3/ui"
	"image"
	"image/color"

	"golang.org/x/exp/shiny/screen"
)

// Operation змінює вхідну текстуру.
type Operation interface {
	Update(state *State)
}

// OperationList групує список операції в одну.
type OperationList []Operation

// UpdateOp операція, яка не змінює текстуру, але сигналізує, що текстуру потрібно розглядати як готову.
var UpdateOp = Update{}

type Update struct{}

func (op Update) Update(_ *State) {}

// Fill зафарбовує текстуру у відповідний колір
type Fill struct {
	Color color.Color
}

func (op Fill) Do(texture screen.Texture) {
	texture.Fill(texture.Bounds(), op.Color, screen.Src)
}

func (op Fill) Update(state *State) {
	state.backgroundColor = &op
}

type Reset struct{}

// ResetOp операція очищує вікно
var ResetOp = Reset{}

func (op Reset) Update(state *State) {
	state.backgroundColor = &Fill{Color: color.Black}
	state.backgroundRect = nil
	state.figureCenters = nil
}

// Bgrect операція додає чорний прямокутник на екран в певних координатах
type Bgrect struct {
	X1 float32
	Y1 float32
	X2 float32
	Y2 float32
}

func (op Bgrect) Do(t screen.Texture) {
	t.Fill(
		image.Rect(
			int(op.X1*float32(t.Size().X)),
			int(op.Y1*float32(t.Size().Y)),
			int(op.X2*float32(t.Size().X)),
			int(op.Y2*float32(t.Size().Y)),
		),
		color.Black,
		screen.Src,
	)
}

func (op Bgrect) Update(state *State) {
	state.backgroundRect = &op
}

// Figure операція додає фігуру варіанту на вказані координати
type Figure struct {
	X float32
	Y float32
}

func (op Figure) Do(texture screen.Texture) {
	ui.DrawFigure(
		texture,
		image.Pt(
			int(op.X*float32(texture.Size().X)),
			int(op.Y*float32(texture.Size().Y)),
		),
	)
}

func (op Figure) Update(state *State) {
	state.figureCenters = append(state.figureCenters, &op)
}

// Move операція переміщує усі на відповідну кількість пікселів
type Move struct {
	X float32
	Y float32
}

func (op Move) Update(state *State) {
	for _, figureCenter := range state.figureCenters {
		figureCenter.X = op.X
		figureCenter.Y = op.Y
	}
}
