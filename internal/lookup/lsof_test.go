package lookup

import (
	"errors"
	"strings"
	"testing"

	"github.com/kishan-thanki/whoisusing/internal/process"
)

func TestParse(t *testing.T) {
	input := []byte(
		"p123\x00" +
			"cnode\x00" +
			"Lalice\x00" +
			"PTCP\x00" +
			"n*:8080\x00" +
			"TST=LISTEN\x00" +
			"p456\x00" +
			"cpython3\x00" +
			"Lbob\x00" +
			"PTCP\x00" +
			"n*:8080\x00" +
			"TST=LISTEN\x00",
	)

	got := parse(input)

	want := []process.Info{
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

	if len(got) != len(want) {
		t.Fatalf("parse() returned %d processes, want %d", len(got), len(want))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Errorf("parse()[%d] = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestParseEmpty(t *testing.T) {
	for _, input := range [][]byte{
		nil,
		{},
		[]byte(""),
		[]byte("\x00"),
	} {
		got := parse(input)

		if len(got) != 0 {
			t.Errorf("parse(%q) = %+v, want empty", input, got)
		}
	}
}

func TestParseMalformedPID(t *testing.T) {
	input := []byte(
		"pabc\x00" +
			"cnode\x00" +
			"Lalice\x00" +
			"PTCP\x00" +
			"n*:8080\x00" +
			"TST=LISTEN\x00",
	)

	got := parse(input)

	if len(got) != 0 {
		t.Errorf("parse() = %+v, want empty", got)
	}
}

func TestParseDeduplicatesPIDs(t *testing.T) {
	input := []byte(
		"p123\x00" +
			"cnode\x00" +
			"Lalice\x00" +
			"PTCP\x00" +
			"n*:8080\x00" +
			"TST=LISTEN\x00" +
			"p123\x00" +
			"cnode\x00" +
			"Lalice\x00" +
			"PTCP\x00" +
			"n127.0.0.1:8080->127.0.0.1:5000\x00" +
			"TST=ESTABLISHED\x00",
	)

	got := parse(input)

	if len(got) != 1 {
		t.Fatalf("parse() returned %d processes, want 1", len(got))
	}

	if got[0].PID != 123 {
		t.Fatalf("PID = %d, want 123", got[0].PID)
	}
}

func TestParseMissingOptionalFields(t *testing.T) {
	input := []byte(
		"p123\x00" +
			"cnode\x00" +
			"Lalice\x00",
	)

	got := parse(input)

	want := []process.Info{
		{
			Command: "node",
			PID:     123,
			User:    "alice",
		},
	}

	if len(got) != 1 {
		t.Fatalf("parse() returned %d processes, want 1", len(got))
	}

	if got[0] != want[0] {
		t.Errorf("parse() = %+v, want %+v", got[0], want[0])
	}
}

func TestFormatConnection(t *testing.T) {
	tests := []struct {
		name     string
		protocol string
		address  string
		state    string
		want     string
	}{
		{
			name:     "complete",
			protocol: "TCP",
			address:  "*:8080",
			state:    "LISTEN",
			want:     "TCP *:8080 (LISTEN)",
		},
		{
			name:     "without state",
			protocol: "TCP",
			address:  "*:8080",
			want:     "TCP *:8080",
		},
		{
			name:    "address only",
			address: "*:8080",
			want:    "*:8080",
		},
		{
			name: "empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatConnection(tt.protocol, tt.address, tt.state)

			if got != tt.want {
				t.Errorf(
					"formatConnection() = %q, want %q",
					got,
					tt.want,
				)
			}
		})
	}
}
func TestLsofFindCommandError(t *testing.T) {
	lsof := Lsof{
		command: "whoisusing-command-that-does-not-exist",
	}

	_, err := lsof.Find(8080)
	if err == nil {
		t.Fatal("Find() error = nil, want error")
	}

	if errors.Is(err, errUnsupported) {
		t.Skip("lsof is unsupported on this platform")
	}

	if !strings.Contains(err.Error(), "lsof failed") {
		t.Fatalf("Find() error = %q, want lsof failed", err)
	}
}

func TestLsofFindEmptyResult(t *testing.T) {
	lsof := Lsof{
		command: "sh",
	}

	_, err := lsof.Find(8080)
	if err == nil {
		t.Fatal("Find() error = nil, want error")
	}
}
