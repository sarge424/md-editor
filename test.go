package main

import (
	"fmt"

	"github.com/sarge424/notes/editor"
)

var (
	path = "C:/Users/arjun/Desktop/vault/crypto.md"
)

func main() {
	ed := editor.New(100)
	ed.LoadFile(path)

	fmt.Println(ed)
}
