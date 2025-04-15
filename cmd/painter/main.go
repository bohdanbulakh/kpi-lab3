package main

import (
	"net/http"

	"github.com/bohdanbulakh/kpi-lab3/painter"
	"github.com/bohdanbulakh/kpi-lab3/painter/lang"
	"github.com/bohdanbulakh/kpi-lab3/ui"
)

func main() {
	var (
		visualizer  ui.Visualizer
		painterLoop painter.Loop
		parser      lang.Parser
	)

	visualizer.Title = "Painter"
	visualizer.OnScreenReady = painterLoop.Start
	painterLoop.Receiver = &visualizer

	go func() {
		http.Handle("/", lang.HttpHandler(&painterLoop, &parser))
		_ = http.ListenAndServe("localhost:17000", nil)
	}()

	visualizer.Main()
}
