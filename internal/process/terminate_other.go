//go:build !unix

package process

import (
	"errors"
	"os"
)

var errUnsupported = errors.New("graceful process termination is not supported on this platform")

func terminate(*os.Process) error {
	return errUnsupported
}
