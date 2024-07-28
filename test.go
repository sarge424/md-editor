package main

import (
	"fmt"

	"github.com/sarge424/notes/editor"
)

func main() {
	tx := editor.NewBody(10)

	tx.Insert("Hello", 0)
	tx.Insert("World", 5)
	tx.Insert("World2", 10)
	tx.Insert("qwer", 10)
	//tx.Delete(9, 4)
	fmt.Println(tx)
}
