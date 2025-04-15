package ui

import (
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"log"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/imageutil"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

type Visualizer struct {
	Title         string
	Debug         bool
	OnScreenReady func(screen screen.Screen)

	window  screen.Window
	texture chan screen.Texture
	done    chan struct{}

	size   size.Event
	pos    image.Rectangle
	center image.Point
}

func (visualizer *Visualizer) Main() {
	visualizer.texture = make(chan screen.Texture)
	visualizer.done = make(chan struct{})
	visualizer.center.X = 400
	visualizer.center.Y = 400
	driver.Main(visualizer.run)
}

func (visualizer *Visualizer) Update(t screen.Texture) {
	visualizer.texture <- t
}

func (visualizer *Visualizer) run(s screen.Screen) {
	if visualizer.OnScreenReady != nil {
		visualizer.OnScreenReady(s)
	}

	window, err := s.NewWindow(&screen.NewWindowOptions{
		Title:  visualizer.Title,
		Width:  800,
		Height: 800,
	})
	if err != nil {
		log.Fatal("Failed to initialize the app window:", err)
	}
	defer func() {
		window.Release()
		close(visualizer.done)
	}()

	visualizer.window = window

	events := make(chan any)
	go func() {
		for {
			event := window.NextEvent()
			if visualizer.Debug {
				log.Printf("new event: %v", event)
			}
			if detectTerminate(event) {
				close(events)
				break
			}
			events <- event
		}
	}()

	var texture screen.Texture

	for {
		select {
		case event, ok := <-events:
			if !ok {
				return
			}
			visualizer.handleEvent(event, texture)

		case texture = <-visualizer.texture:
			window.Send(paint.Event{})
		}
	}
}

func detectTerminate(event any) bool {
	switch event := event.(type) {
	case lifecycle.Event:
		if event.To == lifecycle.StageDead {
			return true // Window destroy initiated.
		}
	case key.Event:
		if event.Code == key.CodeEscape {
			return true // Esc pressed.
		}
	}
	return false
}

func (visualizer *Visualizer) handleEvent(event any, texture screen.Texture) {
	switch event := event.(type) {

	case size.Event: // Оновлення даних про розмір вікна.
		visualizer.size = event
		visualizer.center = image.Pt(
			visualizer.size.WidthPx/2,
			visualizer.size.HeightPx/2,
		)

	case error:
		log.Printf("ERROR: %s", event)

	case mouse.Event:
		if texture == nil {
			if event.Button == mouse.ButtonLeft &&
				event.Direction == mouse.DirPress {
				visualizer.center.Y, visualizer.center.X = int(event.Y), int(event.X)
				visualizer.window.Send(paint.Event{})
			}
		}

	case paint.Event:
		// Малювання контенту вікна.
		if texture == nil {
			visualizer.drawDefaultUI()
		} else {
			// Використання текстури отриманої через виклик Update.
			visualizer.window.Scale(
				visualizer.size.Bounds(),
				texture,
				texture.Bounds(),
				draw.Src,
				nil,
			)
		}
		visualizer.window.Publish()
	}
}

func DrawFigure(uploader screen.Uploader, point image.Point) {
	uploader.Fill(
		image.Rect(point.X+200, point.Y, point.X-200, point.Y-175),
		colornames.Blue,
		draw.Src,
	)

	uploader.Fill(
		image.Rect(point.X-75, point.Y, point.X+75, point.Y+175),
		colornames.Blue,
		draw.Src,
	)
}

func (visualizer *Visualizer) drawDefaultUI() {
	visualizer.window.Fill(
		visualizer.size.Bounds(),
		colornames.Lime,
		draw.Src,
	) // Фон.

	DrawFigure(visualizer.window, visualizer.center)

	// Малювання білої рамки.
	for _, border := range imageutil.Border(visualizer.size.Bounds(), 10) {
		visualizer.window.Fill(border, color.White, draw.Src)
	}
}
