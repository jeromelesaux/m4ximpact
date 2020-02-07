package gui

import (
	"github.com/andlabs/ui"
)

var (
	tableUi            *modelFilesTable
	tableFilesModel    *ui.TableModel
	insertingRow       = false
	updatedByFileTable = false
)

func exportFiles(b *ui.Button) {

}

func sendFilesByMail(b *ui.Button) {

}

func MakeFilesTable() ui.Control {

	tableUi = makeFilesTableUi()
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

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
		return ui.TableString(selectedFiles[row].Directory)
	case 1:
		return ui.TableString(selectedFiles[row].Name)
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
