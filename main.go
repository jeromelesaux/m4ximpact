package main

import (
	"fmt"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
)

var (
	m4url string

	mainwin *ui.Window
)

func main() {

	ui.Main(setupUi)

}

func makeConfigurationPage() ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	vbox.Append(hbox, false)

	m4configuration := ui.NewGroup("M4 address configuration")
	m4configuration.SetMargined(true)
	vbox.Append(m4configuration, true)

	//m4configuration.SetChild(ui.NewNonWrappingMultilineEntry())

	confForm := ui.NewForm()
	confForm.SetPadded(true)
	m4configuration.SetChild(confForm)

	m4urlEntry := ui.NewEntry()
	m4urlEntry.SetReadOnly(false)
	confForm.Append("M4 Address", m4urlEntry, false)

	// m4 button to set url
	m4UrlButton := ui.NewButton(".Save Url.")
	m4UrlButton.OnClicked(func(*ui.Button) {
		m4url = m4urlEntry.Text()
		fmt.Println("Set m4 url : " + m4url)
	})

	confForm.Append("Set", m4UrlButton, false)
	return vbox
}

type modelHandler struct {
	row9Text    string
	yellowRow   int
	checkStates [15]int
}

func newModelHandler() *modelHandler {
	m := new(modelHandler)
	m.row9Text = "You can edit this one"
	m.yellowRow = -1
	return m
}

func (mh *modelHandler) CellValue(m *ui.TableModel, row, column int) ui.TableValue {
	return ui.TableString("")
}

func (mh *modelHandler) ColumnTypes(m *ui.TableModel) []ui.TableValue {
	return []ui.TableValue{
		ui.TableString(""),
	}
}

func (mh *modelHandler) NumRows(m *ui.TableModel) int {
	return 15
}

func (mh *modelHandler) SetCellValue(m *ui.TableModel, row, column int, value ui.TableValue) {
}

func makeM4DiskBrowser() ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	vbox.Append(hbox, false)
	group := ui.NewGroup("Current remote directory")
	group.SetMargined(true)
	vbox.Append(group, true)
	group.SetChild(ui.NewNonWrappingMultilineEntry())
	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	group.SetChild(entryForm)
	entryForm.Append("Directory", ui.NewEntry(), false)
	entryForm.Append("Refresh", ui.NewButton(".Go."), true)

	group2 := ui.NewGroup("Current remote directory")
	group2.SetMargined(true)
	vbox.Append(group2, true)
	entryForm2 := ui.NewForm()
	entryForm2.SetPadded(true)
	group2.SetChild(entryForm2)
	entryForm2.Append("Browser", ui.NewTable(&ui.TableParams{
		Model:                         ui.NewTableModel(newModelHandler()),
		RowBackgroundColorModelColumn: 3,
	}), false)
	entryForm2.Append("Refresh", ui.NewButton(".Back."), true)
	return vbox
}

func setupUi() {
	mainwin = ui.NewWindow("M4xImpact", 300, 300, true)
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

	tab.Append("Browser", makeM4DiskBrowser())
	tab.SetMargined(0, true)
	tab.Append("Configuration", makeConfigurationPage())
	tab.SetMargined(1, true)
	mainwin.Show()
}
