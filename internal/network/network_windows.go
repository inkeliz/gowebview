// +build windows

package network

import (
	"errors"
	"os/exec"
	"strings"
	"syscall"
)

var appName = `microsoft.win32webviewhost_cw5n1h2txyewy`

var (
	ErrExecCommand = errors.New("impossible to execute the command")
	ErrNotEnabled  = errors.New("it's not possible to enable localhost connections")
)

// IsAllowedPrivateConnections will return TRUE if the app is already on the whitelist.
func IsAllowedPrivateConnections() bool {
	cmd := exec.Command(`cmd`, `/c`, `CheckNetIsolation.exe`, `LoopbackExempt`, `-s`)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	result, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}

	if strings.Contains(strings.ToLower(string(result)), appName) {
		return true
	}

	return false
}

// EnablePrivateConnections will make possible to access private network (such as 127.0.0.1) from the WebView.
func EnablePrivateConnections() error {
	return setPrivateConnections(true)
}

// DisablePrivateConnections will remove the app from the whitelist. Then, make impossible to connect with private
// networks.
func DisablePrivateConnections() error {
	return setPrivateConnections(false)
}

func setPrivateConnections(enable bool) error {
	param := `-a`
	if !enable {
		param = `-d`
	}

	// That command MUST run as Administrator, so we are using powershell with `-Verb RunAs`
	// powershell.exe -Command "Start-Process cmd '/c CheckNetIsolation.exe LoopbackExempt -a -n={APP_NAME} > {PATH}'" -Verb RunAs -WindowStyle hidden"
	cmd := exec.Command(
		`powershell.exe`,
		`-Command`,
		`Start-Process cmd '/c CheckNetIsolation.exe LoopbackExempt `+param+` -n=`+appName+`' -Verb RunAs -WindowStyle hidden -Wait`,
	)

	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	if err := cmd.Run(); err != nil {
		return ErrExecCommand
	}

	if IsAllowedPrivateConnections() != enable {
		return ErrNotEnabled
	}

	return nil
}
