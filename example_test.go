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

	// Set $LANG to "C" temporarily
	guard.Setenv("LANG", "C")

	// Do something awesome...
	fmt.Println(os.Getenv("LANG"))

	// Restore state of environment variables. This function is usually called with 'defer'
	if err := guard.Restore(); err != nil {
		panic(err)
	}

	// $LANG was restored
	fmt.Println(os.Getenv("LANG"))
}
