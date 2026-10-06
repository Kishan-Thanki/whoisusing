package output

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/kishan-thanki/whoisusing/internal/process"
)

func Table(w io.Writer, processes []process.Info) error {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)

	if _, err := fmt.Fprintln(tw, "COMMAND\tPID\tUSER\tCONNECTION STATE"); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(tw, "-------\t---\t----\t----------------"); err != nil {
		return err
	}

	for _, p := range processes {
		if _, err := fmt.Fprintf(
			tw,
			"%s\t%d\t%s\t%s\n",
			p.Command,
			p.PID,
			p.User,
			p.Connection,
		); err != nil {
			return err
		}
	}

	return tw.Flush()
}

func PIDs(w io.Writer, pids []int) error {
	for _, pid := range pids {
		if _, err := fmt.Fprintln(w, pid); err != nil {
			return err
		}
	}

	return nil
}

func Free(w io.Writer, port int) error {
	_, err := fmt.Fprintf(w, "Port %d is currently free.\n", port)
	return err
}

func Found(w io.Writer, port int) error {
	_, err := fmt.Fprintf(w, "Processes using port %d:\n\n", port)
	return err
}
