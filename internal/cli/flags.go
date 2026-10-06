package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	minPort = 1
	maxPort = 65535
)

type Options struct {
	Port  int
	Kill  bool
	Quiet bool
}

var ErrHelp = errors.New("help requested")

func Parse(args []string) (Options, error) {
	var opts Options

	if len(args) == 0 {
		return opts, ErrHelp
	}

	portSet := false

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "--help" {
			return opts, ErrHelp
		}

		if len(arg) < 2 || arg[0] != '-' || strings.HasPrefix(arg, "--") {
			return opts, fmt.Errorf("invalid argument: %s", arg)
		}

		needsValue := false

		for _, flag := range arg[1:] {
			switch flag {
			case 'h':
				return opts, ErrHelp

			case 'k':
				opts.Kill = true

			case 'q':
				opts.Quiet = true

			case 'p':
				if needsValue || portSet {
					return opts, fmt.Errorf("port specified more than once")
				}
				needsValue = true

			default:
				return opts, fmt.Errorf("invalid flag: -%c", flag)
			}
		}

		if needsValue {
			if i+1 >= len(args) {
				return opts, fmt.Errorf("missing port")
			}

			i++

			port, err := strconv.Atoi(args[i])
			if err != nil {
				return opts, fmt.Errorf("invalid port: %s", args[i])
			}

			opts.Port = port
			portSet = true
		}
	}

	if !portSet {
		return opts, fmt.Errorf("missing -p <port>")
	}

	if opts.Port < minPort || opts.Port > maxPort {
		return opts, fmt.Errorf("port out of range (%d-%d)", minPort, maxPort)
	}

	return opts, nil
}

func Help() string {
	return fmt.Sprintf(`Usage: whoisusing -p <port> [options]

Options:
  -p <port>    Port number to inspect (required, range: %d-%d)
  -k           Gracefully kill the process using the specified port
  -q           Quiet mode: print only PIDs (one per line), no tables or messages
  -h, --help   Display this help message

Short flags can be combined; -p takes the next argument as its value.

Exit codes:
  0  success (normal mode: also when the port is free)
  1  quiet mode: port is free; or kill failed
  2  invalid arguments

Examples:
  whoisusing -p 8080
  whoisusing -k -p 8080
  whoisusing -kp 8080
  whoisusing -qp 8080
  PID=$(whoisusing -q -p 8080)
  whoisusing -kqp 8080
`, minPort, maxPort)
}
