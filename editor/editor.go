package editor

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type row struct {
	index  int
	length int
}

type Editor struct {
	text content
	rows []row
}

func New(chunkSize int) Editor {
	return Editor{
		text: newContent(chunkSize),
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

func (e Editor) Rows() {
	fmt.Printf("%#v\n", e.rows)
}
