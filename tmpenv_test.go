package tmpenv

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
		"TMPENV_TEST_ADD_FOO": "foo",
		"TMPENV_TEST_ADD_BAR": "",
	}
	defer os.Unsetenv("TMPENV_TEST_ADD_PIYO")

	for k, v := range envs {
		panicIfErr(os.Setenv(k, v))
	}

	g := New()

	// Empty key will be ignored
	g.Add("TMPENV_TEST_ADD_FOO", "TMPENV_TEST_ADD_BAR", "TMPENV_TEST_ADD_PIYO", "")

	if v, ok := os.LookupEnv("TMPENV_TEST_ADD_PIYO"); ok {
		t.Fatal("$TMPENV_TEST_ADD_PIYO should not exist:", v)
	}

	panicIfErr(os.Setenv("TMPENV_TEST_ADD_FOO", "aaa"))
	panicIfErr(os.Setenv("TMPENV_TEST_ADD_BAR", "bbb"))
	panicIfErr(os.Setenv("TMPENV_TEST_ADD_PIYO", "ccc"))

	if err := g.Restore(); err != nil {
		t.Fatal(err)
	}

	for k, v := range envs {
		e := os.Getenv(k)
		if e != v {
			t.Error("Environment variable was not restured:", k, v, e)
		}
	}

	if e, ok := os.LookupEnv("TMPENV_TEST_ADD_PIYO"); ok {
		t.Error("Environment variable was not erased: ", e)
	}
}

func TestSetenv(t *testing.T) {
	envs := map[string]string{
		"TMPENV_TEST_SETENV_FOO": "foo",
		"TMPENV_TEST_SETENV_BAR": "",
	}
	defer os.Unsetenv("TMPENV_TEST_SETENV_PIYO")

	for k, v := range envs {
		panicIfErr(os.Setenv(k, v))
	}

	g := New()

	added := map[string]string{
		"TMPENV_TEST_SETENV_FOO":  "aaa",
		"TMPENV_TEST_SETENV_BAR":  "bbb",
		"TMPENV_TEST_SETENV_PIYO": "ccc",
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

	if e, ok := os.LookupEnv("TMPENV_TEST_SETENV_PIYO"); ok {
		t.Error("Environment variable was not erased: ", e)
	}
}

func TestRemove(t *testing.T) {
	envs := map[string]string{
		"TMPENV_TEST_REMOVE_FOO": "foo",
		"TMPENV_TEST_REMOVE_BAR": "",
	}
	for k, v := range envs {
		panicIfErr(os.Setenv(k, v))
	}
	defer os.Unsetenv("TMPENV_TEST_REMOVE_FOO")
	defer os.Unsetenv("TMPENV_TEST_REMOVE_BAR")

	g := New()
	g.Add("TMPENV_TEST_REMOVE_FOO", "TMPENV_TEST_REMOVE_BAR", "TMPENV_TEST_REMOVE_PIYO")

	panicIfErr(os.Setenv("TMPENV_TEST_REMOVE_FOO", "aaa"))
	panicIfErr(os.Setenv("TMPENV_TEST_REMOVE_BAR", "bbb"))

	if !g.Remove("TMPENV_TEST_REMOVE_BAR", "TMPENV_TEST_REMOVE_PIYO", "TMPENV_TEST_REMOVE_HOGE") {
		t.Fatal("Nothing was removed")
	}
	if g.Remove("TMPENV_TEST_REMOVE_UNKNOWN") {
		t.Fatal("Unknown env var was removed")
	}

	if err := g.Restore(); err != nil {
		t.Fatal(err)
	}

	if v := os.Getenv("TMPENV_TEST_REMOVE_FOO"); v != "foo" {
		t.Error("Env var was not restored", v)
	}
	if v := os.Getenv("TMPENV_TEST_REMOVE_BAR"); v != "bbb" {
		t.Error("Env var was restored", v)
	}
}

func TestSetenvs(t *testing.T) {
	envs := map[string]string{
		"TMPENV_TEST_SETENVS_FOO": "foo",
		"TMPENV_TEST_SETENVS_BAR": "",
	}
	defer os.Unsetenv("TMPENV_TEST_SETENVS_PIYO")

	for k, v := range envs {
		panicIfErr(os.Setenv(k, v))
	}
	defer func() {
		for k := range envs {
			panicIfErr(os.Unsetenv(k))
		}
	}()

	added := map[string]string{
		"TMPENV_TEST_SETENVS_FOO":  "aaa",
		"TMPENV_TEST_SETENVS_BAR":  "bbb",
		"TMPENV_TEST_SETENVS_PIYO": "ccc",
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

	if e, ok := os.LookupEnv("TMPENV_TEST_SETENVS_PIYO"); ok {
		t.Error("Environment variable was not erased: ", e)
	}
}

func TestAll(t *testing.T) {
	panicIfErr(os.Setenv("TMPENV_TEST_ALL_FOO", "foo"))
	defer os.Unsetenv("TMPENV_TEST_ALL_FOO")

	want := os.Environ()
	sort.Strings(want)

	prev := os.Getenv("TMPENV_TEST_ALL_FOO")

	g := All()

	mod := prev + "-modified"
	if err := os.Setenv("TMPENV_TEST_ALL_FOO", mod); err != nil {
		t.Fatal(err)
	}

	if v := os.Getenv("TMPENV_TEST_ALL_FOO"); v != mod {
		t.Fatal("Environment variable was not updated:", v, mod)
	}

	if err := g.Restore(); err != nil {
		t.Fatal(err)
	}

	if v := os.Getenv("TMPENV_TEST_ALL_FOO"); v != prev {
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

func TestNew(t *testing.T) {
	panicIfErr(os.Setenv("TMPENV_TEST_NEW_FOO", "prev"))
	defer os.Unsetenv("TMPENV_TEST_NEW_FOO")

	// Empty key will be ignored
	g := New("TMPENV_TEST_NEW_FOO", "TMPENV_TEST_NEW_BAR", "")

	panicIfErr(os.Setenv("TMPENV_TEST_NEW_FOO", "foo"))
	panicIfErr(os.Setenv("TMPENV_TEST_NEW_BAR", "bar"))

	if err := g.Restore(); err != nil {
		t.Fatal(err)
	}

	if e := os.Getenv("TMPENV_TEST_NEW_FOO"); e != "prev" {
		t.Fatal("env var was not restored:", e)
	}
	if v, ok := os.LookupEnv("TMPENV_TEST_NEW_BAR"); ok {
		t.Fatal("env var which did not exist was not removed", v)
	}
}

func TestSetenvsError(t *testing.T) {
	if _, err := Setenvs(map[string]string{"": "foo"}); err == nil {
		t.Fatal("Error did not occur")
	}
}

func TestRestoreDoesNotCauseErrorOnEmptyKeys(t *testing.T) {
	g := &Envguard{
		maybeMod: map[string]string{"": "foo"},
		maybeAdd: map[string]struct{}{"": struct{}{}},
	}
	if err := g.Restore(); err != nil {
		t.Fatal(err)
	}
}
