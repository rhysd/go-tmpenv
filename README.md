Go package `tmpenv`
===================
[![Linux/Mac Build Status][travisci-badge]][travisci]
[![Windows Build Status][appveyor-badge]][appveyor]
[![Coverage Report][codecov-badge]][codecov]
[![Documentation][doc-badge]][doc]

[`tmpenv`](doc) is a small Go package for replacing environment variables temporarily.

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
	guard := tmpenv.All()

	// Ensure to restore the captured variables with 'defer'
	defer guard.Restore()

	// Modify existing environment variable
	guard.Setenv("LANG", "C")

	// Set new environment variable
	guard.Setenv("FOO_ANSWER", "42")

	// Do some tests...

	// $LANG will be restored to original value and $FOO_ANSWER will be removed.
	// And other all environment variables will be restored to original values.
}
```

If you're interested in specific environment variables:

```go
import (
	"github.com/rhysd/go-tmpenv"
	"testing"
)

func TestFoo(t *testing.T) {
	// $CONFIG_FOO and $CONFIG_BAR are set to "aaa" and "bbb" temporarily.
	// They will be restored to original values when calling Restore() method.
	guard, err := tmpenv.Setenvs(map[string]string{
		"CONFIG_FOO": "aaa",
		"CONFIG_BAR": "bbb",
	})
	if err != nil {
		panic(err)
	}

	// Ensure to restore the captured variables with 'defer'
	defer guard.Restore()

	// Do some tests...

	// $CONFIG_FOO and $CONFIG_BAR will be restored. When they were already existing,
	// they will be restored to original values. If they were not existing, they will
	// will be removed.
}
```

Following is the same as above, but adding variables one by one:

```go
import (
	"github.com/rhysd/go-tmpenv"
	"testing"
)

func TestFoo(t *testing.T) {
	guard := tmpenv.New()

	// Ensure to restore the captured variables with 'defer'
	defer guard.Restore()

	// $CONFIG_FOO and $CONFIG_BAR are set to "aaa" and "bbb" temporarily.
	// They will be restored to original values when calling Restore() method.
	guard.Setenv("CONFIG_FOO", "aaa")
	guard.Setenv("CONFIG_BAR", "bbb")

	// Do some tests...

	// $CONFIG_FOO and $CONFIG_BAR will be restored. When they were already existing,
	// they will be restored to original values. If they were not existing, they will
	// will be removed.
}
```

Please read [the documentation][doc] for more details.

## Repository

This library is developed at [GitHub repository](https://github.com/rhysd/go-tmpenv). If you're facing
some error or have some feature request, please create a new issue at the page.

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
