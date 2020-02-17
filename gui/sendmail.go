package gui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var (
	macos_mail_cmd   = []string{"open", "-a", "Mail"}
	windows_mail_cmd = []string{"start", "outlook.exe", "/a"}
	unix_mail_cmd    = []string{"xdg-open", "mailto://"}
)

func sendmail_cmd() []string {
	switch runtime.GOOS {
	case "windows":
		return windows_mail_cmd
	case "linux":
		return unix_mail_cmd
	case "darwin":
		return macos_mail_cmd
	}
	return []string{""}
}

func Sendmail(attachedFiles []string) error {
	cmds := sendmail_cmd()
	arguments := make([]string, 0)
	arguments = append(arguments, cmds[1:]...)
	arguments = append(arguments, attachedFiles...)
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
