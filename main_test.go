package main

import (
	"errors"
	"strings"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		want        config
		expectError bool
	}{
		{"Valid Port", []string{"whoisusing", "-p", "8080"}, config{port: 8080}, false},
		{"Min Port", []string{"whoisusing", "-p", "1"}, config{port: 1}, false},
		{"Max Port", []string{"whoisusing", "-p", "65535"}, config{port: 65535}, false},

		{"Separate (-k -p)", []string{"whoisusing", "-k", "-p", "8080"}, config{port: 8080, kill: true}, false},
		{"Reverse (-p 8080 -k)", []string{"whoisusing", "-p", "8080", "-k"}, config{port: 8080, kill: true}, false},
		{"Cluster (-pk)", []string{"whoisusing", "-pk", "8080"}, config{port: 8080, kill: true}, false},
		{"Cluster (-kp)", []string{"whoisusing", "-kp", "8080"}, config{port: 8080, kill: true}, false},

		{"Quiet (-q -p)", []string{"whoisusing", "-q", "-p", "8080"}, config{port: 8080, quiet: true}, false},
		{"Quiet cluster (-qp)", []string{"whoisusing", "-qp", "8080"}, config{port: 8080, quiet: true}, false},
		{"Quiet cluster (-pq)", []string{"whoisusing", "-pq", "8080"}, config{port: 8080, quiet: true}, false},
		{"Triple (-kqp)", []string{"whoisusing", "-kqp", "8080"}, config{port: 8080, kill: true, quiet: true}, false},
		{"Triple (-pkq)", []string{"whoisusing", "-pkq", "8080"}, config{port: 8080, kill: true, quiet: true}, false},
		{"Triple (-qkp)", []string{"whoisusing", "-qkp", "8080"}, config{port: 8080, kill: true, quiet: true}, false},
		{"Triple (-kpq)", []string{"whoisusing", "-kpq", "8080"}, config{port: 8080, kill: true, quiet: true}, false},
		{"Mixed (-qk -p)", []string{"whoisusing", "-qk", "-p", "8080"}, config{port: 8080, kill: true, quiet: true}, false},
		{"Mixed (-q -k -p)", []string{"whoisusing", "-q", "-k", "-p", "8080"}, config{port: 8080, kill: true, quiet: true}, false},
		{"Mixed (-p 8080 -q -k)", []string{"whoisusing", "-p", "8080", "-q", "-k"}, config{port: 8080, kill: true, quiet: true}, false},
		{"Mixed (-kq -p)", []string{"whoisusing", "-kq", "-p", "8080"}, config{port: 8080, kill: true, quiet: true}, false},
		{"Repeated flag (-q -q -p)", []string{"whoisusing", "-q", "-q", "-p", "8080"}, config{port: 8080, quiet: true}, false},

		{"No Args", []string{"whoisusing"}, config{}, true},
		{"Help Flag", []string{"whoisusing", "--help"}, config{}, true},
		{"Short Help", []string{"whoisusing", "-h"}, config{}, true},
		{"Help in Cluster", []string{"whoisusing", "-kh", "-p", "8080"}, config{}, true},
		{"Wrong Flag", []string{"whoisusing", "-x", "8080"}, config{}, true},
		{"Typo in Cluster (-kap)", []string{"whoisusing", "-kap", "8080"}, config{}, true},
		{"Unknown Long Flag", []string{"whoisusing", "--foo"}, config{}, true},
		{"Bare Dash", []string{"whoisusing", "-"}, config{}, true},
		{"Bare Number", []string{"whoisusing", "8080"}, config{}, true},
		{"Letters Instead of Port", []string{"whoisusing", "-p", "abc"}, config{}, true},
		{"Flag as Port Value", []string{"whoisusing", "-p", "-k"}, config{}, true},
		{"Port Out of Bounds High", []string{"whoisusing", "-p", "99999"}, config{}, true},
		{"Port Out of Bounds Low", []string{"whoisusing", "-p", "0"}, config{}, true},
		{"Negative Port", []string{"whoisusing", "-p", "-5"}, config{}, true},
		{"Kill Flag but No Port", []string{"whoisusing", "-k"}, config{}, true},
		{"Quiet Flag but No Port", []string{"whoisusing", "-q"}, config{}, true},
		{"Cluster but No Port (-pk)", []string{"whoisusing", "-pk"}, config{}, true},
		{"Cluster but No Port (-kqp)", []string{"whoisusing", "-kqp"}, config{}, true},
		{"Bare Port (-k 8080)", []string{"whoisusing", "-k", "8080"}, config{}, true},
		{"Bare Port (-q 8080)", []string{"whoisusing", "-q", "8080"}, config{}, true},
		{"Bare Port (-qk 8080)", []string{"whoisusing", "-qk", "8080"}, config{}, true},
		{"Bare Port (-kq 8080)", []string{"whoisusing", "-kq", "8080"}, config{}, true},
		{"Duplicate -p in Cluster", []string{"whoisusing", "-pp", "8080"}, config{}, true},
		{"Duplicate -p Separate", []string{"whoisusing", "-p", "8080", "-p", "9090"}, config{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseArgs(tt.args)

			if (err != nil) != tt.expectError {
				t.Fatalf("parseArgs() error = %v, expectError %v", err, tt.expectError)
			}
			if tt.expectError {
				return
			}
			if got != tt.want {
				t.Errorf("parseArgs() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func permutations(items []rune) [][]rune {
	if len(items) <= 1 {
		return [][]rune{append([]rune{}, items...)}
	}
	var out [][]rune
	for i := range items {
		rest := make([]rune, 0, len(items)-1)
		rest = append(rest, items[:i]...)
		rest = append(rest, items[i+1:]...)
		for _, p := range permutations(rest) {
			out = append(out, append([]rune{items[i]}, p...))
		}
	}
	return out
}

func clusterings(seq []rune) [][]string {
	n := len(seq)
	var out [][]string
	for mask := 0; mask < 1<<(n-1); mask++ {
		var groups []string
		cur := "-" + string(seq[0])
		for i := 1; i < n; i++ {
			if mask&(1<<(i-1)) != 0 {
				groups = append(groups, cur)
				cur = "-"
			}
			cur += string(seq[i])
		}
		groups = append(groups, cur)
		out = append(out, groups)
	}
	return out
}

func TestParseArgsExhaustiveValid(t *testing.T) {
	subsets := [][]rune{{}, {'k'}, {'q'}, {'k', 'q'}}

	for _, sub := range subsets {
		want := config{port: 8080}
		for _, f := range sub {
			if f == 'k' {
				want.kill = true
			}
			if f == 'q' {
				want.quiet = true
			}
		}

		flags := append(append([]rune{}, sub...), 'p')
		for _, perm := range permutations(flags) {
			for _, groups := range clusterings(perm) {
				args := []string{"whoisusing"}
				for _, g := range groups {
					args = append(args, g)
					if strings.ContainsRune(g, 'p') {
						args = append(args, "8080")
					}
				}

				got, err := parseArgs(args)
				if err != nil {
					t.Errorf("parseArgs(%v) unexpected error: %v", args, err)
					continue
				}
				if got != want {
					t.Errorf("parseArgs(%v) = %+v, want %+v", args, got, want)
				}
			}
		}
	}
}

func TestParseArgsExhaustiveNoPort(t *testing.T) {
	subsets := [][]rune{{'k'}, {'q'}, {'k', 'q'}}

	for _, sub := range subsets {
		for _, perm := range permutations(sub) {
			for _, groups := range clusterings(perm) {
				base := append([]string{"whoisusing"}, groups...)

				if _, err := parseArgs(base); err == nil {
					t.Errorf("parseArgs(%v) should fail: no -p", base)
				}

				withPort := append(append([]string{}, base...), "8080")
				if _, err := parseArgs(withPort); err == nil {
					t.Errorf("parseArgs(%v) should fail: bare port without -p", withPort)
				}
			}
		}
	}
}

func TestParseArgsHelpSentinel(t *testing.T) {
	helpCases := [][]string{
		{"whoisusing"},
		{"whoisusing", "--help"},
		{"whoisusing", "-h"},
		{"whoisusing", "-kh", "-p", "8080"},
	}

	for _, args := range helpCases {
		if _, err := parseArgs(args); !errors.Is(err, errHelp) {
			t.Errorf("parseArgs(%v) error = %v, want errHelp", args, err)
		}
	}

	errCases := [][]string{
		{"whoisusing", "-x", "8080"},
		{"whoisusing", "-p", "abc"},
		{"whoisusing", "-p", "0"},
		{"whoisusing", "-q", "8080"},
	}

	for _, args := range errCases {
		_, err := parseArgs(args)
		if err == nil || errors.Is(err, errHelp) {
			t.Errorf("parseArgs(%v) error = %v, want non-help error", args, err)
		}
	}
}

func TestFormatLsofOutput(t *testing.T) {
	rawLsof := `COMMAND   PID USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME
python3 12345 user_name   3u  IPv4 0xdeadbeef      0t0  TCP *:http-alt (LISTEN)`

	expectedOutput := `COMMAND   PID     USER        CONNECTION STATE
-------   ---     ----        ----------------
python3   12345   user_name   TCP *:http-alt (LISTEN)
`

	got := formatLsofOutput([]byte(rawLsof))

	if got != expectedOutput {
		t.Errorf("formatLsofOutput() mismatch.\nWant:\n%s\nGot:\n%s", expectedOutput, got)
	}
}
