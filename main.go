package main

import (
	"github.com/sarge424/notes/files"

	g "github.com/AllenDang/giu"
)

var (
	homeDir = "C:/Users/arjun/Desktop/vault/"

	sidebarTitle = "Files"
	taskbarTitle = "Tasks"

	editorTitle   = "File.md"
	editorContent = "Hellooo"

	searchWord = ""
	openSearch = false

	sidebarWidth float32 = 300
	editorWidth  float32 = 1000
	filenames    []string
)

func indexOf(s []string, e string) int {
	for idx, a := range s {
		if a == e {
			return idx
		}
	}
	return -1
}

func openFile(index int) {
	editorTitle = filenames[index]
	editorContent = files.LoadFile(homeDir + editorTitle)
}

func saveFile() {
	files.SaveFile(homeDir+editorTitle, editorContent)
}

func loop() {
	// sidebar
	sidebarLayout := g.Layout{
		g.Style().SetDisabled(true).To(g.InputText(&sidebarTitle).Size(g.Auto)),
	}

	for ix, name := range filenames {
		s := g.Selectable(name)
		s.OnClick(func() {
			openFile(ix)
		})
		sidebarLayout = append(sidebarLayout, s)
	}

	//taskbar
	taskbarLayout := g.Layout{
		g.Style().SetDisabled(true).To(g.InputText(&taskbarTitle).Size(g.Auto)),
	}

	//editor
	editorPane := g.InputTextMultiline(&editorContent)
	editorPane.Size(g.Auto, g.Auto).OnChange(saveFile)

	//main flow
	g.SingleWindow().Layout(
		g.SplitLayout(g.DirectionVertical, &sidebarWidth,
			sidebarLayout,

			g.SplitLayout(g.DirectionVertical, &editorWidth,
				g.Layout{
					g.Style().SetDisabled(true).To(g.InputText(&editorTitle).Size(g.Auto)),

					g.Custom(func() {
						g.SetKeyboardFocusHere()
					}),
					editorPane,
				},

				taskbarLayout,
			),
		),

		g.Custom(func() {
			if openSearch {
				g.OpenPopup("Filesearch")
			}
		}),

		g.PopupModal("Filesearch").Flags(g.WindowFlagsNoDecoration).Layout(
			g.Row(
				g.Label("Open File:"),
				g.Custom(func() {
					g.SetKeyboardFocusHere()
				}),
				g.InputText(&searchWord).
					AutoComplete(filenames).
					Hint("File to Open").
					Size(300).
					OnChange(func() {
						if idx := indexOf(filenames, searchWord); idx >= 0 {
							g.CloseCurrentPopup()
							openSearch = false
							openFile(idx)
						}
					}),
			),
		),
	)

	filenames = files.GetFilenames(homeDir)
}

func addTab() {
	for range 4 {
		g.Context.IO().AddInputCharacter(' ')
	}
}

func main() {

	filenames = files.GetFilenames(homeDir)

	openFile(0)

	wnd := g.NewMasterWindow("Notes", 400, 200, g.MasterWindowFlagsMaximized)
	wnd.RegisterKeyboardShortcuts(
		g.WindowShortcut{
			Key:      g.KeyTab,
			Modifier: g.ModNone,
			Callback: addTab,
		},
		g.WindowShortcut{
			Key:      g.KeyO,
			Modifier: g.ModControl,
			Callback: func() { openSearch = true },
		},
	)
	wnd.SetTargetFPS(60)

	g.Context.FontAtlas.SetDefaultFont("Firacode-regular.ttf", 16)

	wnd.Run(loop)
}
