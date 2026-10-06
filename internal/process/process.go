package process

import (
	"fmt"
	"os"
)

type Info struct {
	Command    string
	PID        int
	User       string
	Connection string
}

func PIDs(processes []Info) []int {
	pids := make([]int, 0, len(processes))

	for _, p := range processes {
		pids = append(pids, p.PID)
	}

	return pids
}

type Manager struct{}

func NewManager() Manager {
	return Manager{}
}

func (Manager) Terminate(pid int) error {
	if pid <= 0 {
		return fmt.Errorf("invalid PID: %d", pid)
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("find process %d: %w", pid, err)
	}

	if err := terminate(proc); err != nil {
		return fmt.Errorf("terminate process %d: %w", pid, err)
	}

	return nil
}
