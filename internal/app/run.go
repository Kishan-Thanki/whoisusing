package app

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/kishan-thanki/whoisusing/internal/cli"
	"github.com/kishan-thanki/whoisusing/internal/output"
	"github.com/kishan-thanki/whoisusing/internal/process"
)

const terminateWait = time.Second

type finder interface {
	Find(port int) ([]process.Info, error)
}

type terminator interface {
	Terminate(pid int) error
}

type Runner struct {
	finder     finder
	terminator terminator
	sleep      func(time.Duration)
}

func New(
	finder finder,
	terminator terminator,
	sleep func(time.Duration),
) *Runner {
	if sleep == nil {
		sleep = time.Sleep
	}

	return &Runner{
		finder:     finder,
		terminator: terminator,
		sleep:      sleep,
	}
}

func (r *Runner) Run(args []string, stdout, stderr io.Writer) int {
	opts, err := cli.Parse(args)
	if err != nil {
		return handleParseError(err, stdout, stderr)
	}

	processes, err := r.finder.Find(opts.Port)
	if err != nil {
		fmt.Fprintf(
			stderr,
			"whoisusing: unable to inspect port %d: %v\n",
			opts.Port,
			err,
		)
		return 1
	}

	if err := display(stdout, opts, processes); err != nil {
		return 1
	}

	if len(processes) == 0 || !opts.Kill {
		if opts.Quiet && len(processes) == 0 {
			return 1
		}

		return 0
	}

	return r.kill(opts.Port, processes, opts.Quiet, stdout)
}

func handleParseError(err error, stdout, stderr io.Writer) int {
	if errors.Is(err, cli.ErrHelp) {
		if _, writeErr := io.WriteString(stdout, cli.Help()); writeErr != nil {
			return 1
		}
		return 0
	}

	fmt.Fprintf(stderr, "whoisusing: %v\n", err)
	fmt.Fprintln(stderr, "See 'whoisusing --help' for more information.")

	return 2
}

func display(
	w io.Writer,
	opts cli.Options,
	processes []process.Info,
) error {
	if opts.Quiet {
		return output.PIDs(w, process.PIDs(processes))
	}

	if len(processes) == 0 {
		return output.Free(w, opts.Port)
	}

	if err := output.Found(w, opts.Port); err != nil {
		return err
	}

	return output.Table(w, processes)
}

func (r *Runner) kill(
	port int,
	processes []process.Info,
	quiet bool,
	stdout io.Writer,
) int {
	if !quiet {
		fmt.Fprintln(stdout)
		fmt.Fprintln(stdout, "Attempting to gracefully terminate processes...")
	}

	for _, p := range processes {
		if err := r.terminator.Terminate(p.PID); err != nil {
			if !quiet {
				fmt.Fprintf(
					stdout,
					"-> Failed to signal PID %d: %v\n",
					p.PID,
					err,
				)
			}
			continue
		}

		if !quiet {
			fmt.Fprintf(stdout, "-> Sent SIGTERM to PID %d\n", p.PID)
		}
	}

	r.sleep(terminateWait)

	return r.verify(port, quiet, stdout)
}

func (r *Runner) verify(
	port int,
	quiet bool,
	stdout io.Writer,
) int {
	if !quiet {
		fmt.Fprintln(stdout)
		fmt.Fprintf(stdout, "Verifying port %d status...\n", port)
	}

	remaining, err := r.finder.Find(port)
	if err != nil {
		if !quiet {
			fmt.Fprintf(
				stdout,
				"\nError: Unable to verify port %d: %v\n",
				port,
				err,
			)
		}
		return 1
	}

	if len(remaining) == 0 {
		if !quiet {
			fmt.Fprintln(stdout)
			fmt.Fprintln(stdout, "Success: All processes were terminated.")
		}
		return 0
	}

	if !quiet {
		fmt.Fprintln(
			stdout,
			"\nError: Unable to kill all processes. They may require root/sudo privileges.",
		)
	}

	return 1
}
