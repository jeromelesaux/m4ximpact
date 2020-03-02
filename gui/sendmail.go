package gui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/uuid"

	"github.com/jeromelesaux/eml"
	"github.com/jeromelesaux/eml/encoding"
)

var (
	macos_mail_cmd   = []string{"open", "-a", "Mail"}
	windows_mail_cmd = []string{"cmd", "/C", "start"}
	unix_mail_cmd    = []string{"xdg-open", "mailto://"}
)

func sendmail_cmd() []string {
	mailApp := config.MailerApp
	switch runtime.GOOS {
	case "windows":
		if mailApp != "" {
			cmds := windows_mail_cmd
			cmds[1] = mailApp
			return cmds
		}
		return windows_mail_cmd
	case "linux":
		return unix_mail_cmd
	case "darwin":
		if mailApp != "" {
			cmds := macos_mail_cmd
			cmds[2] = mailApp
			return cmds
		}
		return macos_mail_cmd
	}
	return []string{""}
}

func Sendmail(attachedFiles []string) error {
	cmds := sendmail_cmd()
	arguments := make([]string, 0)
	arguments = append(arguments, cmds[1:]...)
	if runtime.GOOS == "windows" {
		mailFilepath, err := CreateEml(attachedFiles)
		if err != nil {
			return err
		}
		arguments = append(arguments, mailFilepath)
	} else {
		arguments = append(arguments, attachedFiles...)
	}

	fmt.Fprintf(os.Stdout, "Executing %s with arguments %s\n", cmds[0], strings.Join(arguments, " "))
	cmd := exec.Command(cmds[0], arguments...)
	err := cmd.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while opening mail with error %v\n", err)
		return err
	}
	err = cmd.Wait()
	fmt.Fprintf(os.Stderr, "cmd mail finished with :%v\n", err)
	return nil
}

func CreateEml(attachedFiles []string) (string, error) {
	e := eml.NewEml()
	e.From = "change@me.net"
	e.To = "change@me.net"
	e.XSender = "change@me.net"
	e.XReceiver = "change@me.net"
	for _, v := range attachedFiles {
		e.AddAttachment(v)
	}
	path, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while getting current path with error %v\n", err)
		return "", err
	}
	// create temporary folder to store the mails files
	mailsFolders := filepath.Join(path, "m4backup_mails")
	_, err = os.Stat(mailsFolders)
	// create folder and sub folder if not exists
	if os.IsNotExist(err) {
		if err = os.MkdirAll(mailsFolders, os.ModePerm); err != nil {
			fmt.Fprintf(os.Stderr, "Error while creating directory %s error %v \n", mailsFolders, err)
			return "", err
		}
	}
	mailFilename := uuid.New()
	filename := filepath.Join(mailsFolders, mailFilename.String()+".eml")
	f, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create mail file %s error :%v\n", filename, err)
		return "", err
	}
	defer f.Close()
	if err := encoding.NewEncoder(f).Encode(e); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot encode mail file %s error :%v\n", filename, err)
		return "", err
	}
	return filename, nil
}
