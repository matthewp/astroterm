package actors

import (
	"os"
	"syscall"
)

func killPid(pid int) error {
	proc, err := os.FindProcess(pid)
	if err == nil {
		err = syscall.Kill(-proc.Pid, syscall.SIGKILL)
		if err != nil {
			return err
		}
		proc.Wait()
	}
	return nil
}
