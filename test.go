package main

import (
	"github.com/sarge424/notes/editor"
)

var (
	path = "C:/Users/arjun/Desktop/vault/test.md"
)

func main() {
	ed := editor.New(100)
	ed.LoadFile(path)

	ed.Rows()
}
