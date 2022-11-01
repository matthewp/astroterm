package util

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

func OpenBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}

func KillPid(pid int) error {
	proc, err := os.FindProcess(pid)
	if err == nil {
		err = proc.Signal(syscall.SIGKILL)
		if err != nil {
			return err
		}
		proc.Wait()
	}
	return nil
}

func TermPid(pid int) error {
	proc, err := os.FindProcess(pid)
	if err == nil {
		err = proc.Signal(syscall.SIGTERM)
		if err != nil {
			return err
		}
		proc.Wait()
	}
	return nil
}
