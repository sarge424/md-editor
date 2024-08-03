package main

import (
	"fmt"
	"image"
	"os"
	"slices"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/sarge424/notes/config"
	"github.com/sarge424/notes/editor"
	"github.com/sarge424/notes/kb"
	"github.com/tfriedel6/canvas/glfwcanvas"
)

var (
	font = "C:/windows/fonts/liberationmono-regular.ttf"
	file = "C:/users/arjun/Desktop/vault/test.md"

	kbState       kb.State
	curr          kb.Keyboarder
	liftChans     = make(map[int]chan interface{})
	heldShortcuts = []string{
		"BCKSP",
		"SHF BCKSP",
		"ENTER",
		"SHF ENTER",
		"328",
		"331",
		"333",
		"336",
	}

	mouseX, mouseY int
	moused         bool
)

func holdKey(code int, signal <-chan interface{}) {
	ts := time.Now()
	t := time.Now()

	for {
		select {
		case <-signal:
			return
		default:
			if time.Since(t) < time.Millisecond*50 || time.Since(ts) < time.Millisecond*500 {
				continue
			}
			curr.HandleShortcut(kbState.Emit(code, 0))
			t = time.Now()
		}
	}
}

func main() {
	// WINDOW INIT
	win, cv, err := glfwcanvas.CreateWindow(14*80+24, 24*35, "mdedit")
	win.Window.SetAttrib(glfw.Resizable, 0)
	win.Window.SetAttrib(glfw.Decorated, 0)
	if err != nil {
		panic(err)
	}

	//LOAD RESOURCES
	config.Initialize(cv, font)

	iconFile, _ := os.Open("icon.png")
	defer iconFile.Close()
	icon, _, _ := image.Decode(iconFile)
	win.Window.SetIcon([]image.Image{icon})

	//EVENT HANDLING
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
			notAlphaNum := rn == 0 || rn == '\n' || rn == '\t'
			modPressed := kbState.Ctrl > 0 || kbState.Alt > 0

			//create a SC if mods are held or special keys are pressed
			if notAlphaNum || modPressed {
				sc := kbState.Emit(scancode, rn)
				curr.HandleShortcut(sc)

				if slices.Contains(heldShortcuts, sc.String()) {
					if _, ok := liftChans[scancode]; !ok {
						liftChans[scancode] = make(chan interface{})
						fmt.Println("made hold channel for ", scancode)
					}
					go holdKey(scancode, liftChans[scancode])
				}
			}
		}
	}

	win.KeyChar = func(rn rune) {
		if rn == '\n' {
			return
		}
		ks := kb.Keystroke(rn)
		curr.HandleKeystroke(ks)
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
		} else if slices.Contains(heldShortcuts, kbState.Emit(scancode, rn).String()) {
			liftChans[scancode] <- 0
		}

		// cap counters at 0
		kbState.HandleUnderflow()
	}

	//MOUSE EVENTS
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

	//EDITOR SETUP
	ed := editor.NewEditor(50, 33, 100)
	ed.LoadFile(file)

	curr = &ed

	// MAIN LOOP
	win.MainLoop(func() {

		w, h := cv.Size()
		cv.SetFillStyle(config.Color.AppBG)
		cv.FillRect(0, 0, float64(w), float64(h))

		ed.Render(cv)
	})

}
