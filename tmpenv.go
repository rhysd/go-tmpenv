package tmpenv

import (
	"os"
	"strings"
)

// Envguard is a state of environment variables which will be restored at Restore() method call.
type Envguard struct {
	maybeMod map[string]string
	maybeAdd map[string]struct{}
}

// Setenv sets the environment variable to the given value. The set variable will be restored when
// calling Restore().
func (guard *Envguard) Setenv(key, val string) error {
	if v, ok := os.LookupEnv(key); ok {
		guard.maybeMod[key] = v
	} else {
		guard.maybeAdd[key] = struct{}{}
	}
	return os.Setenv(key, val)
}

// Add adds an environment variable by key.
func (guard *Envguard) Add(keys ...string) {
	for _, k := range keys {
		if k == "" {
			continue
		}
		if v, ok := os.LookupEnv(k); ok {
			guard.maybeMod[k] = v
		} else {
			guard.maybeAdd[k] = struct{}{}
		}
	}
}

// Remove removes given keys from stored environment variables state.
func (guard *Envguard) Remove(keys ...string) (deleted bool) {
	for _, k := range keys {
		if _, ok := guard.maybeAdd[k]; ok {
			delete(guard.maybeAdd, k)
			deleted = true
		}
		if _, ok := guard.maybeMod[k]; ok {
			delete(guard.maybeMod, k)
			deleted = true
		}
	}
	return
}

// Unsetenv unsets the environment variable. When Restore() method is called, the unset variable will
// be restored. When underlying os.Unsetenv() call failed, it returns an error
func (guard *Envguard) Unsetenv(key string) error {
	if key == "" {
		return nil
	}
	guard.Add(key)
	return os.Unsetenv(key)
}

// Restore restores stored environment variable values. This method is usually called with 'defer' to
// ensure the state to be restored. This function returns an error when underlying Setenv() calls
// returned an error except for on Windows. On Windows this function returns nil always since some
// environment variables are not set on Windows and it is intentional.
func (guard *Envguard) Restore() error {
	for k, v := range guard.maybeMod {
		if k == "" {
			continue
		}
		if err := os.Setenv(k, v); err != nil {
			return err
		}
	}
	for k := range guard.maybeAdd {
		if k == "" {
			continue
		}
		if err := os.Unsetenv(k); err != nil {
			return err
		}
	}
	return nil
}

// New creates a new Envguard instance. If one ore more keys are given, corresponding environment variables
// will be restored at Restore() method call.
func New(keys ...string) *Envguard {
	g := &Envguard{map[string]string{}, map[string]struct{}{}}
	if len(keys) == 0 {
		return g
	}
	g.Add(keys...)
	return g
}

// Setenvs sets environment variables with given key-values. And creates a new Envguard instance to restore
// them when calling Restore().
func Setenvs(m map[string]string) (*Envguard, error) {
	g := New()
	for k, v := range m {
		if err := g.Setenv(k, v); err != nil {
			return nil, err
		}
	}
	return g, nil
}

// All creates a new Envguard instance with all existing environment variables. All environment variables
// will be restored by calling Restore().
func All() *Envguard {
	kv := os.Environ()
	m := make(map[string]string, len(kv))
	for _, s := range kv {
		if idx := strings.IndexRune(s, '='); idx > 0 {
			if k := s[:idx]; k != "" {
				m[k] = s[idx+1:]
			}
		}
	}
	return &Envguard{m, map[string]struct{}{}}
}

// Unset creates a new Envguard instance with unsetting all given anvironment variables. The unset variables
// will be restored by calling Restore() method. When underlying os.Unsetenv() failed, it returns an error.
func Unset(keys ...string) (*Envguard, error) {
	g := New()
	for _, k := range keys {
		if err := g.Unsetenv(k); err != nil {
			return nil, err
		}
	}
	return g, nil
}

// UnsetAll creates a new Envguard instance with clearing all environment variables.
func UnsetAll() (*Envguard, error) {
	g := All()
	for k := range g.maybeMod {
		if err := os.Unsetenv(k); err != nil {
			return nil, err
		}
	}
	return g, nil
}
