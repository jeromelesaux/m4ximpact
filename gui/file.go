package gui

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/andlabs/ui"
)

var (
	tableUi            *modelFilesTable
	tableFilesModel    *ui.TableModel
	insertingRow       = false
	updatedByFileTable = false
)

// export files to local harddrive
func exportFiles(b *ui.Button) {
	downloadFiles()
}

func completeM4Backup(b *ui.Button) {
	path, err := os.Getwd()
	if err != nil {
		ui.MsgBoxError(Mainwin, "Error in folder",
			err.Error())
		return
	}
	// make root directory
	t := time.Now()
	folderName := t.Format("2006-01-02")
	fmt.Fprintf(os.Stdout, "Creating folder %s\n", folderName)
	rootpath := filepath.Join(path, folderName)
	m4BackupFolder("/", rootpath)
}

func m4BackupFolder(remotefolder, localfolder string) {
	// create local folder
	if err := os.MkdirAll(localfolder, os.ModePerm); err != nil {
		ui.MsgBoxError(Mainwin, "Error in folder creation",
			err.Error())
		return
	}
	// call remote m4 to get the remotefolder content
	err, dir := m4Browser.m4client.GetDir(remotefolder)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while calling m4 (%s) error : %v\n", m4Browser.m4client.IPClient, err)
		ui.MsgBoxError(Mainwin, "Error while calling M4",
			"Error while calling M4 "+m4Browser.m4client.IPClient+", error : "+err.Error())
		return
	}
	for i := 0; i < len(dir.Nodes); i++ {
		node := dir.Nodes[i]
		if node.IsDirectory {
			// directory get content folder
			m4BackupFolder(remotefolder+"/"+node.Name, filepath.Join(localfolder, node.Name))
		} else {
			// save file in localfolder
			folder := dir.CurrentPath
			filename := node.Name
			fmt.Fprintf(os.Stdout, "folder %s file %s will be donwloaded.\n", folder, filename)
			content, err := m4Browser.m4client.DownloadContent(folder + "/" + filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while getting file (%s/%s) error : %v\n", folder, filename, err)
				continue
			}
			folderFilename := filepath.Join(localfolder, folder)
			_, err = os.Stat(folderFilename)
			// create folder and sub folder if not exists
			if os.IsNotExist(err) {
				if err = os.MkdirAll(folderFilename, os.ModePerm); err != nil {
					fmt.Fprintf(os.Stderr, "Error while creating directory %s error %v \n", folderFilename, err)
					continue
				}
			}
			// copy file locally
			if err = ioutil.WriteFile(filepath.Join(folderFilename, filename), content, os.ModePerm); err != nil {
				fmt.Fprintf(os.Stderr, "Error while creating file %s error %v \n", filename, err)
			}
		}
	}
}

func downloadFiles() {
	path, err := os.Getwd()
	if err != nil {
		ui.MsgBoxError(Mainwin, "Error in folder",
			err.Error())
		return
	}
	// make root directory
	t := time.Now()
	folderName := t.Format("2006-01-02")
	fmt.Fprintf(os.Stdout, "Creating folder %s\n", folderName)
	rootpath := filepath.Join(path, folderName)
	if err := os.MkdirAll(rootpath, os.ModePerm); err != nil {
		ui.MsgBoxError(Mainwin, "Error in folder creation",
			err.Error())
		return
	}
	onError := false
	// download all selected files
	for i := 0; i < tableUi.NumRows(tableFilesModel); i++ {
		folder := string(tableUi.CellValue(tableFilesModel, i, 0).(ui.TableString))
		filename := string(tableUi.CellValue(tableFilesModel, i, 1).(ui.TableString))
		fmt.Fprintf(os.Stdout, "folder %s file %s will be donwloaded.\n", folder, filename)
		content, err := m4Browser.m4client.DownloadContent(folder + "/" + filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while getting file (%s/%s) error : %v\n", folder, filename, err)
			onError = true
			continue
		}
		folderFilename := filepath.Join(rootpath, folder)
		_, err = os.Stat(folderFilename)
		// create folder and sub folder if not exists
		if os.IsNotExist(err) {
			if err = os.MkdirAll(folderFilename, os.ModePerm); err != nil {
				fmt.Fprintf(os.Stderr, "Error while creating directory %s error %v \n", folderFilename, err)
				onError = true
				continue
			}
		}
		// copy file locally
		if err = ioutil.WriteFile(filepath.Join(folderFilename, filename), content, os.ModePerm); err != nil {
			fmt.Fprintf(os.Stderr, "Error while creating file %s error %v \n", filename, err)
			onError = true
		}
	}
	if onError {
		ui.MsgBoxError(Mainwin, "Download Error !",
			"Errors occur when downloading files, check log to know why.")
	}
}

func files() []string {
	filespaths := make([]string, 0)

	path, err := os.Getwd()
	if err != nil {
		ui.MsgBoxError(Mainwin, "Error in folder",
			err.Error())
		return filespaths
	}
	t := time.Now()
	folderName := t.Format("2006-01-02")
	fmt.Fprintf(os.Stdout, "Creating folder %s\n", folderName)
	rootpath := filepath.Join(path, folderName)
	if err := os.MkdirAll(rootpath, os.ModePerm); err != nil {
		ui.MsgBoxError(Mainwin, "Error in folder creation",
			err.Error())
		return filespaths
	}
	// download all selected files
	for i := 0; i < tableUi.NumRows(tableFilesModel); i++ {
		folder := string(tableUi.CellValue(tableFilesModel, i, 0).(ui.TableString))
		filename := string(tableUi.CellValue(tableFilesModel, i, 1).(ui.TableString))
		folderFilename := filepath.Join(rootpath, folder)
		localFilepath := filepath.Join(folderFilename, filename)
		filespaths = append(filespaths, localFilepath)
	}
	return filespaths
}

func sendFilesByMail(b *ui.Button) {
	downloadFiles()
	filespaths := files()
	Sendmail(filespaths)
}

func MakeFilesTable() ui.Control {

	tableUi = makeFilesTableUi()
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	backupAll := ui.NewButton("Backup all your M4 content")
	backupAll.OnClicked(completeM4Backup)
	vbox.Append(backupAll, false)
	grid := ui.NewGrid()
	grid.SetPadded(true)
	vbox.Append(grid, false)

	export := ui.NewButton("Save")
	export.OnClicked(exportFiles)
	//vbox.Append(currentDirectory, false)
	sendByMail := ui.NewButton("Send by Mail")
	sendByMail.OnClicked(sendFilesByMail)
	//	hbox.Append(browse, false)
	grid.Append(export,
		0, 1, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)
	grid.Append(sendByMail,
		1, 1, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)

	tableFilesModel = ui.NewTableModel(tableUi)
	table := ui.NewTable(&ui.TableParams{
		Model:                         tableFilesModel,
		RowBackgroundColorModelColumn: 3,
	})
	vbox.Append(table, true)
	table.AppendTextColumn("Filepath", 0, ui.TableModelColumnNeverEditable, nil)
	table.AppendTextColumn("Filename", 1, ui.TableModelColumnNeverEditable, nil)
	table.AppendButtonColumn("remove", 2, ui.TableModelColumnAlwaysEditable)
	return vbox
}

func makeFilesTableUi() *modelFilesTable {
	m := new(modelFilesTable)
	return m
}

type modelFilesTable struct {
}

// type des colonnes dans le tableau browser
func (mb *modelFilesTable) ColumnTypes(m *ui.TableModel) []ui.TableValue {
	return []ui.TableValue{
		ui.TableString(""), // chemin du fichier
		ui.TableString(""), // nom du fichier
		ui.TableString(""), // selection du fichier pour récupération
	}
}

// nombre de fichier dans le répertoire
func (mb *modelFilesTable) NumRows(m *ui.TableModel) int {
	return len(selectedFiles)
}

func (mb *modelFilesTable) CellValue(m *ui.TableModel, row, column int) ui.TableValue {
	switch column {
	case 0:
		if row < len(selectedFiles) {
			return ui.TableString(selectedFiles[row].Directory)
		}
		return ui.TableString("")
	case 1:
		if row < len(selectedFiles) {
			return ui.TableString(selectedFiles[row].Name)
		}
		return ui.TableString("")
	case 2:
		return ui.TableString("remove")
	}

	return ui.TableString("")
}

func (mb *modelFilesTable) SetCellValue(m *ui.TableModel, row, column int, value ui.TableValue) {

	switch column {
	case 2:
		if !insertingRow && !updatedByFileTable {
			path := selectedFiles[row].Directory
			file := selectedFiles[row].Name
			unselectFile(path, file)
			//		removeFileWithData(path, file)
			//	m.RowDeleted(row)
		}
	}

}

func insertSelectedFile(path, name string) {
	insertingRow = true
	index := len(selectedFiles) - 1
	tableUi.SetCellValue(tableFilesModel, index, 0, ui.TableString(path))
	tableUi.SetCellValue(tableFilesModel, index, 1, ui.TableString(name))
	tableUi.SetCellValue(tableFilesModel, index, 2, ui.TableString(""))
	tableFilesModel.RowInserted(index)
	insertingRow = false
}

func unselectFileInFilesUi(path, name string) {
	if path == m4Browser.m4Dir.CurrentPath {
		for i, v := range selectedFiles {
			if v.Name == name {
				tableFilesModel.RowDeleted(i)
				break
			}
		}
	}
}
