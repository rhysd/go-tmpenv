package tmpenv

import (
	"fmt"
	"os"
)

func Example() {
	// Existing environment variable
	os.Setenv("TMPENV_EXAMPLE_EXISTING", "prev")

	// Create a new environment variable replacements state
	guard := New()

	// Set $TMPENV_EXAMPLE_NEW to "new"
	guard.Setenv("TMPENV_EXAMPLE_NEW", "new")

	// Set $TMPENV_EXAMPLE_EXISTING to "updated"
	guard.Setenv("TMPENV_EXAMPLE_EXISTING", "updated")

	fmt.Println(os.Getenv("TMPENV_EXAMPLE_NEW"))
	// Output:
	// new

	fmt.Println(os.Getenv("TMPENV_EXAMPLE_EXISTING"))
	// Output:
	// new
	// updated

	// Restore state of environment variables. This function is usually called with 'defer'
	if err := guard.Restore(); err != nil {
		panic(err)
	}

	// Environment variables set via guard.Setenv() were restored
	fmt.Println(os.Getenv("TMPENV_EXAMPLE_EXISTING"))
	// Output:
	// new
	// updated
	// prev

	// All environment variables previously did not exist are removed
	_, ok := os.LookupEnv("TMPENV_EXAMPLE_NEW")
	fmt.Println(ok)
	// Output:
	// new
	// updated
	// prev
	// false
}

func ExampleSetenvs() {
	// Existing environment variable
	os.Setenv("TMPENV_EXAMPLE_EXISTING", "prev")

	// Set temporary environment variables by map
	guard, err := Setenvs(map[string]string{
		"TMPENV_EXAMPLE_NEW":      "new",
		"TMPENV_EXAMPLE_EXISTING": "updated",
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("before new: %q\n", os.Getenv("TMPENV_EXAMPLE_NEW"))
	fmt.Printf("before existing: %q\n", os.Getenv("TMPENV_EXAMPLE_EXISTING"))

	// Restore state of environment variables. This function is usually called with 'defer'
	if err := guard.Restore(); err != nil {
		panic(err)
	}

	fmt.Printf("after new: %q\n", os.Getenv("TMPENV_EXAMPLE_NEW"))
	fmt.Printf("after existing: %q\n", os.Getenv("TMPENV_EXAMPLE_EXISTING"))
	// Output:
	// before new: "new"
	// before existing: "updated"
	// after new: ""
	// after existing: "prev"
}

func ExampleAll() {
	// Capture all environment variables
	guard := All()

	// Set $LANG to "C" temporarily with os.Setenv()
	os.Setenv("LANG", "C")

	// Do something awesome...
	fmt.Println(os.Getenv("LANG"))

	// Restore state of environment variables. This function is usually called with 'defer'
	if err := guard.Restore(); err != nil {
		panic(err)
	}

	// $LANG was restored
	fmt.Println(os.Getenv("LANG"))
}

func ExampleEnvguard_Setenv() {
	guard := New()

	// $CONFIG_FOO and $CONFIG_BAR are set to "aaa" and "bbb" temporarily.
	// They will be restored to original values when calling Restore() method.
	guard.Setenv("CONFIG_FOO", "aaa")
	guard.Setenv("CONFIG_BAR", "bbb")

	// Do some tests...
	fmt.Println("foo:", os.Getenv("CONFIG_FOO"))
	fmt.Println("bar:", os.Getenv("CONFIG_BAR"))

	// $CONFIG_FOO and $CONFIG_BAR will be restored. When they were already existing,
	// they will be restored to original values. If they were not existing, they will
	// will be removed. This function is usually called with 'defer'
	guard.Restore()

	// Both variables were restored
	_, fooExists := os.LookupEnv("CONFIG_FOO")
	_, barExists := os.LookupEnv("CONFIG_BAR")
	fmt.Println(fooExists, barExists)
	// Output:
	// foo: aaa
	// bar: bbb
	// false false
}

func ExampleNew() {
	// Some value is already set
	os.Setenv("CONFIG_FOO", "prev-foo")
	os.Setenv("CONFIG_BAR", "prev-bar")

	// Declare to modify $CONFIG_FOO and $CONFIG_BAR temporarily
	guard := New("CONFIG_FOO", "CONFIG_BAR", "CONFIG_PIYO")

	// Modify them
	os.Setenv("CONFIG_FOO", "tmp-foo")
	os.Setenv("CONFIG_BAR", "tmp-bar")
	os.Setenv("CONFIG_PIYO", "tmp-piyo")

	fmt.Println(os.Getenv("CONFIG_FOO"), os.Getenv("CONFIG_BAR"), os.Getenv("CONFIG_PIYO"))

	// Restore previous values
	guard.Restore()

	// Now all values are restored
	fmt.Println(os.Getenv("CONFIG_FOO"), os.Getenv("CONFIG_BAR"))

	// $CONFIG_PIYO was previously not existing. So it does not exist after the restore.
	_, ok := os.LookupEnv("CONFIG_PIYO")
	fmt.Println("$CONFIG_PIYO exists?", ok)
	// Output:
	// tmp-foo tmp-bar tmp-piyo
	// prev-foo prev-bar
	// $CONFIG_PIYO exists? false
}

func ExampleUnset() {
	// Some value is already set
	os.Setenv("CONFIG_FOO", "prev-foo")
	os.Setenv("CONFIG_BAR", "prev-bar")

	// Remember values of $CONFIG_FOO and $CONFIG_BAR and clear them
	guard, err := Unset("CONFIG_FOO", "CONFIG_BAR")
	if err != nil {
		panic(err)
	}

	var ok bool

	fmt.Println("Before")
	_, ok = os.LookupEnv("CONFIG_FOO")
	fmt.Println("  foo exists:", ok)
	_, ok = os.LookupEnv("CONFIG_BAR")
	fmt.Println("  bar exists:", ok)

	// Restore the state. $CONFIG_FOO and $CONFIG_BAR are set to the original values
	guard.Restore()

	fmt.Println("After")
	_, ok = os.LookupEnv("CONFIG_FOO")
	fmt.Println("  foo exists:", ok)
	_, ok = os.LookupEnv("CONFIG_BAR")
	fmt.Println("  bar exists:", ok)
	// Output:
	// Before
	//   foo exists: false
	//   bar exists: false
	// After
	//   foo exists: true
	//   bar exists: true
}

func ExampleUnsetAll() {
	// Remember all environment variables and clear them all
	guard, err := UnsetAll()
	if err != nil {
		panic(err)
	}

	fmt.Println("Environment variables empty?", len(os.Environ()) == 0)

	// All cleared environment variables are restored to the original values
	guard.Restore()

	fmt.Println("Environment variables empty?", len(os.Environ()) == 0)
	// Output:
	// Environment variables empty? true
	// Environment variables empty? false
}
