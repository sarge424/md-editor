package editor

import (
	"math"

	"github.com/sarge424/notes/internal/buffer"
	"github.com/tfriedel6/canvas"
)

type Editor struct {
	buf buffer.Buffer
}

func New() Editor {
	return Editor{
		buf: buffer.New(500),
	}
}

func DrawAppBG(width, height, rnd int, cv *canvas.Canvas) {
	cv.SetFillStyle("#0008")
	cv.FillRect(float64(rnd), 0, float64(width-2*rnd), float64(height))
	cv.FillRect(0, float64(rnd), float64(width), float64(height-2*rnd))

	cv.BeginPath()
	cv.Arc(float64(rnd), float64(rnd), float64(rnd), 0, 6, true)
	cv.Fill()
	cv.BeginPath()
	cv.Arc(float64(width-rnd), float64(rnd), float64(rnd), 0, 2*math.Pi, true)
	cv.Fill()
	cv.BeginPath()
	cv.Arc(float64(rnd), float64(height-rnd), float64(rnd), 0, 2*math.Pi, true)
	cv.Fill()
	cv.BeginPath()
	cv.Arc(float64(width-rnd), float64(height-rnd), float64(rnd), 0, 2*math.Pi, true)
	cv.Fill()
}

func (e *Editor) LoadFile(fp string) error {
	return e.buf.LoadFile(fp)
}

func (e Editor) DrawBuffer(cv *canvas.Canvas) {
	p := e.buf.Parser()
	for p.Next() {
		cv.FillText(p.Data, 0, float64(p.RowNo+1)*20)
	}
}
