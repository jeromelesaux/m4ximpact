package main

import (
	"fmt"
	"strconv"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/jeromelesaux/m4client/m4"
)

var (
	m4url            string
	m4remoteLocation string
	selectedFiles    = make([]selectedFile, 0)
	mainwin          *ui.Window
)

type selectedFile struct {
	Name      string
	Directory string
}

func addFile(i int, m *modelBrowser) {
	if i < len(m.m4Dir.Nodes) {
		f := m.m4Dir.Nodes[i].Name
		d := m.m4Dir.CurrentPath
		fmt.Println("file " + f + " directory :" + d + " selected.")
		selectedFiles = append(selectedFiles, selectedFile{
			Name:      f,
			Directory: d,
		})
	}
}

func removeFile(i int, m *modelBrowser) {
	name := m.m4Dir.Nodes[i].Name
	directory := m.m4Dir.CurrentPath
	indexToRemove := -1
	for index, v := range selectedFiles {
		if v.Name == name && directory == v.Directory {
			indexToRemove = index
			break
		}
	}
	if indexToRemove >= 0 {
		fmt.Println("Remove selected file at index : " + strconv.Itoa(indexToRemove) + " name " + selectedFiles[indexToRemove].Name)
		copy(selectedFiles[indexToRemove:], selectedFiles[indexToRemove+1:])
	}
}

func main() {

	ui.Main(setupUI)

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

type modelBrowser struct {
	m4Dir            m4.M4Dir
	filesCheckStates []bool
	row9Text         string
	yellowRow        int
	m4client         *m4.M4Client
}

func newModelBrowser() *modelBrowser {
	m := new(modelBrowser)
	m.filesCheckStates = make([]bool, 0)
	return m
}

// nombre de fichier dans le répertoire
func (mb *modelBrowser) NumRows(m *ui.TableModel) int {
	return len(mb.m4Dir.Nodes)
}

// type des colonnes dans le tableau browser
func (mb *modelBrowser) ColumnTypes(m *ui.TableModel) []ui.TableValue {
	return []ui.TableValue{
		ui.TableString(""), // chemin du fichier
		ui.TableString(""), // nom du fichier
		ui.TableString(""), // selection du fichier pour récupération
		ui.TableString(""), // browse le répertoire bouton
	}
}

func (mb *modelBrowser) SetCellValue(m *ui.TableModel, row, column int, value ui.TableValue) {

	switch column {
	case 2:
		if !mb.m4Dir.Nodes[row].IsDirectory {
			if !mb.filesCheckStates[row] {
				mb.filesCheckStates[row] = true
				addFile(row, mb)
			} else {
				mb.filesCheckStates[row] = false
				removeFile(row, mb)
			}
			m.RowChanged(row)
		}
	case 3:
		if row < len(mb.m4Dir.Nodes) {
			if mb.m4Dir.Nodes[row].IsDirectory {
				fmt.Println("Go to navigate into directory " + m4remoteLocation + "/" + mb.m4Dir.Nodes[row].Name)
			}
		}
	}

}
func (mb *modelBrowser) CellValue(m *ui.TableModel, row, column int) ui.TableValue {
	switch column {
	case 0:
		return ui.TableString(mb.m4Dir.CurrentPath)
	case 1:
		var t ui.TableString
		if row < len(mb.m4Dir.Nodes) {
			t = ui.TableString(mb.m4Dir.Nodes[row].Name)
		} else {
			t = ui.TableString("")
		}
		return t
	case 2:
		if mb.filesCheckStates[row] {
			return ui.TableString("selected")
		} else {
			return ui.TableString("")
		}
	case 3:
		t := ui.TableString("")
		if row < len(mb.m4Dir.Nodes) {
			if mb.m4Dir.Nodes[row].IsDirectory {
				t = ui.TableString("Navigate")
			}
		}
		return t
	}

	return ui.TableString("")
}

func makeSampleBrowser() *modelBrowser {
	m := newModelBrowser()
	m.m4Dir.CurrentPath = "/home/home/documents"
	for i := 0; i < 5; i++ {
		m.m4Dir.Nodes = append(m.m4Dir.Nodes, m4.M4Node{Name: "fichier", Size: "10 ko", IsDirectory: false})
	}
	m.m4Dir.Nodes = append(m.m4Dir.Nodes, m4.M4Node{Name: "repertoire", Size: "0 ko", IsDirectory: true})
	m.filesCheckStates = make([]bool, len(m.m4Dir.Nodes))
	m4remoteLocation = m.m4Dir.CurrentPath
	return m
}

func makeM4DiskBrowser() ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	vbox.Append(hbox, false)
	currentDirectory := ui.NewEntry()
	currentDirectory.SetReadOnly(false)
	currentDirectory.SetText(m4remoteLocation)
	vbox.Append(currentDirectory, false)
	browse := ui.NewButton("Browse")
	vbox.Append(browse, false)
	/*	group := ui.NewGroup("Current remote directory")
		group.SetMargined(true)
		vbox.Append(group, true)
		group.SetChild(ui.NewNonWrappingMultilineEntry())
		entryForm := ui.NewForm()
		entryForm.SetPadded(true)
		group.SetChild(entryForm)
		entryForm.Append("Directory", ui.NewEntry(), false)
		entryForm.Append("Refresh", ui.NewButton(".Go."), true) */

	/*group2 := ui.NewGroup("Current remote directory")
	group2.SetMargined(true)
	vbox.Append(group2, true)*/
	/*entryForm2 := ui.NewForm()
	entryForm2.SetPadded(true)
	vbox.Append(entryForm2, true) */
	browser := makeSampleBrowser()
	table := ui.NewTable(&ui.TableParams{
		Model:                         ui.NewTableModel(browser),
		RowBackgroundColorModelColumn: 3,
	})
	vbox.Append(table, true)
	table.AppendTextColumn("Filepath", 0, ui.TableModelColumnNeverEditable, nil)
	table.AppendTextColumn("Filename", 1, ui.TableModelColumnNeverEditable, nil)
	table.AppendButtonColumn("Select", 2, ui.TableModelColumnAlwaysEditable)
	table.AppendButtonColumn("Browse", 3, ui.TableModelColumnAlwaysEditable)
	//entryForm2.Append("", table, false)
	//	entryForm2.Append("Refresh", ui.NewButton(".Back."), true)

	return vbox
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

	tab.Append("Browser", makeM4DiskBrowser())
	tab.SetMargined(0, true)
	tab.Append("Configuration", makeConfigurationPage())
	tab.SetMargined(1, true)
	mainwin.Show()
}
