package gui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

var (
	macos_mail_cmd   = "open -a Mail"
	windows_mail_cmd = "start outlook.exe /a"
	unix_mail_cmd    = "xdg-open mailto://"
)

func sendmail_cmd() string {
	switch runtime.GOOS {
	case "windows":
		return windows_mail_cmd
	case "linux":
		return unix_mail_cmd
	case "darwin":
		return macos_mail_cmd
	}
	return ""
}

func Sendmail(attachedFiles []string) error {
	cmd := exec.Command(sendmail_cmd(), attachedFiles...)
	err := cmd.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while opening mail with error %v\n", err)
		return err
	}
	err = cmd.Wait()
	fmt.Fprintf(os.Stderr, "cmd mail finished with :%v\n", err)
	return nil
}
