# expect ![test](https://github.com/crumbandbase/expect/workflows/test/badge.svg?event=push)

A (very) simple expectation, and `diff`, package for use in test suites.

## Prerequisites

You will need the following things properly installed on your computer.

* [Go](https://golang.org/): any one of the **three latest major**
  [releases](https://golang.org/doc/devel/release.html)

## Installation

With [Go module](https://github.com/golang/go/wiki/Modules) support (Go 1.11+),
simply add the following import

```go
import "github.com/crumbandbase/expect"
```

to your code, and then `go [build|run|test]` will automatically fetch the
necessary dependencies.

Otherwise, to install the `expect` package, run the following command:

```bash
$ go get -u github.com/crumbandbase/expect
```

## Usage

In the simple case where two like values need to be tested for equality the
`expect.Equal` function can be used.

```go
package main_test

import (
  "testing"

  "github.com/crumbandbase/expect"
)

func TestGreeting(t *testing.T) {
	t.Run("succeeds when the greetings are equal", func(t *testing.T) {
		expected := "Hello, Picard"
		got := greeting("Picard")

		expect.Equal(t, got, expected)
	})

	t.Run("fails when the greeting are not equal", func(t *testing.T) {
		expected := "Bonjour, Picard"
		got := greeting("Picard")

		expect.Equal(t, got, expected)
	})
}
```

In the above the second test will fail and will produce a `diff` describing the
differences. This is handled by the brilliant
[go-cmp](https://github.com/google/go-cmp) package. For example the output
that will be printed for this test suite is:

```bash
    string(
  -       "Hello, Picard",
  +       "Bonjour, Picard",
    )
```

## License

This project is licensed under the [MIT License](LICENSE.md).
