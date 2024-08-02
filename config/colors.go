package config

import "github.com/tfriedel6/canvas"

//https://lospec.com/palette-list/akc12

type colorConfig struct {
	AppBG string

	EditorHighlight string
	EditorPanel     string

	NavPointer  string
	EditPointer string

	RowText        string
	CurrentRowText string
	Text           string

	H1 string
	H2 string
	H3 string
}

var Color = colorConfig{
	AppBG: "#201127",

	EditorHighlight: "#2b3e44",
	EditorPanel:     "#201433",

	NavPointer:  "#152d68",
	EditPointer: "#ffffff",

	RowText:        "#355d68",
	CurrentRowText: "#94c5ac",
	Text:           "#fff",

	H1: "#c24b6e",
	H2: "#d9626b",
	H3: "#ec9a6d",
}

var (
	cv          *canvas.Canvas
	CurrentFont *canvas.Font
)

func Initialize(c *canvas.Canvas, font string) {
	f, err := c.LoadFont(font)
	if err != nil {
		panic("Error loading font")
	}

	CurrentFont = f
	cv = c

	SetFontSize(24)
}

func SetFontSize(size int) {
	cv.SetFont(CurrentFont, float64(size))
}
