package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/sarge424/notes/editor"
	"github.com/sarge424/notes/kb"
	"github.com/tfriedel6/canvas/glfwcanvas"
)

type Shortcuter interface {
	Handle(s kb.Shortcut)
}

var (
	font = "C:/windows/fonts/liberationmono-regular.ttf"
	file = "C:/users/arjun/Desktop/vault/crypto.md"

	ed      editor.Editor
	kbState kb.State
	curr    Shortcuter
)

func main() {
	// Initialize a window
	win, cv, err := glfwcanvas.CreateWindow(1000, 800, "Canvas Example")
	win.Window.SetAttrib(glfw.Resizable, 0)
	if err != nil {
		panic(err)
	}

	//load resources
	f, err := cv.LoadFont(font)
	if err != nil {
		panic("Error loading font")
	}

	//setup editor
	ed.Pointer = editor.Cursor{X: 5, Y: 7}
	ed.LoadFile(file)

	//setup event handling
	win.KeyDown = func(scancode int, rn rune, name string) {
		if scancode == 29 || scancode == 285 {
			//control
			kbState.Ctrl++
		} else if scancode == 42 || scancode == 54 {
			//shift
			kbState.Shift++
		} else if scancode == 56 || scancode == 312 {
			//alt
			kbState.Alt++
		} else {
			//Not a modifier
			sc, ok := kbState.Emit(scancode, rn)
			if ok {
				curr.Handle(sc)
			}
		}
	}

	win.KeyUp = func(scancode int, rn rune, name string) {
		if scancode == 29 || scancode == 285 {
			//CTRL
			kbState.Ctrl--
		} else if scancode == 42 || scancode == 54 {
			//SHIFT
			kbState.Shift--
		} else if scancode == 56 || scancode == 312 {
			//ALT
			kbState.Alt--
		} else if scancode == 58 {
			//CAPS LOCK
			kbState.Caps = !kbState.Caps
		}

		// cap counters at 0
		kbState.HandleUnderflow()
	}

	//set current window
	curr = &ed

	// Main loop
	win.MainLoop(func() {
		w, h := cv.Size()
		cv.SetFillStyle("#242424")
		cv.FillRect(0, 0, float64(w), float64(h))

		cv.SetFont(f, 24)
		ed.DrawToCanvas(cv)

	})

}
