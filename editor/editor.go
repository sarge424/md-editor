package editor

import (
	"log"
	"os"
	"strings"
)

type Editor struct {
	text content
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

	// load the file as chunks
	for len(data) >= e.text.chunkSize {
		e.text.append(data[:e.text.chunkSize])
		data = data[e.text.chunkSize:]
	}

	// load leftover characters
	if len(data) > 0 {
		e.text.append(data)
	}
}

func (e Editor) String() string {
	return e.text.String()
}
