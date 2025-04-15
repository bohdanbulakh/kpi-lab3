package painter

import (
	"errors"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/colornames"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	"testing"
)

type MockReceiver struct {
	callCount int
}

func (r *MockReceiver) Update(_ screen.Texture) {
	r.callCount++
}

type MockScreen struct{}

func (s MockScreen) NewBuffer(_ image.Point) (screen.Buffer, error) {
	return nil, errors.New("not implemented")
}
func (s MockScreen) NewTexture(_ image.Point) (screen.Texture, error) {
	return nil, errors.New("not implemented")
}
func (s MockScreen) NewWindow(_ *screen.NewWindowOptions) (screen.Window, error) {
	return nil, errors.New("not implemented")
}

type MockTexture struct{}

func (t MockTexture) Release()                                                 {}
func (t MockTexture) Size() image.Point                                        { return image.Point{} }
func (t MockTexture) Bounds() image.Rectangle                                  { return image.Rectangle{} }
func (t MockTexture) Upload(_ image.Point, _ screen.Buffer, _ image.Rectangle) {}
func (t MockTexture) Fill(_ image.Rectangle, _ color.Color, _ draw.Op)         {}

type Checker struct {
	expect int
	ch     chan struct{}
}

func (c Checker) wait() {
	for i := 0; i < c.expect; i++ {
		<-c.ch
	}
}

func (c Checker) signal() {
	c.ch <- struct{}{}
}

func newChecker(count int) Checker {
	return Checker{expect: count, ch: make(chan struct{}, count)}
}

func TestSingleFillColor(t *testing.T) {
	ops := OperationList{
		Fill{Color: color.RGBA{G: 0x44, A: 0xff}},
		Fill{Color: color.White},
	}

	checker := newChecker(len(ops))
	loop := Loop{Receiver: &MockReceiver{}, doneFunc: checker.signal}
	loop.Start(MockScreen{})
	loop.Post(ops)
	checker.wait()

	if loop.state.backgroundColor.Color != color.White {
		t.Errorf("Expected white background, got %v", loop.state.backgroundColor.Color)
	}
}

func TestEmptyOperationList(t *testing.T) {
	loop := Loop{Receiver: &MockReceiver{}}
	loop.Start(MockScreen{})
	loop.Post(OperationList{})

	if loop.state.backgroundColor.Color != colornames.Lime ||
		loop.state.backgroundRect != nil ||
		loop.state.figureCenters != nil {
		t.Error("Initial state is incorrect")
	}
}

func TestMultipleFills(t *testing.T) {
	ops := OperationList{
		Fill{Color: color.RGBA{G: 0x33, A: 0xdd}},
		Fill{Color: color.RGBA{R: 0xaa, A: 0xff}},
		Fill{Color: color.Black},
	}

	checker := newChecker(len(ops))
	loop := Loop{Receiver: &MockReceiver{}, doneFunc: checker.signal}
	loop.Start(MockScreen{})
	loop.Post(ops)
	checker.wait()

	if loop.state.backgroundColor.Color != color.Black {
		t.Errorf("Expected black background, got %v", loop.state.backgroundColor.Color)
	}
}

func TestStoreLastRect(t *testing.T) {
	ops := OperationList{
		Bgrect{X1: 0.1, Y1: 0.1, X2: 0.4, Y2: 0.4},
		Bgrect{X1: 0.2, Y1: 0.2, X2: 0.5, Y2: 0.5},
	}

	checker := newChecker(len(ops))
	loop := Loop{Receiver: &MockReceiver{}, doneFunc: checker.signal}
	loop.Start(MockScreen{})
	loop.Post(ops)
	checker.wait()

	expected := Bgrect{X1: 0.2, Y1: 0.2, X2: 0.5, Y2: 0.5}
	if *loop.state.backgroundRect != expected {
		t.Errorf("Expected %v, got %v", expected, *loop.state.backgroundRect)
	}
}

func TestAddAndMoveFigures(t *testing.T) {
	ops := OperationList{
		Figure{X: 0.1, Y: 0.1},
		Figure{X: 0.3, Y: 0.3},
		Move{X: 0.7, Y: 0.2},
	}

	checker := newChecker(len(ops))
	loop := Loop{Receiver: &MockReceiver{}, doneFunc: checker.signal}
	loop.Start(MockScreen{})
	loop.Post(ops)
	checker.wait()

	moved := Figure{X: 0.7, Y: 0.2}
	if *loop.state.figureCenters[0] != moved || *loop.state.figureCenters[1] != moved {
		t.Error("Figures not moved correctly")
	}
}

func TestResetFunctionality(t *testing.T) {
	ops := OperationList{
		Figure{X: 0.1, Y: 0.1},
		Bgrect{X1: 0.1, Y1: 0.1, X2: 0.2, Y2: 0.2},
		Fill{Color: color.RGBA{G: 0xbb, A: 0xee}},
		Reset{},
	}

	checker := newChecker(len(ops))
	loop := Loop{Receiver: &MockReceiver{}, doneFunc: checker.signal}
	loop.Start(MockScreen{})
	loop.Post(ops)
	checker.wait()

	if loop.state.figureCenters != nil || loop.state.backgroundRect != nil || loop.state.backgroundColor.Color != color.Black {
		t.Error("Reset did not clear state properly")
	}
}

func TestUpdateCall(t *testing.T) {
	ops := OperationList{
		Figure{X: 0.4, Y: 0.4},
		Update{},
	}

	receiver := &MockReceiver{}
	checker := newChecker(len(ops))
	loop := Loop{Receiver: receiver, doneFunc: checker.signal}
	loop.Start(MockScreen{})
	loop.currentTexture = MockTexture{}
	loop.Post(ops)
	checker.wait()

	if receiver.callCount != 1 {
		t.Errorf("Expected Update to be called once, got %d", receiver.callCount)
	}
}

func TestAddFigureOnly(t *testing.T) {
	ops := OperationList{
		Figure{X: 0.1, Y: 0.1},
		Figure{X: 0.2, Y: 0.2},
	}

	checker := newChecker(len(ops))
	loop := Loop{Receiver: &MockReceiver{}, doneFunc: checker.signal}
	loop.Start(MockScreen{})
	loop.Post(ops)
	checker.wait()

	if len(loop.state.figureCenters) != 2 {
		t.Errorf("Expected 2 figures, got %d", len(loop.state.figureCenters))
	}
}

func TestMoveWithoutFigures(t *testing.T) {
	ops := OperationList{
		Move{X: 0.5, Y: 0.5},
	}

	checker := newChecker(len(ops))
	loop := Loop{Receiver: &MockReceiver{}, doneFunc: checker.signal}
	loop.Start(MockScreen{})
	loop.Post(ops)
	checker.wait()

	if loop.state.figureCenters != nil {
		t.Error("Expected no figures to move, but got some")
	}
}

func TestResetOnEmptyState(t *testing.T) {
	ops := OperationList{
		Reset{},
	}

	checker := newChecker(len(ops))
	loop := Loop{Receiver: &MockReceiver{}, doneFunc: checker.signal}
	loop.Start(MockScreen{})
	loop.Post(ops)
	checker.wait()

	if loop.state.figureCenters != nil || loop.state.backgroundRect != nil || loop.state.backgroundColor.Color != color.Black {
		t.Error("Reset on empty state failed")
	}
}
