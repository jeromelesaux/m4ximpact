package gui

import (
	"fmt"
	"os"
	"strconv"

	"github.com/andlabs/ui"
	"github.com/jeromelesaux/m4client/m4"
)

var (
	m4remoteLocation string = "/"
	currentDirectory        = ui.NewEntry()
	selectedFiles           = make([]selectedFile, 0)
	m4Browser        *modelBrowser
	m4BrowserModel   *ui.TableModel
)

type FileSelectModel interface {
}

type selectedFile struct {
	Name      string
	Directory string
}

type modelBrowser struct {
	m4Dir            *m4.M4Dir
	filesCheckStates []bool
	m4client         *m4.M4Client
}

func newModelBrowser() *modelBrowser {
	m := new(modelBrowser)
	m.m4Dir = &m4.M4Dir{}
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

func (mb *modelBrowser) Clean() {
	mb.m4Dir.Nodes = mb.m4Dir.Nodes[:0]
	mb.filesCheckStates = make([]bool, len(mb.m4Dir.Nodes))
}

func (mb *modelBrowser) Navigate(m *ui.TableModel, row int) {
	rowsBefore := len(mb.m4Dir.Nodes)
	remotePath := mb.m4Dir.CurrentPath + "/" + mb.m4Dir.Nodes[row].Name
	mb.Clean()
	callM4AndUpdateBrowser(mb)
	mb.m4Dir.CurrentPath = remotePath
	rowsAfter := len(mb.m4Dir.Nodes)
	for i := 0; i < rowsBefore; i++ {
		m.RowChanged(i)
	}
	if rowsBefore > rowsAfter {
		for i := rowsAfter; i < rowsBefore; i++ {
			m.RowDeleted(i)
		}
	} else {
		for i := rowsBefore; i < rowsAfter; i++ {
			m.RowInserted(i)
		}
	}
	currentDirectory.SetText(mb.m4Dir.CurrentPath)
}

func (mb *modelBrowser) SetCellValue(m *ui.TableModel, row, column int, value ui.TableValue) {

	switch column {
	case 2:
		if !mb.m4Dir.Nodes[row].IsDirectory {
			if !mb.filesCheckStates[row] {
				mb.filesCheckStates[row] = true
				addFile(row, mb)
				insertSelectedFile(mb.m4Dir.CurrentPath, mb.m4Dir.Nodes[row].Name)
			} else {
				dir := mb.m4Dir.CurrentPath
				name := mb.m4Dir.Nodes[row].Name
				unselectFileInFilesUi(dir, name)
				mb.filesCheckStates[row] = false
				removeFile(row, mb)
			}
			m.RowChanged(row)
		}
	case 3:
		if row < len(mb.m4Dir.Nodes) {
			if mb.m4Dir.Nodes[row].IsDirectory {
				fmt.Println("Go to navigate into directory " + m4remoteLocation + "/" + mb.m4Dir.Nodes[row].Name)
				mb.Navigate(m, row)
			}
		}
	}

}

func unselectFile(directory, name string) {
	if directory == m4Browser.m4Dir.CurrentPath {
		for i, v := range m4Browser.m4Dir.Nodes {
			if v.Name == name {
				m4Browser.SetCellValue(m4BrowserModel, i, 2, ui.TableString(""))
				m4BrowserModel.RowChanged(i)
				break
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
		if mb.m4Dir.Nodes[row].IsDirectory {
			return ui.TableString("x")
		}
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
		//tableUi.CellValue(tableFilesModel, row, column)
		return t
	}

	return ui.TableString("")
}

func callM4AndUpdateBrowser(m *modelBrowser) {
	err, dir := m.m4client.GetDir(m.m4Dir.CurrentPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while calling m4 (%s) error : %v\n", m.m4client.IPClient, err)
		ui.MsgBoxError(Mainwin, "Error while calling M4",
			"Error while calling M4 "+m.m4client.IPClient+", error : "+err.Error())
		return
	}
	m.m4Dir = dir
	//m.m4Dir.CurrentPath = "/"
	m.filesCheckStates = make([]bool, len(m.m4Dir.Nodes))
	m4remoteLocation = m.m4Dir.CurrentPath
}

func makeDefaultBrowser() *modelBrowser {
	m := newModelBrowser()
	m.m4client = &m4.M4Client{}
	m.m4Dir.CurrentPath = "/"
	m.filesCheckStates = make([]bool, len(m.m4Dir.Nodes))
	m4remoteLocation = m.m4Dir.CurrentPath
	return m
}

func makeSampleBrowser() *modelBrowser {
	m := newModelBrowser()
	m.m4Dir.CurrentPath = "/home/home/documents"
	for i := 0; i < 5; i++ {
		m.m4Dir.Nodes = append(m.m4Dir.Nodes, m4.M4Node{Name: "fichier" + strconv.Itoa(i), Size: "10 ko", IsDirectory: false})
	}
	m.m4Dir.Nodes = append(m.m4Dir.Nodes, m4.M4Node{Name: "repertoire", Size: "0 ko", IsDirectory: true})
	m.filesCheckStates = make([]bool, len(m.m4Dir.Nodes))
	m4remoteLocation = m.m4Dir.CurrentPath
	return m
}

func updateSampleBrowser(m *modelBrowser) {
	m.m4Dir.CurrentPath = "/home/home/documents/repertoire"
	m.m4Dir.Nodes = m.m4Dir.Nodes[:0]
	for i := 0; i < 15; i++ {
		m.m4Dir.Nodes = append(m.m4Dir.Nodes, m4.M4Node{Name: "fichier" + strconv.Itoa(i), Size: "10 ko", IsDirectory: false})
	}
	m.m4Dir.Nodes = append(m.m4Dir.Nodes, m4.M4Node{Name: "repertoire2", Size: "0 ko", IsDirectory: true})
	m.m4Dir.Nodes = append(m.m4Dir.Nodes, m4.M4Node{Name: "repertoire3", Size: "0 ko", IsDirectory: true})
	m.filesCheckStates = make([]bool, len(m.m4Dir.Nodes))
	m4remoteLocation = m.m4Dir.CurrentPath
}

func MakeM4DiskBrowser() ui.Control {

	m4Browser = makeDefaultBrowser()

	m4Browser.m4client.IPClient = config.M4Url

	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	//
	//	hbox := ui.NewHorizontalBox()
	//	hbox.SetPadded(true)
	//	vbox.Append(hbox, false)
	//
	grid := ui.NewGrid()
	grid.SetPadded(true)
	vbox.Append(grid, false)

	//	hbox.Append(currentDirectory, false)
	currentDirectory.SetReadOnly(false)
	currentDirectory.SetText(m4remoteLocation)
	//vbox.Append(currentDirectory, false)
	browse := ui.NewButton("Browse")
	browse.OnClicked(browseM4)
	//	hbox.Append(browse, false)
	goBackButton := ui.NewButton("Go back")
	goBackButton.OnClicked(goBack)

	grid.Append(browse,
		0, 1, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)
	grid.Append(currentDirectory,
		1, 1, 1, 1,
		true, ui.AlignFill, false, ui.AlignFill)
	vbox.Append(goBackButton, false)
	m4BrowserModel = ui.NewTableModel(m4Browser)
	table := ui.NewTable(&ui.TableParams{
		Model:                         m4BrowserModel,
		RowBackgroundColorModelColumn: 3,
	})
	vbox.Append(table, true)
	table.AppendTextColumn("Filepath", 0, ui.TableModelColumnNeverEditable, nil)
	table.AppendTextColumn("Filename", 1, ui.TableModelColumnNeverEditable, nil)
	table.AppendButtonColumn("Select", 2, ui.TableModelColumnAlwaysEditable)
	table.AppendButtonColumn("Browse", 3, ui.TableModelColumnAlwaysEditable)

	return vbox
}

func addFile(i int, m *modelBrowser) {
	if i < len(m.m4Dir.Nodes) {
		f := m.m4Dir.Nodes[i].Name
		d := m.m4Dir.CurrentPath
		isPresent := false
		for _, v := range selectedFiles {
			if v.Name == f && v.Directory == d {
				isPresent = true
				break
			}
		}
		if !isPresent {
			fmt.Println("file " + f + " directory :" + d + " selected.")
			selectedFiles = append(selectedFiles, selectedFile{
				Name:      f,
				Directory: d,
			})

		}
	}
}

func removeFileWithData(directory, name string) {
	indexToRemove := -1
	for index, v := range selectedFiles {
		if v.Name == name && directory == v.Directory {
			indexToRemove = index
			break
		}
	}
	if indexToRemove >= 0 {
		fmt.Println("Remove selected file at index : " + strconv.Itoa(indexToRemove) + " name " + selectedFiles[indexToRemove].Name)
		selectedFiles = append(selectedFiles[:indexToRemove], selectedFiles[indexToRemove+1:]...)
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
		selectedFiles = append(selectedFiles[:indexToRemove], selectedFiles[indexToRemove+1:]...)
	}
}

func browseM4(*ui.Button) {
	callM4AndUpdateBrowser(m4Browser)
}

func goBack(*ui.Button) {

}
