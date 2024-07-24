package main

import (
	"log"
	"runtime"
	"time"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/nullboundary/glfont"
)

const windowWidth = 1000
const windowHeight = 600
const fontPath = "C:/users/arjun/appdata/local/microsoft/windows/fonts/firacode-regular.ttf"

func init() {
	runtime.LockOSThread()
}

func main() {

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, _ := glfw.CreateWindow(int(windowWidth), int(windowHeight), "notedit", nil, nil)

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	if err := gl.Init(); err != nil {
		panic(err)
	}

	//load font (fontfile, font scale, window width, window height
	font, err := glfont.LoadFont(fontPath, int32(16), windowWidth, windowHeight)
	if err != nil {
		log.Panicf("LoadFont: %v", err)
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		t := time.Now()

		//set color and draw text
		font.SetColor(1.0, 0.0, 1.0, 1.0)                                                   //r,g,b,a font color
		font.Printf(0, 16, 1.0, "Lorem ipsum dolor sit amet, consectetur adipiscing elit.") //x,y,scale,string,printf args
		font.Printf(0, 32, 1.0, t.Format("2006-01-02 15:04:05"))                            //x,y,scale,string,printf args

		window.SwapBuffers()
		glfw.PollEvents()

	}
}
