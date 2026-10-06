package app

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/kishan-thanki/whoisusing/internal/process"
)

type fakeFinder struct {
	results [][]process.Info
	errs    []error
	calls   int
}

func (f *fakeFinder) Find(int) ([]process.Info, error) {
	index := f.calls
	f.calls++

	if len(f.results) == 0 {
		return nil, nil
	}

	if index >= len(f.results) {
		index = len(f.results) - 1
	}

	if len(f.errs) > 0 && index < len(f.errs) && f.errs[index] != nil {
		return nil, f.errs[index]
	}

	return f.results[index], nil
}

type fakeTerminator struct {
	pids []int
	err  error
}

func (f *fakeTerminator) Terminate(pid int) error {
	f.pids = append(f.pids, pid)
	return f.err
}

type failingWriter struct {
	err error
}

func (w failingWriter) Write([]byte) (int, error) {
	return 0, w.err
}

func TestRunHelp(t *testing.T) {
	runner := New(
		&fakeFinder{},
		&fakeTerminator{},
		nil,
	)

	var stdout, stderr strings.Builder

	code := runner.Run(
		[]string{"--help"},
		&stdout,
		&stderr,
	)

	if code != 0 {
		t.Fatalf("Run() = %d, want 0", code)
	}

	if !strings.Contains(stdout.String(), "Usage: whoisusing") {
		t.Errorf("missing help output: %q", stdout.String())
	}

	if stderr.Len() != 0 {
		t.Errorf("unexpected stderr: %q", stderr.String())
	}
}

func TestRunInvalidArguments(t *testing.T) {
	runner := New(
		&fakeFinder{},
		&fakeTerminator{},
		nil,
	)

	var stdout, stderr strings.Builder

	code := runner.Run(
		[]string{"-p", "99999"},
		&stdout,
		&stderr,
	)

	if code != 2 {
		t.Fatalf("Run() = %d, want 2", code)
	}

	if !strings.Contains(stderr.String(), "whoisusing:") {
		t.Errorf("unexpected stderr: %q", stderr.String())
	}

	if !strings.Contains(stderr.String(), "See 'whoisusing --help'") {
		t.Errorf("missing help hint: %q", stderr.String())
	}
}

func TestRunPortFree(t *testing.T) {
	finder := &fakeFinder{
		results: [][]process.Info{
			nil,
		},
	}

	runner := New(
		finder,
		&fakeTerminator{},
		func(time.Duration) {},
	)

	var stdout, stderr strings.Builder

	code := runner.Run(
		[]string{"-p", "8080"},
		&stdout,
		&stderr,
	)

	if code != 0 {
		t.Fatalf("Run() = %d, want 0", code)
	}

	want := "Port 8080 is currently free.\n"
	if stdout.String() != want {
		t.Errorf("stdout = %q, want %q", stdout.String(), want)
	}
}

func TestRunQuietPortFree(t *testing.T) {
	finder := &fakeFinder{
		results: [][]process.Info{
			nil,
		},
	}

	runner := New(
		finder,
		&fakeTerminator{},
		func(time.Duration) {},
	)

	var stdout, stderr strings.Builder

	code := runner.Run(
		[]string{"-q", "-p", "8080"},
		&stdout,
		&stderr,
	)

	if code != 1 {
		t.Fatalf("Run() = %d, want 1", code)
	}

	if stdout.Len() != 0 {
		t.Errorf("stdout = %q, want empty", stdout.String())
	}
}

func TestRunFoundProcesses(t *testing.T) {
	finder := &fakeFinder{
		results: [][]process.Info{
			{
				{
					Command:    "node",
					PID:        123,
					User:       "alice",
					Connection: "TCP *:8080 (LISTEN)",
				},
			},
		},
	}

	runner := New(
		finder,
		&fakeTerminator{},
		func(time.Duration) {},
	)

	var stdout, stderr strings.Builder

	code := runner.Run(
		[]string{"-p", "8080"},
		&stdout,
		&stderr,
	)

	if code != 0 {
		t.Fatalf("Run() = %d, want 0", code)
	}

	for _, want := range []string{
		"Processes using port 8080:",
		"node",
		"123",
		"alice",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("stdout does not contain %q:\n%s", want, stdout.String())
		}
	}
}

func TestRunQuietFoundProcesses(t *testing.T) {
	finder := &fakeFinder{
		results: [][]process.Info{
			{
				{PID: 123},
				{PID: 456},
			},
		},
	}

	runner := New(
		finder,
		&fakeTerminator{},
		func(time.Duration) {},
	)

	var stdout, stderr strings.Builder

	code := runner.Run(
		[]string{"-q", "-p", "8080"},
		&stdout,
		&stderr,
	)

	if code != 0 {
		t.Fatalf("Run() = %d, want 0", code)
	}

	want := "123\n456\n"
	if stdout.String() != want {
		t.Errorf("stdout = %q, want %q", stdout.String(), want)
	}
}

func TestRunFinderError(t *testing.T) {
	finder := &fakeFinder{
		results: [][]process.Info{
			nil,
		},
		errs: []error{
			errors.New("lsof unavailable"),
		},
	}

	runner := New(
		finder,
		&fakeTerminator{},
		nil,
	)

	var stdout, stderr strings.Builder

	code := runner.Run(
		[]string{"-p", "8080"},
		&stdout,
		&stderr,
	)

	if code != 1 {
		t.Fatalf("Run() = %d, want 1", code)
	}

	if !strings.Contains(stderr.String(), "lsof unavailable") {
		t.Errorf("stderr = %q", stderr.String())
	}
}

func TestRunKillSuccess(t *testing.T) {
	finder := &fakeFinder{
		results: [][]process.Info{
			{
				{PID: 123},
				{PID: 456},
			},
			nil,
		},
	}

	terminator := &fakeTerminator{}

	var sleepCalls int

	runner := New(
		finder,
		terminator,
		func(time.Duration) {
			sleepCalls++
		},
	)

	var stdout, stderr strings.Builder

	code := runner.Run(
		[]string{"-k", "-p", "8080"},
		&stdout,
		&stderr,
	)

	if code != 0 {
		t.Fatalf("Run() = %d, want 0", code)
	}

	wantPIDs := []int{123, 456}

	if len(terminator.pids) != len(wantPIDs) {
		t.Fatalf(
			"terminated PIDs = %v, want %v",
			terminator.pids,
			wantPIDs,
		)
	}

	for i := range wantPIDs {
		if terminator.pids[i] != wantPIDs[i] {
			t.Errorf(
				"terminated PIDs[%d] = %d, want %d",
				i,
				terminator.pids[i],
				wantPIDs[i],
			)
		}
	}

	if sleepCalls != 1 {
		t.Errorf("sleep called %d times, want 1", sleepCalls)
	}

	for _, want := range []string{
		"Attempting to gracefully terminate processes...",
		"Sent SIGTERM to PID 123",
		"Sent SIGTERM to PID 456",
		"Verifying port 8080 status...",
		"Success: All processes were terminated.",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("stdout does not contain %q:\n%s", want, stdout.String())
		}
	}
}

func TestRunKillFailure(t *testing.T) {
	finder := &fakeFinder{
		results: [][]process.Info{
			{
				{PID: 123},
			},
			{
				{PID: 123},
			},
		},
	}

	terminator := &fakeTerminator{
		err: errors.New("permission denied"),
	}

	runner := New(
		finder,
		terminator,
		func(time.Duration) {},
	)

	var stdout, stderr strings.Builder

	code := runner.Run(
		[]string{"-k", "-p", "8080"},
		&stdout,
		&stderr,
	)

	if code != 1 {
		t.Fatalf("Run() = %d, want 1", code)
	}

	if !strings.Contains(stdout.String(), "Failed to signal PID 123") {
		t.Errorf("stdout = %q", stdout.String())
	}

	if !strings.Contains(stdout.String(), "Unable to kill all processes") {
		t.Errorf("stdout = %q", stdout.String())
	}
}

func TestRunVerificationError(t *testing.T) {
	finder := &fakeFinder{
		results: [][]process.Info{
			{
				{PID: 123},
			},
			nil,
		},
		errs: []error{
			nil,
			errors.New("verification failed"),
		},
	}

	runner := New(
		finder,
		&fakeTerminator{},
		func(time.Duration) {},
	)

	var stdout, stderr strings.Builder

	code := runner.Run(
		[]string{"-k", "-p", "8080"},
		&stdout,
		&stderr,
	)

	if code != 1 {
		t.Fatalf("Run() = %d, want 1", code)
	}

	if !strings.Contains(stdout.String(), "Unable to verify port 8080") {
		t.Errorf("stdout = %q", stdout.String())
	}
}

func TestRunDoesNotKillWhenPortIsFree(t *testing.T) {
	finder := &fakeFinder{
		results: [][]process.Info{
			nil,
		},
	}

	terminator := &fakeTerminator{}

	runner := New(
		finder,
		terminator,
		func(time.Duration) {},
	)

	var stdout, stderr strings.Builder

	code := runner.Run(
		[]string{"-k", "-p", "8080"},
		&stdout,
		&stderr,
	)

	if code != 0 {
		t.Fatalf("Run() = %d, want 0", code)
	}

	if len(terminator.pids) != 0 {
		t.Errorf("terminated PIDs = %v, want none", terminator.pids)
	}
}

func TestRunQuietKillSuccess(t *testing.T) {
	finder := &fakeFinder{
		results: [][]process.Info{
			{
				{PID: 123},
			},
			nil,
		},
	}

	terminator := &fakeTerminator{}

	runner := New(
		finder,
		terminator,
		func(time.Duration) {},
	)

	var stdout, stderr strings.Builder

	code := runner.Run(
		[]string{"-kq", "-p", "8080"},
		&stdout,
		&stderr,
	)

	if code != 0 {
		t.Fatalf("Run() = %d, want 0", code)
	}

	if stdout.String() != "123\n" {
		t.Errorf("stdout = %q, want %q", stdout.String(), "123\n")
	}

	if len(terminator.pids) != 1 || terminator.pids[0] != 123 {
		t.Errorf("terminated PIDs = %v, want [123]", terminator.pids)
	}
}

func TestRunOutputError(t *testing.T) {
	finder := &fakeFinder{
		results: [][]process.Info{
			{
				{PID: 123},
			},
		},
	}

	runner := New(
		finder,
		&fakeTerminator{},
		nil,
	)

	stdout := failingWriter{
		err: errors.New("write failed"),
	}

	var stderr strings.Builder

	code := runner.Run(
		[]string{"-p", "8080"},
		stdout,
		&stderr,
	)

	if code != 1 {
		t.Fatalf("Run() = %d, want 1", code)
	}
}
