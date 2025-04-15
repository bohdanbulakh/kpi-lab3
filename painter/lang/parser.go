package lang

import (
	"bufio"
	"errors"
	"image/color"
	"io"
	"strconv"
	"strings"

	"github.com/bohdanbulakh/kpi-lab3/painter"
)

// Parser уміє прочитати дані з вхідного io.Reader та повернути список операцій представлені вхідним скриптом.
type Parser struct{}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	var res []painter.Operation
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		commandLine := scanner.Text()
		op, err := parseCommand(commandLine)

		if err != nil {
			return res, err
		}

		res = append(res, op)
	}

	return res, nil
}

func parseCommand(commandLine string) (painter.Operation, error) {
	words := strings.Fields(commandLine)
	if len(words) == 0 {
		return nil, errors.New("empty line")
	}

	cmd := words[0]
	args := words[1:]

	switch cmd {
	case "white":
		return painter.Fill{Color: color.White}, nil
	case "green":
		return painter.Fill{Color: color.RGBA{G: 255, A: 255}}, nil
	case "update":
		return painter.UpdateOp, nil
	case "bgrect":
		return handleBgRect(args)
	case "figure":
		return handleFigure(args)
	case "move":
		return handleMove(args)
	case "reset":
		return painter.ResetOp, nil
	default:
		return nil, errors.New("unknown command: " + cmd)
	}
}

func parseParams(input []string, requiredCount int) ([]float32, error) {
	if len(input) != requiredCount {
		return nil, errors.New("expected " + strconv.Itoa(requiredCount) + " arguments")
	}

	values := make([]float32, 0, requiredCount)
	for _, val := range input {
		num, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return nil, errors.New("invalid number: " + val)
		}
		if num < 0 || num > 1 {
			return nil, errors.New("value out of range: " + val)
		}
		values = append(values, float32(num))
	}

	return values, nil
}

func handleBgRect(args []string) (painter.Operation, error) {
	coords, err := parseParams(args, 4)
	if err != nil {
		return nil, err
	}
	return painter.Bgrect{
		X1: coords[0],
		Y1: coords[1],
		X2: coords[2],
		Y2: coords[3],
	}, nil
}

func handleFigure(args []string) (painter.Operation, error) {
	coords, err := parseParams(args, 2)
	if err != nil {
		return nil, err
	}
	return painter.Figure{X: coords[0], Y: coords[1]}, nil
}

func handleMove(args []string) (painter.Operation, error) {
	coords, err := parseParams(args, 2)
	if err != nil {
		return nil, err
	}
	return painter.Move{X: coords[0], Y: coords[1]}, nil
}
