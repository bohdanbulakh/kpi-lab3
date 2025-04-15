package painter

type State struct {
	backgroundColor *Fill
	backgroundRect  *Bgrect
	figureCenters   []*Figure
}
