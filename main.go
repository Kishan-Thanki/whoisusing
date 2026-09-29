package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"text/tabwriter"
	"time"
)

const (
	MinPort = 1
	MaxPort = 65535
)

var errHelp = errors.New("help requested")

type config struct {
	port  int
	kill  bool
	quiet bool
}

func help() {
	fmt.Println("Usage: whoisusing -p <port> [options]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Printf("  -p <port>    Port number to inspect (required, range: %d-%d)\n", MinPort, MaxPort)
	fmt.Println("  -k           Gracefully kill the process using the specified port")
	fmt.Println("  -q           Quiet mode: print only PIDs (one per line), no tables or messages")
	fmt.Println("  -h, --help   Display this help message")
	fmt.Println("")
	fmt.Println("Short flags can be combined; -p takes the next argument as its value.")
	fmt.Println("")
	fmt.Println("Exit codes:")
	fmt.Println("  0  success (normal mode: also when the port is free)")
	fmt.Println("  1  quiet mode: port is free; or kill failed")
	fmt.Println("  2  invalid arguments")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  whoisusing -p 8080")
	fmt.Println("  whoisusing -k -p 8080")
	fmt.Println("  whoisusing -kp 8080")
	fmt.Println("  whoisusing -qp 8080")
	fmt.Println("  PID=$(whoisusing -q -p 8080)")
	fmt.Println("  whoisusing -kqp 8080")
}

func parseArgs(args []string) (config, error) {
	var cfg config

	if len(args) <= 1 {
		return cfg, errHelp
	}

	portSet := false

	for i := 1; i < len(args); i++ {
		arg := args[i]

		if arg == "--help" {
			return cfg, errHelp
		}
		if len(arg) < 2 || arg[0] != '-' || strings.HasPrefix(arg, "--") {
			return cfg, fmt.Errorf("invalid argument: %s", arg)
		}

		needsValue := false
		for _, c := range arg[1:] {
			switch c {
			case 'h':
				return cfg, errHelp
			case 'k':
				cfg.kill = true
			case 'q':
				cfg.quiet = true
			case 'p':
				if needsValue || portSet {
					return cfg, fmt.Errorf("port specified more than once")
				}
				needsValue = true
			default:
				return cfg, fmt.Errorf("invalid flag: -%c", c)
			}
		}

		if needsValue {
			if i+1 >= len(args) {
				return cfg, fmt.Errorf("missing port")
			}
			i++
			p, err := strconv.Atoi(args[i])
			if err != nil {
				return cfg, fmt.Errorf("invalid port: %s", args[i])
			}
			cfg.port = p
			portSet = true
		}
	}

	if !portSet {
		return cfg, fmt.Errorf("missing -p <port>")
	}
	if cfg.port < MinPort || cfg.port > MaxPort {
		return cfg, fmt.Errorf("port out of range (%d-%d)", MinPort, MaxPort)
	}

	return cfg, nil
}

func formatLsofOutput(output []byte) string {
	var buf strings.Builder
	w := tabwriter.NewWriter(&buf, 0, 0, 3, ' ', 0)

	fmt.Fprintf(w, "COMMAND\tPID\tUSER\tCONNECTION STATE\n")
	fmt.Fprintf(w, "-------\t---\t----\t----------------\n")

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	for i := 1; i < len(lines); i++ {
		fields := strings.Fields(lines[i])
		if len(fields) >= 9 {
			details := strings.Join(fields[7:], " ")
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", fields[0], fields[1], fields[2], details)
		}
	}

	w.Flush()
	return buf.String()
}

func lookupPort(port int) ([]byte, []int) {
	cmd := exec.Command("lsof", "-i", fmt.Sprintf(":%d", port))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, nil
	}

	var pids []int
	seen := make(map[int]bool)
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for i := 1; i < len(lines); i++ {
		fields := strings.Fields(lines[i])
		if len(fields) >= 9 {
			if pid, err := strconv.Atoi(fields[1]); err == nil && !seen[pid] {
				seen[pid] = true
				pids = append(pids, pid)
			}
		}
	}
	return output, pids
}

func inspectPort(port int, quiet bool) []int {
	output, pids := lookupPort(port)

	if quiet {
		return pids
	}

	if len(pids) == 0 {
		fmt.Printf("Port %d is currently free.\n", port)
		return nil
	}

	fmt.Printf("Processes using port %d:\n\n", port)
	fmt.Print(formatLsofOutput(output))
	return pids
}

func main() {
	cfg, err := parseArgs(os.Args)
	if err != nil {
		if errors.Is(err, errHelp) {
			help()
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "whoisusing: %v\n", err)
		fmt.Fprintln(os.Stderr, "See 'whoisusing --help' for more information.")
		os.Exit(2)
	}

	pids := inspectPort(cfg.port, cfg.quiet)

	if cfg.quiet {
		for _, pid := range pids {
			fmt.Println(pid)
		}
		if len(pids) == 0 {
			os.Exit(1)
		}
	}

	if !cfg.kill || len(pids) == 0 {
		os.Exit(0)
	}

	if !cfg.quiet {
		fmt.Println("\nAttempting to gracefully terminate processes...")
	}

	for _, pid := range pids {
		proc, err := os.FindProcess(pid)
		if err != nil {
			continue
		}
		if err := proc.Signal(syscall.SIGTERM); err != nil {
			if !cfg.quiet {
				fmt.Printf("-> Failed to signal PID %d: %v\n", pid, err)
			}
			continue
		}
		if !cfg.quiet {
			fmt.Printf("-> Sent SIGTERM to PID %d\n", pid)
		}
	}

	time.Sleep(1 * time.Second)

	if !cfg.quiet {
		fmt.Printf("\nVerifying port %d status...\n", cfg.port)
	}

	remaining := inspectPort(cfg.port, cfg.quiet)

	if len(remaining) == 0 {
		if !cfg.quiet {
			fmt.Println("\nSuccess: All processes were terminated.")
		}
		os.Exit(0)
	}

	if !cfg.quiet {
		fmt.Println("\nError: Unable to kill all processes. They may require root/sudo privileges.")
	}
	os.Exit(1)
}
