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

		switch arg {
		case "-h", "--help":
			return opts, ErrHelp

		case "-k":
			opts.Kill = true

		case "-q":
			opts.Quiet = true

		case "-p":
			if portSet {
				return opts, fmt.Errorf("port specified more than once")
			}

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

		default:
			if strings.HasPrefix(arg, "-") {
				return opts, fmt.Errorf("invalid flag: %s", arg)
			}

			return opts, fmt.Errorf("invalid argument: %s", arg)
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
  -p <port>    TCP port to inspect (required, range: %d-%d)
  -k           Gracefully terminate processes listening on the specified TCP port
  -q           Quiet mode: print only PIDs (one per line), no tables or messages
  -h, --help   Display this help message

Use each short flag separately; -p takes the next argument as its value.

Exit codes:
  0  success (normal mode: also when the port is free)
  1  quiet mode: no TCP listener found; or kill failed
  2  invalid arguments

Examples:
  whoisusing -p 8080
  whoisusing -k -p 8080
  whoisusing -q -p 8080
  PID=$(whoisusing -q -p 8080)
  whoisusing -k -q -p 8080
`, minPort, maxPort)
}
