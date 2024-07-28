package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/sarge424/notes/editor"
	"github.com/sarge424/notes/kb"
	"github.com/tfriedel6/canvas/glfwcanvas"
)

var (
	font = "C:/windows/fonts/liberationmono-regular.ttf"
	file = "C:/users/arjun/Desktop/vault/crypto.md"

	kbState kb.State
	curr    kb.Keyboarder

	mouseX, mouseY int
	moused         bool
)

func main() {
	// Initialize a window
	win, cv, err := glfwcanvas.CreateWindow(1000, 800, "Noter")
	win.Window.SetAttrib(glfw.Resizable, 1)
	if err != nil {
		panic(err)
	}

	//load resources
	f, err := cv.LoadFont(font)
	if err != nil {
		panic("Error loading font")
	}

	//setup event handling
	// win.KeyDown = func(scancode int, rn rune, name string) {
	// 	if scancode == 29 || scancode == 285 {
	// 		//control
	// 		kbState.Ctrl++
	// 	} else if scancode == 42 || scancode == 54 {
	// 		//shift
	// 		kbState.Shift++
	// 	} else if scancode == 56 || scancode == 312 {
	// 		//alt
	// 		kbState.Alt++
	// 	} else {
	// 		notAlphaNum := rn == 0 || rn == '\n' || rn == '\t'
	// 		modPressed := kbState.Ctrl > 0 || kbState.Alt > 0

	// 		//create a SC if mods are held or special keys are pressed
	// 		if notAlphaNum || modPressed {
	// 			sc := kbState.Emit(scancode, rn)
	// 			curr.HandleShortcut(sc)
	// 		}
	// 	}
	// }

	// win.KeyChar = func(rn rune) {
	// 	ks := kb.Keystroke(rn)
	// 	curr.HandleKeystroke(ks)
	// }

	// win.KeyUp = func(scancode int, rn rune, name string) {
	// 	if scancode == 29 || scancode == 285 {
	// 		//CTRL
	// 		kbState.Ctrl--
	// 	} else if scancode == 42 || scancode == 54 {
	// 		//SHIFT
	// 		kbState.Shift--
	// 	} else if scancode == 56 || scancode == 312 {
	// 		//ALT
	// 		kbState.Alt--
	// 	}

	// 	// cap counters at 0
	// 	kbState.HandleUnderflow()
	// }

	//mouse events
	win.MouseDown = func(button, x, y int) {
		if button != 0 {
			return
		}
		mouseX = x
		mouseY = y
		moused = true
	}

	win.MouseMove = func(x, y int) {
		//move window
		if moused {
			wx, wy := win.Window.GetPos()
			dx := x - mouseX
			dy := y - mouseY
			win.Window.SetPos(wx+dx, wy+dy)
		}
	}

	win.MouseUp = func(button, x, y int) {
		moused = false
	}

	//set context
	ed := editor.New(100)
	ed.LoadFile(file)

	// Main loop
	win.MainLoop(func() {

		w, h := cv.Size()
		cv.SetFillStyle("#242424")
		cv.FillRect(0, 0, float64(w), float64(h))

		cv.SetFont(f, 24)
		cv.SetFillStyle("#FFF")
		ed.Render(cv)

		win.Close()
	})

}
