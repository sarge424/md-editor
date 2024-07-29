package editor

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/sarge424/notes/kb"
	"github.com/tfriedel6/canvas"
)

const (
	NavMode  int = iota
	EditMode int = iota
)

type row struct {
	index  int
	length int
}

type pointer struct {
	x    int
	y    int
	oldx int
}

type Editor struct {
	text content
	rows []row

	mode int
	p    pointer

	scroll int
	Height int
}

func New(height, chunkSize int) Editor {
	return Editor{
		text:   newContent(chunkSize),
		Height: height,
	}
}

func (e *Editor) MoveX(dx int) {
	newX := clamp(e.p.x+dx, 0, e.rows[e.p.y].length)
	if e.p.x != newX {
		e.p.x = newX
		e.p.oldx = newX
	}
}

func (e *Editor) MoveY(dy int) {
	newY := clamp(e.p.y+dy, 0, len(e.rows)-1)
	e.p.y = newY

	// if the row is too short
	e.p.x = min(e.rows[e.p.y].length, e.p.oldx)

	//if the cursor goes off the screen, scroll
	if e.p.y < e.scroll {
		e.scroll = e.p.y
	}
	if e.p.y >= e.scroll+e.Height {
		e.scroll = e.p.y - e.Height + 1
	}
}

func (e *Editor) HandleKeystroke(k kb.Keystroke) {
	// standardize letters to uppercase
	if e.mode == NavMode {

		switch k.Std() {
		// movement
		case 'H':
			e.MoveX(-1)
		case 'L':
			e.MoveX(1)
		case 'J':
			e.MoveY(1)
		case 'K':
			e.MoveY(-1)

		// mode switch
		case 'I':
			e.mode = EditMode

		// scroll
		case '[':
			e.scroll = min(e.scroll+1, len(e.rows)-e.Height)
		case ']':
			e.scroll = max(e.scroll-1, 0)
		}

	} else if e.mode == EditMode {
		e.InsertText(string(k))
	}
}

func (e *Editor) HandleShortcut(k kb.Shortcut) {
	if e.mode == EditMode {

		switch fmt.Sprint(k) { // switch on the string representation
		case "1": //ESC
			e.mode = NavMode

		case "BCKSP", "SHF BCKSP":
			e.DeleteText(1)

		case "ENTER":
			e.InsertText("\n")

		default:
			fmt.Println(k)
		}

	} else {
		fmt.Println(k)
	}
}

func (e *Editor) newLine() {
	// add a new line at cursor pos
	text := "\n"
	pointerPos := e.rows[e.p.y].index + e.p.x
	e.text.Insert(text, pointerPos)

	// split row
	r := row{
		index:  pointerPos + 1,
		length: e.rows[e.p.y].index + e.rows[e.p.y].length - pointerPos,
	}
	e.rows[e.p.y].length = pointerPos - e.rows[e.p.y].index
	e.rows = slices.Insert(e.rows, e.p.y+1, r)

	//move cursor
	e.p.x = 0
	e.p.oldx = 0
	e.MoveY(1)

	//offset the start of all following rows
	for i := e.p.y + 1; i < len(e.rows); i++ {
		e.rows[i].index += len(text)
	}
}

func (e *Editor) InsertText(text string) {
	//split text by newline and add them separately
	if nl := strings.Index(text, "\n"); nl >= 0 {
		if nl > 0 {
			e.InsertText(text[:nl])
		}
		e.newLine()
		if nl+1 < len(text) {
			e.InsertText(text[nl+1:])
		}

		return
	}

	pointerPos := e.rows[e.p.y].index + e.p.x
	e.text.Insert(text, pointerPos)

	e.p.x += len(text)
	e.rows[e.p.y].length += len(text)

	//offset the start of all following rows
	for i := e.p.y + 1; i < len(e.rows); i++ {
		e.rows[i].index += len(text)
	}
}

func (e *Editor) DeleteText(length int) {
	pointerPos := e.rows[e.p.y].index + e.p.x - 1
	if pointerPos < 0 {
		return
	}

	e.text.Delete(pointerPos, length)

	e.p.x -= length
	e.rows[e.p.y].length -= length
	//offset the start of all following rows
	for i := e.p.y + 1; i < len(e.rows); i++ {
		e.rows[i].index -= length
	}

	//merge this row with prev if newline was removed
	if e.p.x < 0 {
		e.p.x += e.rows[e.p.y-1].length + 1
		e.p.oldx = e.p.x

		e.rows[e.p.y-1].length += e.rows[e.p.y].length + 1
		e.rows = slices.Concat(e.rows[:e.p.y], e.rows[min(e.p.y+1, len(e.rows)):])

		e.p.y--
	}

}

func (e *Editor) LoadFile(file string) {
	dat, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	data := strings.ReplaceAll(string(dat), "\t", "    ")

	for len(data) >= e.text.chunkSize {
		e.text.append(data[:e.text.chunkSize])
		data = data[e.text.chunkSize:]
	}

	// load leftover characters
	if len(data) > 0 {
		e.text.append(data)
	}

	// create rows
	e.MakeRows()
}

func (e *Editor) MakeRows() {
	e.rows = nil
	chunkIndex := 0
	last := -1
	for i, ch := range e.text.chunks {
		nextLn := strings.Index(ch, "\n")
		for ; nextLn >= 0; nextLn = strings.Index(ch, "\n") {
			r := row{
				index:  last + 1,
				length: chunkIndex + nextLn - last - 1,
			}
			last = chunkIndex + nextLn
			ch = strings.Replace(ch, "\n", "Q", 1)

			e.rows = append(e.rows, r)
		}

		//last row
		if i == len(e.text.chunks)-1 {
			r := row{
				index:  last + 1,
				length: chunkIndex + len(ch) - last - 1,
			}
			e.rows = append(e.rows, r)
		}

		chunkIndex += len(ch)
	}
}

func (e Editor) String() string {
	return e.text.String()
}

func (e Editor) DrawPointer(cv *canvas.Canvas) {
	if e.p.y < e.scroll || e.p.y >= e.scroll+e.Height {
		return
	}
	switch e.mode {
	case NavMode:
		cv.SetFillStyle("#ff4242")
	case EditMode:
		cv.SetFillStyle("#4242ff")
	}
	// y extra +1 formatting to leave space for border
	cv.FillRect(float64(e.p.x+8)*14, float64(e.p.y-e.scroll+1)*24, 14, 24)
}

func (e Editor) Render(cv *canvas.Canvas) {
	// panel
	cv.SetFillStyle("#444")
	cv.FillRect(12, 12, 80*14, 34*24)

	// line no divider
	cv.SetStrokeStyle("#888")
	cv.BeginPath()
	cv.MoveTo(14*7, 24)
	cv.LineTo(14*7, float64(24*(1+e.Height)))
	cv.Stroke()
	cv.BeginPath()
	cv.MoveTo(14*8+14*50, 24)
	cv.LineTo(14*8+14*50, float64(24*(1+e.Height)))
	cv.Stroke()

	e.DrawPointer(cv)

	rowNo := e.scroll
	chunkStart := 0
	rowBuffer := ""

outer:
	for _, ch := range e.text.chunks {
		chunkEnd := chunkStart + len(ch)

		// as long as rows can start in this chunk
		for e.rows[rowNo].index < chunkEnd {
			rowEnd := e.rows[rowNo].index + e.rows[rowNo].length

			// the row ends in this chunk
			st := max(0, e.rows[rowNo].index-chunkStart)
			if rowEnd <= chunkEnd {
				rowBuffer += ch[st : rowEnd-chunkStart]

				//row numbers
				cv.SetFillStyle("#888")
				cv.FillText(fmt.Sprintf("%04d", rowNo+1), 14*2, float64(rowNo-e.scroll+1+1)*24)

				// y +1 extra formatting to leave space for border
				cv.SetFillStyle("#FFF")
				cv.FillText(rowBuffer, 14*8, float64(rowNo-e.scroll+1+1)*24)
				//fmt.Println("<", rowBuffer, ">")

				rowBuffer = ""
				rowNo++
				if rowNo >= e.scroll+e.Height {
					break outer
				}

			} else { // the row does not end in this chunk
				rowBuffer += ch[st:]
				break
			}
		}

		chunkStart = chunkEnd
	}

	// the file ends in a newline
	if rowNo < len(e.rows) {
		// fmt.Println("<", rowBuffer, ">")
		cv.FillText(rowBuffer, 0, float64(rowNo-e.scroll+1)*24)
	}
}

func clamp(x, lo, hi int) int {
	return min(max(lo, x), hi)
}
