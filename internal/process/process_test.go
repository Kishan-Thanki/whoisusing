package process

import "testing"

func TestPIDs(t *testing.T) {
	processes := []Info{
		{PID: 123},
		{PID: 456},
		{PID: 789},
	}

	got := PIDs(processes)
	want := []int{123, 456, 789}

	if len(got) != len(want) {
		t.Fatalf("PIDs() = %v, want %v", got, want)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Errorf("PIDs()[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

func TestPIDsEmpty(t *testing.T) {
	got := PIDs(nil)

	if len(got) != 0 {
		t.Errorf("PIDs(nil) = %v, want empty", got)
	}
}

func TestTerminateInvalidPID(t *testing.T) {
	manager := NewManager()

	for _, pid := range []int{0, -1} {
		if err := manager.Terminate(pid); err == nil {
			t.Errorf("Terminate(%d) error = nil, want error", pid)
		}
	}
}
