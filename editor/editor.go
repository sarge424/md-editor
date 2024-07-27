package editor

import (
	"bufio"
	"fmt"
	"os"

	"github.com/sarge424/notes/kb"
	"github.com/tfriedel6/canvas"
)

type Cursor struct {
	X    int
	Y    int
	oldX int
}

type Editor struct {
	Lines   []string
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
	for scanner.Scan() {
		e.Lines = append(e.Lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (e *Editor) MovePointerX(dx int) {
	e.Pointer.X = clamp(e.Pointer.X+dx, 0, len(e.Lines[e.Pointer.Y]))
	e.Pointer.oldX = e.Pointer.X
}

func (e *Editor) MovePointerY(dy int) {
	if e.Pointer.Y == 0 && dy < 0 { //first line -> go to SOF
		e.MovePointerX(-e.Pointer.X)
	} else if e.Pointer.Y == len(e.Lines)-1 && dy > 0 { // last line -> goto EOF
		e.MovePointerX(len(e.Lines[e.Pointer.Y]))
	} else { // regular case
		e.Pointer.Y = clamp(e.Pointer.Y+dy, 0, len(e.Lines)-1)

		e.Pointer.X = e.Pointer.oldX
		e.Pointer.X = clamp(e.Pointer.X, 0, len(e.Lines[e.Pointer.Y]))
	}
}

func (e Editor) DrawToCanvas(cv *canvas.Canvas) {
	cv.SetFillStyle("#44E")
	cv.FillRect(float64(e.Pointer.X*14), float64((e.Pointer.Y-e.ScrollY)*24), 14, 24)

	cv.SetFillStyle("#FFF")
	for i, line := range e.Lines {
		i++
		cv.FillText(line, 0, float64((i-e.ScrollY)*24))
	}
}
