package main

import (
	"notes/files"

	g "github.com/AllenDang/giu"
)

var (
	homeDir               = "C:/Users/arjun/Desktop/vault/"
	sidebarHeader         = "Tasks"
	editorTitle           = "File.md"
	editorContent         = "Hellooo"
	sidebarWidth  float32 = 400
)

func saveFile() {
	files.SaveFile(homeDir+editorTitle, editorContent)
}

func loop() {
	g.SingleWindow().Layout(
		g.Row(
			g.Column(
				g.Style().SetDisabled(true).To(g.InputText(&sidebarHeader).Size(sidebarWidth)),
			),
			g.Column(
				g.Style().SetDisabled(true).To(g.InputText(&editorTitle).Size(g.Auto)),
				g.InputTextMultiline(&editorContent).Size(g.Auto, g.Auto).OnChange(saveFile),
			),
		),
	)
}

func main() {

	editorTitle = files.GetFilenames(homeDir)[0]
	editorContent = files.LoadFile(homeDir + editorTitle)

	wnd := g.NewMasterWindow("Notes", 400, 200, g.MasterWindowFlagsMaximized)

	g.Context.FontAtlas.SetDefaultFont("Firacode-regular.ttf", 16)

	wnd.Run(loop)
}
