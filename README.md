Go package `tmpenv`
===================
[![Linux/Mac Build Status][travisci-badge]][travisci]
[![Windows Build Status][appveyor-badge]][appveyor]
[![Coverage Report][codecov-badge]][codecov]
[![Documentation][doc-badge]][doc]

`tmpenv` is a small Go package for replacing environment variables temporarily.

Setting temporary environment variables and restoring them are common pattern on testing.
However, they have a pitfall since setting empty value to environment variable is different from
unsetting environment variable.

`tmpenv` helps you to do the pattern correctly with less lines.

## Installation

```
$ go get -u github.com/rhysd/go-tmpenv
```

## Usage

To ensure to restore all existing environment variables and to remove all temporary environment variables,
`tmpvar.All()` is useful.

```go
import (
	"github.com/rhysd/go-tmpenv"
	"testing"
)

func TestFoo(t *testing.T) {
	// Captures all environment variables
	g := tmpenv.All()

	// Ensure to restore the captured variables with 'defer'
	defer g.Restore()

	// Modify existing environment variable
	g.Setenv("LANG", "C")

	// Set new environment variable
	g.Setenv("FOO_ANSWER", "42")

	// Do some tests...

	// $LANG will be restored to original value and $FOO_ANSWER will be removed.
	// And other all environment variables will be restored to original values.
}
```

If you're interested in specific variables:

```go
import (
	"github.com/rhysd/go-tmpenv"
	"testing"
)

func TestFoo(t *testing.T) {
	// $CONFIG_FOO and $CONFIG_BAR are set to "aaa" and "bbb" temporarily.
	// They will be restored to original values when calling Restore() method.
	g, err := tmpenv.Setenvs(map[string]string{
		"CONFIG_FOO": "aaa",
		"CONFIG_BAR": "bbb",
	})
	if err != nil {
		panic(err)
	}

	// Ensure to restore the captured variables with 'defer'
	defer g.Restore()

	// Do some tests...

	// $CONFIG_FOO and $CONFIG_BAR will be restored. When they were already existing,
	// they will be restored to original values. If they were not existing, they will
	// will be removed.
}
```

## License

[MIT License](LICENSE.txt)


[doc-badge]: https://godoc.org/github.com/rhysd/go-tmpenv?status.svg
[doc]: http://godoc.org/github.com/rhysd/go-tmpenv
[travisci-badge]: https://travis-ci.org/rhysd/go-tmpenv.svg?branch=master
[travisci]: https://travis-ci.org/rhysd/go-tmpenv
[appveyor-badge]: https://ci.appveyor.com/api/projects/status/5pbcku1buw8gnqu9/branch/master?svg=true
[appveyor]: https://ci.appveyor.com/project/rhysd/go-tmpenv
[codecov-badge]: https://codecov.io/gh/rhysd/go-tmpenv/branch/master/graph/badge.svg
[codecov]: https://codecov.io/gh/rhysd/go-tmpenv
