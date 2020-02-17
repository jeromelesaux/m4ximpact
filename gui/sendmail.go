package gui

import "runtime"

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
