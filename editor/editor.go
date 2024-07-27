package editor

import (
	"bufio"
	"os"

	"github.com/tfriedel6/canvas"
)

type Cursor struct {
	X int
	Y int
}

type Editor struct {
	Lines   []string
	Pointer Cursor
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

func (e Editor) Draw(cv *canvas.Canvas) {
	cv.SetFillStyle("#BBB")
	cv.FillRect(float64(e.Pointer.X*14), float64(e.Pointer.Y*24), 14, 24)

	cv.SetFillStyle("#FFF")
	for i, line := range e.Lines {
		i++
		cv.FillText(line, 0, float64(i*24))
	}
}
