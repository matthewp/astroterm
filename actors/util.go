package actors

import (
	"os"
	"syscall"
)

func killPid(pid int) error {
	proc, err := os.FindProcess(pid)
	if err == nil {
		proc.Signal(syscall.SIGKILL)
		proc.Wait()
	}
	return err
}
