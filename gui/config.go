package gui

import (
	"fmt"

	"github.com/andlabs/ui"
	"github.com/jeromelesaux/m4ximpact/common"
)

var (
	config      = common.NewConfig()
	m4urlEntry  *ui.Entry
	mailerEntry *ui.Entry
	mailFrom    *ui.Entry
)

// nolint: ireturn
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

	m4urlEntry = ui.NewEntry()
	m4urlEntry.SetReadOnly(false)
	m4urlEntry.SetText(config.M4Url)
	confForm.Append("M4 Address", m4urlEntry, false)

	mailerEntry = ui.NewEntry()
	mailerEntry.SetReadOnly(false)
	mailerEntry.SetText(config.MailerApp)
	confForm.Append("Mail application", mailerEntry, false)

	mailFrom = ui.NewEntry()
	mailFrom.SetReadOnly(false)
	mailFrom.SetText(config.MailFrom)
	confForm.Append("Your mail:", mailFrom, false)
	// m4 button to set url
	saveConfigButton := ui.NewButton(".Save Configuration.")
	saveConfigButton.OnClicked(saveConfiguration)
	confForm.Append("Set", saveConfigButton, false)

	return vbox
}

func saveConfiguration(b *ui.Button) {
	config.M4Url = m4urlEntry.Text()
	m4Browser.m4client.IPClient = config.M4Url
	config.MailerApp = mailerEntry.Text()
	config.MailFrom = mailFrom.Text()
	m4urlEntry2.SetText(config.M4Url)
	err := config.Save()
	if err != nil {
		fmt.Printf("error while saving %v\n", err)
	}
	fmt.Println("Set m4 url : " + config.M4Url + " , and MailerApp : " + config.MailerApp + " , and your mail :" + config.MailFrom)
}
