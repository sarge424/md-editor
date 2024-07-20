package main

import (
	"strings"

	"github.com/sarge424/notes/files"

	imgui "github.com/AllenDang/cimgui-go"
	g "github.com/AllenDang/giu"
)

var (
	homeDir = "C:/Users/arjun/Desktop/vault/"

	sidebarTitle = "Files"

	currentFile   = "File.md"
	editorTitle   = "File.md"
	editorContent = "Hellooo"
	editMode      = true

	searchWord = ""
	openSearch = false

	sidebarWidth float32 = 300
	filenames    []string

	bigFont *g.FontInfo
)

func attemptFileRename() {
	if editorTitle != currentFile {
		if files.RenameFile(currentFile, editorTitle, homeDir) {
			currentFile = editorTitle
		}
	}
}

func closeSearchModal() {
	if idx := indexOf(filenames, searchWord); openSearch && idx >= 0 {
		openSearch = false
		openFile(idx)
		g.CloseCurrentPopup()
	}
}

func indexOf(s []string, e string) int {
	for idx, a := range s {
		if a == e {
			return idx
		}
	}
	return -1
}

func openFile(index int) {
	currentFile = filenames[index]
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

	var mainPane g.Layout

	//editor / viewer
	if editMode {
		editorPane := g.InputTextMultiline(&editorContent)
		editorPane.Size(g.Auto, g.Auto).OnChange(saveFile)
		mainPane = append(mainPane, editorPane)
	} else {
		for _, line := range strings.Split(editorContent, "\n") {
			viewPane := g.Label(line)
			mainPane = append(mainPane, viewPane)
		}
	}

	//main flow
	g.SingleWindow().Layout(
		g.Label("Title line").Font(bigFont),
		g.SplitLayout(g.DirectionVertical, &sidebarWidth,
			sidebarLayout,

			g.Layout{
				g.InputText(&editorTitle).Size(g.Auto).OnChange(attemptFileRename),
				mainPane,
			},
		),

		g.Custom(func() {
			if openSearch {
				searchWord = editorTitle
				g.OpenPopup("Filesearch")
			}
		}),

		g.PopupModal("Filesearch").IsOpen(&openSearch).Flags(g.WindowFlagsNoDecoration).Layout(
			g.Row(
				g.Label("Open File:"),
				g.Custom(func() {
					g.SetKeyboardFocusHere()
				}),
				g.InputText(&searchWord).
					AutoComplete(filenames).
					Hint("File to Open").
					Size(300).
					OnChange(closeSearchModal),
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
		//tabs
		g.WindowShortcut{
			Key:      g.KeyTab,
			Modifier: g.ModNone,
			Callback: addTab,
		},

		//search modal
		g.WindowShortcut{
			Key:      g.KeyO,
			Modifier: g.ModControl,
			Callback: func() { openSearch = true },
		},
		g.WindowShortcut{
			Key:      g.KeyEnter,
			Modifier: g.ModNone,
			Callback: closeSearchModal,
		},

		//sidebar resizing
		g.WindowShortcut{
			Key:      g.KeyLeftBracket,
			Modifier: g.ModControl,
			Callback: func() { sidebarWidth = max(sidebarWidth-50, 0) },
		},
		g.WindowShortcut{
			Key:      g.KeyRightBracket,
			Modifier: g.ModControl,
			Callback: func() { sidebarWidth = max(sidebarWidth+50, 0) },
		},

		//edit mode
		g.WindowShortcut{
			Key:      g.KeyE,
			Modifier: g.ModControl,
			Callback: func() { editMode = !editMode },
		},
	)
	wnd.SetTargetFPS(60)

	g.Context.IO().SetNavVisible(false)

	g.Context.IO().SetConfigFlags(imgui.ConfigFlagsNavNoCaptureKeyboard)
	g.Context.FontAtlas.SetDefaultFont("Firamono-regular.ttf", 16)
	bigFont = g.Context.FontAtlas.AddFont("Arial.ttf", 24)

	wnd.Run(loop)
}
