package editor

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/sarge424/notes/config"
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
	filename string
	text     content

	rows   []row
	rowLen int

	mode int
	p    pointer

	scroll int
	Height int
}

func NewEditor(width, height, chunkSize int) Editor {
	return Editor{
		text:   newContent(chunkSize),
		Height: height,
		rowLen: width,
	}
}

func (e *Editor) MoveX(dx int) {
	// move cursor horizontally (clamp to row length)
	newX := clamp(e.p.x+dx, 0, e.rows[e.p.y].length)
	if e.p.x != newX {
		e.p.x = newX
		e.p.oldx = newX
	}
}

func (e *Editor) MoveY(dy int) {
	//move the cursor vertically (clamp to no. of lines)
	newY := clamp(e.p.y+dy, 0, len(e.rows)-1)
	e.p.y = newY

	// if the row is too short, change x
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
	if e.mode == NavMode {

		// standardize letters to uppercase
		switch k {

		// movement
		case 'h', 'H':
			e.MoveX(-1)
		case 'l', 'L':
			e.MoveX(1)
		case 'j', 'J':
			e.MoveY(1)
		case 'k', 'K':
			e.MoveY(-1)

		case 'w':
			// rem is all the text after pointer in the row
			rem := e.text.Get(e.rows[e.p.y].index+e.p.x+1, e.rows[e.p.y].length-e.p.x)
			if ix := strings.Index(rem, " "); ix >= 0 {
				e.MoveX(ix + 1)
			} else {
				e.MoveX(e.rows[e.p.y].length)
			}

		case 'W':
			// rem is all the text before pointer in the row
			rem := e.text.Get(e.rows[e.p.y].index, e.p.x)

			// find last space in rem
			ix := -1
			for i := len(rem) - 1; i >= 0; i-- {
				if rem[i] == ' ' {
					ix = i
					break
				}
			}

			if ix >= 0 {
				e.MoveX(-(len(rem) - ix))
			} else {
				e.MoveX(-e.rows[e.p.y].length)
			}

		// mode switch
		case 'i', 'I':
			e.mode = EditMode

		// scroll
		case '[':
			e.scroll = min(e.scroll+1, max(e.scroll, len(e.rows)-e.Height))
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
		//uses returns to avoid the all-modes switch
		case "1": //ESC
			e.mode = NavMode
			return

		case "BCKSP", "SHF BCKSP":
			e.DeleteText(1)
			return

		case "ENTER":
			e.InsertText("\n")
			return

		}
	}

	// ALL MODES
	switch fmt.Sprint(k) {
	case "CTL S":
		e.SaveFile()
		fmt.Println("File saved.")

	// Cursor movement
	case "328":
		e.MoveY(-1)
	case "336":
		e.MoveY(1)
	case "331":
		e.MoveX(-1)
	case "333":
		e.MoveX(1)

	default:
		fmt.Println(k)
	}
}

func (e *Editor) insertNewLine() {
	// add a new line at cursor pos

	text := "\n"
	pointerPos := e.rows[e.p.y].index + e.p.x
	e.text.Insert(text, pointerPos)

	// split row at cursor
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

	//offset the start index of all following rows
	for i := e.p.y + 1; i < len(e.rows); i++ {
		e.rows[i].index += len(text)
	}
}

func (e *Editor) InsertText(text string) {
	//if the text has newlines, split around it and add it separately
	if nl := strings.Index(text, "\n"); nl >= 0 {
		if nl > 0 {
			e.InsertText(text[:nl])
		}
		e.insertNewLine()
		if nl+1 < len(text) {
			e.InsertText(text[nl+1:])
		}

		return
	}

	//get pointer index
	pointerPos := e.rows[e.p.y].index + e.p.x
	e.text.Insert(text, pointerPos)

	//move pointer
	e.rows[e.p.y].length += len(text)
	e.MoveX(len(text))

	//offset the start of all following rows
	for i := e.p.y + 1; i < len(e.rows); i++ {
		e.rows[i].index += len(text)
	}
}

func (e *Editor) DeleteText(length int) {
	// dont allow backspace on the first character in the file
	pointerPos := e.rows[e.p.y].index + e.p.x - 1
	if pointerPos < 0 {
		return
	}

	// clear from the piece table
	e.text.Delete(pointerPos, length)

	// move cursor
	e.p.x -= length // this is checked later when merging rows - dont use the function
	e.rows[e.p.y].length -= length

	// offset the start of all following rows
	for i := e.p.y + 1; i < len(e.rows); i++ {
		e.rows[i].index -= length
	}

	//merge this row with prev if newline was removed
	if e.p.x < 0 {
		// if empty row, remove it

		//fix pointer xpos
		e.p.x += e.rows[e.p.y-1].length + 1
		e.p.oldx = e.p.x

		//merge rows
		e.rows[e.p.y-1].length += e.rows[e.p.y].length + 1
		e.rows = slices.Concat(e.rows[:e.p.y], e.rows[min(e.p.y+1, len(e.rows)):])

		//fix pointer y
		e.p.y--

	} else {
		e.p.oldx = e.p.x // since e.p.x was set manually earlier
	}
}

func (e *Editor) LoadFile(file string) {
	// open the file
	e.filename = file
	dat, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	// cant handle tabs...
	data := strings.ReplaceAll(string(dat), "\t", "    ")

	// append text as pieces
	for len(data) >= e.text.chunkSize {
		e.text.append(data[:e.text.chunkSize])
		data = data[e.text.chunkSize:]
	}

	// load leftover characters
	if len(data) > 0 {
		e.text.append(data)
	}

	//generating rows
	// only one row edge case
	if !strings.Contains(data, "\n") {
		e.rows = append(e.rows, row{
			index:  0,
			length: len(data),
		})
	} else {
		// create rows
		e.MakeRows()
	}
}

func (e *Editor) SaveFile() {
	f, err := os.Create(e.filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for _, ch := range e.text.chunks {
		ch = strings.ReplaceAll(ch, "    ", "\t")
		f.WriteString(ch)
	}
	f.Sync()
}

func (e *Editor) MakeRows() {
	// generate the rows

	e.rows = nil
	chunkIndex := 0 // index of first character in chunk
	last := -1      // index after previous newline

	// loop through each chunk
	for i, ch := range e.text.chunks {

		// loop until no newlines are found in this chunk
		nextLn := strings.Index(ch, "\n")
		for ; nextLn >= 0; nextLn = strings.Index(ch, "\n") {
			// create a row
			r := row{
				index:  last + 1,
				length: chunkIndex + nextLn - last - 1,
			}

			// update index pointers
			last = chunkIndex + nextLn

			// dummy replace so that the same newline isnt found again
			ch = strings.Replace(ch, "\n", "Q", 1)

			// add the new row to rows
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

		// update the chunk start index
		chunkIndex += len(ch)
	}
}

func (e Editor) String() string {
	return e.text.String() + "\n\n" + fmt.Sprint(e.rows)
}

func (e Editor) DrawPointer(cv *canvas.Canvas, yloc int) {
	if e.p.y < e.scroll || e.p.y >= e.scroll+e.Height {
		return
	}
	switch e.mode {
	case NavMode:
		cv.SetFillStyle(config.Color.NavPointer)
	case EditMode:
		if time.Now().UnixMilli()%1000 > 500 {
			cv.SetStrokeStyle(config.Color.EditPointer)
		} else {
			cv.SetStrokeStyle("#0000")
		}
		cv.SetLineWidth(2)
	}

	// xpos -> offset % max length for long rows
	// ypos -> additional y offset in a long row
	// yloc -> additional offset from previous rows spanning extra lines
	// extra constants to maintain spacing

	xpos := e.p.x % e.rowLen
	ypos := e.p.x/e.rowLen + yloc

	switch e.mode {
	case NavMode:
		cv.FillRect(float64(xpos+8)*14, float64(e.p.y+ypos-e.scroll+1)*24, 14, 24)
	case EditMode:
		cv.BeginPath()
		cv.MoveTo(float64(xpos+8)*14, float64(e.p.y+ypos-e.scroll+1)*24)
		cv.LineTo(float64(xpos+8)*14, float64(e.p.y+ypos-e.scroll+2)*24)
		cv.Stroke()
	}
}

func (e Editor) Render(cv *canvas.Canvas) {

	// draw the background pane, rownumbers divider, etc
	e.DrawPanel(cv)

	rowsDrawn := 0    // number of rows drawn (includes extra height of multiline rows)
	rowNo := e.scroll // current row to be drawn
	chunkStart := 0   // start index of the current chunk
	rowBuffer := ""   // string to be displayed in this row

	// in case the file is empty
	if rowNo == e.p.y {
		e.DrawPointer(cv, 0)
	}

outer:
	for _, ch := range e.text.chunks {
		chunkEnd := chunkStart + len(ch)

		// as long as rows can start in this chunk
		for e.rows[rowNo].index < chunkEnd {

			rowEnd := e.rows[rowNo].index + e.rows[rowNo].length

			// if the row ends in this chunk, finalize buffer and display
			st := max(0, e.rows[rowNo].index-chunkStart)
			if rowEnd <= chunkEnd {
				rowBuffer += ch[st : rowEnd-chunkStart]

				rowsDrawn = e.DrawLine(rowBuffer, rowNo, rowsDrawn, cv)

				rowBuffer = ""
				rowNo++
				if rowNo == e.p.y {
					e.DrawPointer(cv, rowsDrawn-rowNo)
				}
				if rowNo >= min(e.scroll+e.Height, len(e.rows)) || rowsDrawn >= e.Height {
					break outer
				}

			} else { // the row does not end in this chunk, continue to the next one
				rowBuffer += ch[st:]
				break
			}
		}

		chunkStart = chunkEnd
	}

	// the file ends in a newline
	if rowNo < len(e.rows) && rowNo < e.scroll+e.Height {
		e.DrawLine(rowBuffer, rowNo, rowsDrawn, cv)
	}
}

func (e Editor) DrawLine(rowBuffer string, rowNo, rowsDrawn int, cv *canvas.Canvas) int {
	//row number
	if e.p.y == rowNo {
		cv.SetFillStyle(config.Color.CurrentRowText)
	} else {
		cv.SetFillStyle(config.Color.RowText)
	}
	cv.FillText(fmt.Sprintf("%04d", rowNo+1), 14*2, float64(rowsDrawn-e.scroll+1+1)*24)

	//row style
	if strings.HasPrefix(rowBuffer, "# ") {
		cv.SetFillStyle(config.Color.H1)
	} else if strings.HasPrefix(rowBuffer, "## ") {
		cv.SetFillStyle(config.Color.H2)
	} else if strings.HasPrefix(rowBuffer, "### ") {
		cv.SetFillStyle(config.Color.H3)
	} else {
		cv.SetFillStyle(config.Color.Text)
	}

	for len(rowBuffer) > e.rowLen {
		cv.FillText(rowBuffer[:e.rowLen], 14*8, float64(rowsDrawn-e.scroll+1+1)*24)
		rowsDrawn++

		rowBuffer = rowBuffer[e.rowLen:]
	}
	cv.FillText(rowBuffer, 14*8, float64(rowsDrawn-e.scroll+1+1)*24)
	rowsDrawn++

	return rowsDrawn
}

func (e Editor) DrawPanel(cv *canvas.Canvas) {
	// panel
	cv.SetFillStyle(config.Color.EditorPanel)
	cv.FillRect(12, 12, 80*14, 34*24)

	// line no divider
	cv.SetStrokeStyle(config.Color.EditorHighlight)
	cv.BeginPath()
	cv.MoveTo(14*7, 24)
	cv.LineTo(14*7, float64(24*(1+e.Height)))
	cv.Stroke()
	cv.BeginPath()
	cv.MoveTo(14*8+14*50, 24)
	cv.LineTo(14*8+14*50, float64(24*(1+e.Height)))
	cv.Stroke()
}

func clamp(x, lo, hi int) int {
	return min(max(lo, x), hi)
}
