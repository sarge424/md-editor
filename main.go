package main

import (
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

	searchWord = ""
	openSearch = false

	sidebarWidth float32 = 300
	filenames    []string
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

	//editor
	editorPane := g.InputTextMultiline(&editorContent)
	editorPane.Size(g.Auto, g.Auto).OnChange(saveFile)

	//main flow
	g.SingleWindow().Layout(
		g.SplitLayout(g.DirectionVertical, &sidebarWidth,
			sidebarLayout,

			g.Layout{
				g.InputText(&editorTitle).Size(g.Auto).OnChange(attemptFileRename),

				editorPane,
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
		g.WindowShortcut{
			Key:      g.KeyEnter,
			Modifier: g.ModNone,
			Callback: closeSearchModal,
		},
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
	)
	wnd.SetTargetFPS(60)

	g.Context.IO().SetConfigFlags(imgui.ConfigFlagsNavNoCaptureKeyboard)
	g.Context.FontAtlas.SetDefaultFont("Firacode-regular.ttf", 16)

	wnd.Run(loop)
}
