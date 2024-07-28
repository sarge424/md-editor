package editor

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sarge424/notes/kb"
	"github.com/tfriedel6/canvas"
)

type row struct {
	index  int
	length int
}

type Editor struct {
	text   content
	rows   []row
	scroll int
}

func New(chunkSize int) Editor {
	return Editor{
		text: newContent(chunkSize),
	}
}

func (e *Editor) HandleKeystroke(k kb.Keystroke) {
	// standardize letters to uppercase
	switch k.Std() {
	case ';':
		e.scroll++
	case '\'':
		e.scroll--
	}
}

func (e *Editor) HandleShortcut(k kb.Shortcut) {
	fmt.Println(k)
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

func (e Editor) Render(cv *canvas.Canvas) {
	rowNo := 0
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

				cv.FillText(rowBuffer, 0, float64(rowNo-e.scroll+1)*24)
				//fmt.Println("<", rowBuffer, ">")

				rowBuffer = ""
				rowNo++
				if rowNo >= len(e.rows) {
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
