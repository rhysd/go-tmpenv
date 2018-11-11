// +build windows

package tmpenv

import (
	"os"
)

func (guard *Envguard) restore() error {
	// Returned errors should not be checked since on Windows some environment variables cannot be
	// modified.
	for k, v := range guard.maybeMod {
		os.Setenv(k, v)
	}
	for k := range guard.maybeAdd {
		os.Unsetenv(k)
	}
	return nil
}
