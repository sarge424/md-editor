package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/sarge424/notes/editor"
	"github.com/tfriedel6/canvas/glfwcanvas"
)

var (
	font = "C:/windows/fonts/liberationmono-regular.ttf"
	home = "C:/users/arjun/Desktop/vault"
	file = "C:/users/arjun/Desktop/vault/crypto.md"
)

func main() {
	// Initialize a window
	win, cv, err := glfwcanvas.CreateWindow(1000, 800, "Canvas Example")
	win.Window.SetAttrib(glfw.Resizable, 0)

	if err != nil {
		panic(err)
	}

	f, err := cv.LoadFont(font)
	if err != nil {
		panic("Error loading font")
	}

	var ed editor.Editor
	ed.Pointer = editor.Cursor{X: 5, Y: 7}
	ed.LoadFile(file)

	// Main loop
	win.MainLoop(func() {
		w, h := cv.Size()
		cv.SetFillStyle("#242424")
		cv.FillRect(0, 0, float64(w), float64(h))

		cv.SetFont(f, 24)
		ed.Draw(cv)

	})

}
