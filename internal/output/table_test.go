package output

import (
	"strings"
	"testing"

	"github.com/kishan-thanki/whoisusing/internal/process"
)

func TestTable(t *testing.T) {
	var buf strings.Builder

	processes := []process.Info{
		{
			Command:    "node",
			PID:        123,
			User:       "alice",
			Connection: "TCP *:8080 (LISTEN)",
		},
		{
			Command:    "python3",
			PID:        456,
			User:       "bob",
			Connection: "TCP *:8080 (LISTEN)",
		},
	}

	if err := Table(&buf, processes); err != nil {
		t.Fatalf("Table() error = %v", err)
	}

	got := buf.String()

	for _, want := range []string{
		"COMMAND",
		"PID",
		"USER",
		"CONNECTION STATE",
		"node",
		"123",
		"alice",
		"python3",
		"456",
		"bob",
		"TCP *:8080 (LISTEN)",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("Table() output does not contain %q:\n%s", want, got)
		}
	}
}

func TestTableEmpty(t *testing.T) {
	var buf strings.Builder

	if err := Table(&buf, nil); err != nil {
		t.Fatalf("Table() error = %v", err)
	}

	got := buf.String()

	if !strings.Contains(got, "COMMAND") {
		t.Errorf("missing header: %q", got)
	}
}

func TestPIDs(t *testing.T) {
	var buf strings.Builder

	err := PIDs(&buf, []int{123, 456, 789})
	if err != nil {
		t.Fatalf("PIDs() error = %v", err)
	}

	want := "123\n456\n789\n"
	if got := buf.String(); got != want {
		t.Errorf("PIDs() = %q, want %q", got, want)
	}
}

func TestPIDsEmpty(t *testing.T) {
	var buf strings.Builder

	if err := PIDs(&buf, nil); err != nil {
		t.Fatalf("PIDs() error = %v", err)
	}

	if buf.Len() != 0 {
		t.Errorf("PIDs() output = %q, want empty", buf.String())
	}
}

func TestFree(t *testing.T) {
	var buf strings.Builder

	if err := Free(&buf, 8080); err != nil {
		t.Fatalf("Free() error = %v", err)
	}

	want := "Port 8080 is currently free.\n"
	if got := buf.String(); got != want {
		t.Errorf("Free() = %q, want %q", got, want)
	}
}

func TestFound(t *testing.T) {
	var buf strings.Builder

	if err := Found(&buf, 8080); err != nil {
		t.Fatalf("Found() error = %v", err)
	}

	want := "Processes using port 8080:\n\n"
	if got := buf.String(); got != want {
		t.Errorf("Found() = %q, want %q", got, want)
	}
}
