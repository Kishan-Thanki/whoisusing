package cli

import (
	"errors"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want Options
	}{
		{
			name: "port only",
			args: []string{"-p", "8080"},
			want: Options{Port: 8080},
		},
		{
			name: "kill",
			args: []string{"-k", "-p", "8080"},
			want: Options{Port: 8080, Kill: true},
		},
		{
			name: "quiet",
			args: []string{"-q", "-p", "8080"},
			want: Options{Port: 8080, Quiet: true},
		},
		{
			name: "kill and quiet",
			args: []string{"-kq", "-p", "8080"},
			want: Options{
				Port:  8080,
				Kill:  true,
				Quiet: true,
			},
		},
		{
			name: "combined flags",
			args: []string{"-kqp", "8080"},
			want: Options{
				Port:  8080,
				Kill:  true,
				Quiet: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if got != tt.want {
				t.Errorf("Parse() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"missing arguments", nil},
		{"missing port", []string{"-k"}},
		{"missing port value", []string{"-p"}},
		{"invalid port", []string{"-p", "abc"}},
		{"port too small", []string{"-p", "0"}},
		{"port too large", []string{"-p", "65536"}},
		{"duplicate port", []string{"-p", "8080", "-p", "8081"}},
		{"invalid flag", []string{"-x", "-p", "8080"}},
		{"positional argument", []string{"8080"}},
		{"unknown long option", []string{"--port", "8080"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.args)

			if tt.name == "missing arguments" {
				if !errors.Is(err, ErrHelp) {
					t.Fatalf("Parse() error = %v, want ErrHelp", err)
				}
				return
			}

			if err == nil {
				t.Fatal("Parse() error = nil, want error")
			}
		})
	}
}

func TestParseHelp(t *testing.T) {
	for _, args := range [][]string{
		{"-h"},
		{"--help"},
	} {
		_, err := Parse(args)

		if !errors.Is(err, ErrHelp) {
			t.Fatalf("Parse(%v) error = %v, want ErrHelp", args, err)
		}
	}
}

func TestParsePortBoundaries(t *testing.T) {
	for _, port := range []string{"1", "65535"} {
		opts, err := Parse([]string{"-p", port})
		if err != nil {
			t.Fatalf("Parse(%q) error = %v", port, err)
		}

		if opts.Port == 0 {
			t.Fatalf("Parse(%q) returned zero port", port)
		}
	}
}

func TestHelp(t *testing.T) {
	got := Help()

	for _, want := range []string{
		"Usage: whoisusing -p <port> [options]",
		"-p <port>",
		"-k",
		"-q",
		"-h, --help",
		"Exit codes:",
		"Examples:",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("Help() does not contain %q", want)
		}
	}
}
