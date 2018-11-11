package envguard

import (
	"os"
	"sort"
	"testing"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func fatalIfErr(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestAdd(t *testing.T) {
	envs := map[string]string{
		"ENVGUARD_TEST_ADD_FOO": "foo",
		"ENVGUARD_TEST_ADD_BAR": "",
	}
	defer os.Unsetenv("ENVGUARD_TEST_ADD_PIYO")

	for k, v := range envs {
		panicIfErr(os.Setenv(k, v))
	}

	g := New()
	g.Add("ENVGUARD_TEST_ADD_FOO", "ENVGUARD_TEST_ADD_BAR", "ENVGUARD_TEST_ADD_PIYO")

	if v, ok := os.LookupEnv("ENVGUARD_TEST_ADD_PIYO"); ok {
		t.Fatal("$ENVGUARD_TEST_ADD_PIYO should not exist:", v)
	}

	panicIfErr(os.Setenv("ENVGUARD_TEST_ADD_FOO", "aaa"))
	panicIfErr(os.Setenv("ENVGUARD_TEST_ADD_BAR", "bbb"))
	panicIfErr(os.Setenv("ENVGUARD_TEST_ADD_PIYO", "ccc"))

	if err := g.Restore(); err != nil {
		t.Fatal(err)
	}

	for k, v := range envs {
		e := os.Getenv(k)
		if e != v {
			t.Error("Environment variable was not restured:", k, v, e)
		}
	}

	if e, ok := os.LookupEnv("ENVGUARD_TEST_ADD_PIYO"); ok {
		t.Error("Environment variable was not erased: ", e)
	}
}

func TestSetenv(t *testing.T) {
	envs := map[string]string{
		"ENVGUARD_TEST_SETENV_FOO": "foo",
		"ENVGUARD_TEST_SETENV_BAR": "",
	}
	defer os.Unsetenv("ENVGUARD_TEST_SETENV_PIYO")

	for k, v := range envs {
		panicIfErr(os.Setenv(k, v))
	}

	g := New()

	added := map[string]string{
		"ENVGUARD_TEST_SETENV_FOO":  "aaa",
		"ENVGUARD_TEST_SETENV_BAR":  "bbb",
		"ENVGUARD_TEST_SETENV_PIYO": "ccc",
	}
	for k, v := range added {
		if err := g.Setenv(k, v); err != nil {
			t.Fatal(err, k, v)
		}
	}

	for k, v := range added {
		e, ok := os.LookupEnv(k)
		if !ok {
			t.Fatal("Env var was not set", k, v)
		}
		if e != v {
			t.Fatal("Env var is unexpected. wanted", v, "but have", e)
		}
	}

	if err := g.Restore(); err != nil {
		t.Fatal(err)
	}

	for k, v := range envs {
		e := os.Getenv(k)
		if e != v {
			t.Error("Environment variable was not restured:", k, v, e)
		}
	}

	if e, ok := os.LookupEnv("ENVGUARD_TEST_SETENV_PIYO"); ok {
		t.Error("Environment variable was not erased: ", e)
	}
}

func TestRemove(t *testing.T) {
	envs := map[string]string{
		"ENVGUARD_TEST_REMOVE_FOO": "foo",
		"ENVGUARD_TEST_REMOVE_BAR": "",
	}
	for k, v := range envs {
		panicIfErr(os.Setenv(k, v))
	}
	defer os.Unsetenv("ENVGUARD_TEST_REMOVE_FOO")
	defer os.Unsetenv("ENVGUARD_TEST_REMOVE_BAR")

	g := New()
	g.Add("ENVGUARD_TEST_REMOVE_FOO", "ENVGUARD_TEST_REMOVE_BAR", "ENVGUARD_TEST_REMOVE_PIYO")

	panicIfErr(os.Setenv("ENVGUARD_TEST_REMOVE_FOO", "aaa"))
	panicIfErr(os.Setenv("ENVGUARD_TEST_REMOVE_BAR", "bbb"))

	if !g.Remove("ENVGUARD_TEST_REMOVE_BAR", "ENVGUARD_TEST_REMOVE_PIYO", "ENVGUARD_TEST_REMOVE_HOGE") {
		t.Fatal("Nothing was removed")
	}
	if g.Remove("ENVGUARD_TEST_REMOVE_UNKNOWN") {
		t.Fatal("Unknown env var was removed")
	}

	if err := g.Restore(); err != nil {
		t.Fatal(err)
	}

	if v := os.Getenv("ENVGUARD_TEST_REMOVE_FOO"); v != "foo" {
		t.Error("Env var was not restored", v)
	}
	if v := os.Getenv("ENVGUARD_TEST_REMOVE_BAR"); v != "bbb" {
		t.Error("Env var was restored", v)
	}
}

func TestSetenvs(t *testing.T) {
	envs := map[string]string{
		"ENVGUARD_TEST_SETENVS_FOO": "foo",
		"ENVGUARD_TEST_SETENVS_BAR": "",
	}
	defer os.Unsetenv("ENVGUARD_TEST_SETENVS_PIYO")

	for k, v := range envs {
		panicIfErr(os.Setenv(k, v))
	}
	defer func() {
		for k := range envs {
			panicIfErr(os.Unsetenv(k))
		}
	}()

	added := map[string]string{
		"ENVGUARD_TEST_SETENVS_FOO":  "aaa",
		"ENVGUARD_TEST_SETENVS_BAR":  "bbb",
		"ENVGUARD_TEST_SETENVS_PIYO": "ccc",
	}
	g, err := Setenvs(added)
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range added {
		e, ok := os.LookupEnv(k)
		if !ok {
			t.Fatal("Env var was not set", k, v)
		}
		if e != v {
			t.Fatal("Env var is unexpected. wanted", v, "but have", e)
		}
	}

	if err := g.Restore(); err != nil {
		t.Fatal(err)
	}

	for k, v := range envs {
		e := os.Getenv(k)
		if e != v {
			t.Error("Environment variable was not restured:", k, v, e)
		}
	}

	if e, ok := os.LookupEnv("ENVGUARD_TEST_SETENVS_PIYO"); ok {
		t.Error("Environment variable was not erased: ", e)
	}
}

func TestAll(t *testing.T) {
	want := os.Environ()
	sort.Strings(want)

	prev := os.Getenv("LANG")
	defer func() {
		os.Setenv("LANG", prev)
	}()

	g := All()

	mod := prev + "-modified"
	g.Setenv("LANG", mod)

	if v := os.Getenv("LANG"); v != mod {
		t.Fatal("Environment variable was not updated:", v, mod)
	}

	if err := g.Restore(); err != nil {
		t.Fatal(err)
	}

	if v := os.Getenv("LANG"); v != prev {
		t.Fatal("Value of $LANG was not restored:", v, prev)
	}

	have := os.Environ()
	sort.Strings(have)

	if len(want) != len(have) {
		t.Fatalf("Number of environment variables does not match. Wanted %#v but have %#v", want, have)
	}

	for i, w := range want {
		h := have[i]
		if w != h {
			t.Fatalf("Value mismatch at index %d: Want %#v but have %#v", i, w, h)
		}
	}
}
