package main

import (
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/jeromelesaux/m4Ximpact/gui"
)

var (
	mainwin *ui.Window
)

func main() {

	ui.Main(setupUI)

}

func setupUI() {
	mainwin = ui.NewWindow("M4xImpact", 600, 400, true)
	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})

	tab := ui.NewTab()
	mainwin.SetChild(tab)
	mainwin.SetMargined(true)

	tab.Append("Browser", gui.MakeM4DiskBrowser())
	tab.SetMargined(0, true)
	tab.Append("Files", gui.MakeFilesTable())
	tab.SetMargined(1, true)
	tab.Append("Configuration", gui.MakeConfigurationPage())
	tab.SetMargined(2, true)
	mainwin.Show()
}
