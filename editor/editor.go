package editor

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/sarge424/notes/kb"
	"github.com/tfriedel6/canvas"
)

type Cursor struct {
	X    int
	Y    int
	oldX int
}

type Line struct {
	Text  string
	Rows  int
	AtRow int
}

type Editor struct {
	Lines   []Line
	Pointer Cursor
	ScrollY int
}

func clamp(x, minX, maxX int) int {
	return max(minX, min(x, maxX))
}

func (e *Editor) HandleKeystroke(ks kb.Keystroke) {
	//convert letters to uppercase, others remain unchanged
	let := ks.Std()

	switch let {
	case 'J':
		e.MovePointerY(1)
	case 'K':
		e.MovePointerY(-1)
	case 'H':
		e.MovePointerX(-1)
	case 'L':
		e.MovePointerX(1)
	case ';':
		e.ScrollY++
	case '\'':
		e.ScrollY--
	}
}

func (e *Editor) HandleShortcut(sc kb.Shortcut) {
	fmt.Println(sc)
}

func (e *Editor) LoadFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	e.Lines = nil
	rowCount := 0
	for scanner.Scan() {
		tx := strings.ReplaceAll(scanner.Text(), "\t", "    ")
		line := Line{
			Text:  tx,
			Rows:  int(math.Ceil(float64(len(tx)) / 40)),
			AtRow: rowCount,
		}

		e.Lines = append(e.Lines, line)
		rowCount += line.Rows
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (e *Editor) MovePointerX(dx int) {
	e.Pointer.X = clamp(e.Pointer.X+dx, 0, len(e.Lines[e.Pointer.Y].Text))
	e.Pointer.oldX = e.Pointer.X
}

func (e *Editor) MovePointerY(dy int) {
	if e.Pointer.Y == 0 && dy < 0 { //first line -> go to SOF
		e.MovePointerX(-e.Pointer.X)
	} else if e.Pointer.Y == len(e.Lines)-1 && dy > 0 { // last line -> goto EOF
		e.MovePointerX(len(e.Lines[e.Pointer.Y].Text))
	} else { // regular case
		e.Pointer.Y = clamp(e.Pointer.Y+dy, 0, len(e.Lines)-1)

		e.Pointer.X = e.Pointer.oldX
		e.Pointer.X = clamp(e.Pointer.X, 0, len(e.Lines[e.Pointer.Y].Text))
	}
}

func (e Editor) DrawToCanvas(cv *canvas.Canvas) {
	cv.SetFillStyle("#44E")
	row := e.Lines[e.Pointer.Y].AtRow + e.Pointer.X/40
	col := e.Pointer.X % 40
	cv.FillRect(float64(col*14), float64((row-e.ScrollY)*24), 14, 24)

	cv.SetFillStyle("#FFF")
	for _, line := range e.Lines {
		for r := range line.Rows {
			en := min(len(line.Text), (r+1)*40)
			cv.FillText(line.Text[r*40:en], 0, float64((line.AtRow+r+1-e.ScrollY)*24))
		}
	}
}
