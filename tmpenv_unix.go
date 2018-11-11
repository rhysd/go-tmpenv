// +build !windows

package tmpenv

import (
	"os"
)

func (guard *Envguard) restore() error {
	for k, v := range guard.maybeMod {
		if err := os.Setenv(k, v); err != nil {
			return err
		}
	}
	for k := range guard.maybeAdd {
		if err := os.Unsetenv(k); err != nil {
			return err
		}
	}
	return nil
}
