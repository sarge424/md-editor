package main

import (
	"notes/files"

	g "github.com/AllenDang/giu"
)

var (
	homeDir               = "C:/Users/arjun/Desktop/vault/"
	sidebarHeader         = "Files"
	editorTitle           = "File.md"
	taskbarTitle          = "Tasks"
	editorContent         = "Hellooo"
	sidebarWidth  float32 = 300
	editorWidth   float32 = 1000
	filenames     []string
)

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
		g.Style().SetDisabled(true).To(g.InputText(&sidebarHeader).Size(g.Auto)),
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

					editorPane,
				},

				taskbarLayout,
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
	wnd.RegisterKeyboardShortcuts(g.WindowShortcut{
		Key:      g.KeyTab,
		Modifier: g.ModNone,
		Callback: addTab,
	})

	g.Context.FontAtlas.SetDefaultFont("Firacode-regular.ttf", 16)

	wnd.Run(loop)
}
