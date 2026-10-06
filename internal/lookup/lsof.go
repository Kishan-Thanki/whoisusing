package lookup

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/kishan-thanki/whoisusing/internal/process"
)

const (
	lsofFieldFormat = "pcLPnT0"
	lsofPortFormat  = ":%d"
)

var (
	errNoMatch     = errors.New("no match")
	errUnsupported = errors.New("unsupported platform")
)

type Lsof struct {
	command string
}

func NewLsof() Lsof {
	return Lsof{
		command: "lsof",
	}
}

func (l Lsof) Find(port int) ([]process.Info, error) {
	output, stderr, err := execute(l.command, port)
	if err != nil {
		if errors.Is(err, errNoMatch) {
			return nil, nil
		}

		if errors.Is(err, errUnsupported) {
			return nil, err
		}

		return nil, fmt.Errorf(
			"lsof failed: %w: %s",
			err,
			strings.TrimSpace(stderr),
		)
	}

	return parse(output), nil
}

func parse(output []byte) []process.Info {
	var (
		processes []process.Info
		current   process.Info
		havePID   bool

		protocol string
		address  string
		state    string

		seen = make(map[int]struct{})
	)

	flush := func() {
		if !havePID {
			return
		}

		if _, exists := seen[current.PID]; exists {
			return
		}

		current.Connection = formatConnection(
			protocol,
			address,
			state,
		)

		seen[current.PID] = struct{}{}
		processes = append(processes, current)
	}

	for _, field := range bytes.Split(output, []byte{0}) {
		if len(field) < 2 {
			continue
		}

		value := string(field[1:])

		switch field[0] {
		case 'p':
			flush()

			pid, err := strconv.Atoi(value)
			if err != nil || pid <= 0 {
				havePID = false
				continue
			}

			current = process.Info{
				PID: pid,
			}

			protocol = ""
			address = ""
			state = ""
			havePID = true

		case 'c':
			if havePID {
				current.Command = value
			}

		case 'L':
			if havePID {
				current.User = value
			}

		case 'P':
			if havePID && protocol == "" {
				protocol = value
			}

		case 'n':
			if havePID && address == "" {
				address = value
			}

		case 'T':
			if havePID && state == "" && strings.HasPrefix(value, "ST=") {
				state = strings.TrimPrefix(value, "ST=")
			}
		}
	}

	flush()

	return processes
}

func formatConnection(protocol, address, state string) string {
	parts := make([]string, 0, 3)

	if protocol != "" {
		parts = append(parts, protocol)
	}

	if address != "" {
		parts = append(parts, address)
	}

	if state != "" {
		parts = append(parts, "("+state+")")
	}

	return strings.Join(parts, " ")
}
