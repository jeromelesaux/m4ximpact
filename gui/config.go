package gui

import (
	"fmt"

	"github.com/andlabs/ui"
	"github.com/jeromelesaux/m4Ximpact/common"
)

var (
	config = common.NewConfig()
)

func MakeConfigurationPage() ui.Control {
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
	m4urlEntry.SetText(config.M4Url)
	confForm.Append("M4 Address", m4urlEntry, false)

	// m4 button to set url
	m4UrlButton := ui.NewButton(".Save Url.")
	m4UrlButton.OnClicked(func(*ui.Button) {
		config.M4Url = m4urlEntry.Text()
		config.Save()
		fmt.Println("Set m4 url : " + config.M4Url)
	})

	confForm.Append("Set", m4UrlButton, false)
	return vbox
}
