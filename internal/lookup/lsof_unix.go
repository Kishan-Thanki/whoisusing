//go:build darwin || freebsd || illumos || linux || netbsd || openbsd || solaris

package lookup

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)

func execute(command string, port int) ([]byte, string, error) {
	cmd := exec.Command(
		command,
		"-nP",
		"-F",
		lsofFieldFormat,
		"-Ts",
		"-i",
		fmt.Sprintf(lsofPortFormat, port),
	)

	var stderr bytes.Buffer

	cmd.Stderr = &stderr

	output, err := cmd.Output()
	if err == nil {
		return output, stderr.String(), nil
	}

	var exitErr *exec.ExitError

	if errors.As(err, &exitErr) &&
		exitErr.ExitCode() == 1 &&
		len(output) == 0 &&
		stderr.Len() == 0 {
		return nil, "", errNoMatch
	}

	return nil, stderr.String(), err
}
