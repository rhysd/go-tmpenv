package envguard

import (
	"os"
	"strings"
)

// Envguard is a state of restored environment variables.
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

// Restore restores stored environment variable values.
func (guard *Envguard) Restore() error {
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

// New creates a new Envguard instance.
func New() *Envguard {
	return &Envguard{map[string]string{}, map[string]struct{}{}}
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
		if idx := strings.IndexRune(s, '='); idx >= 0 {
			m[s[:idx]] = s[idx+1:]
		}
	}
	return &Envguard{m, map[string]struct{}{}}
}
