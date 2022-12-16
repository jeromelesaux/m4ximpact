package main

import (
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/jeromelesaux/m4ximpact/gui"
)

var version = "0.1rc"

func main() {
	err := ui.Main(setupUI)
	if err != nil {
		panic(err)
	}
}

func setupUI() {
	gui.Mainwin = ui.NewWindow("M4 backup (Impact) Version "+version, 800, 500, true)
	gui.Mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		gui.Mainwin.Destroy()
		return true
	})

	tab := ui.NewTab()
	gui.Mainwin.SetChild(tab)
	gui.Mainwin.SetMargined(true)

	tab.Append("Browser", gui.MakeM4DiskBrowser())
	tab.SetMargined(0, true)
	tab.Append("Files", gui.MakeFilesTable())
	tab.SetMargined(1, true)
	tab.Append("Configuration", gui.MakeConfigurationPage())
	tab.SetMargined(2, true)
	gui.Mainwin.Show()
}
