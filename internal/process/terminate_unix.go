//go:build unix

package process

import (
	"os"
	"syscall"
)

func terminate(proc *os.Process) error {
	return proc.Signal(syscall.SIGTERM)
}
