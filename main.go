package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/sarge424/notes/internal/editor"
	"github.com/tfriedel6/canvas/glfwcanvas"
)

var (
	fontpath = "C:/windows/fonts/liberationmono-regular.ttf"
	filepath = "C:\\Users\\arjun\\Desktop\\vault\\test.md"
)

func main() {
	// WINDOW INIT
	glfw.Init()
	//glfw.WindowHint(glfw.Decorated, glfw.False)
	glfw.WindowHint(glfw.Floating, glfw.True)
	glfw.WindowHint(glfw.TransparentFramebuffer, glfw.True)

	win, cv, err := glfwcanvas.CreateWindow(14*80+24, 24*35, "mdedit")
	if err != nil {
		panic(err)
	}

	cv.SetFont(fontpath, 20)

	ed := editor.New()
	ed.LoadFile(filepath)

	// MAIN LOOP
	win.MainLoop(func() {

		w, h := cv.Size()
		editor.DrawAppBG(w, h, 10, cv)

		cv.SetFillStyle("#FFF")
		cv.SetLineWidth(2.5)
		ed.DrawBuffer(cv)
	})

}
