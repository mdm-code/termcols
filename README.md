<h1 align="center">
  <div>
    <img src="https://raw.githubusercontent.com/mdm-code/mdm-code.github.io/main/termcols_logo.png" alt="logo"/>
  </div>
</h1>

<h4 align="center">Colorful text in your terminal implemented in Go</h4>

<div align="center">
<p>
    <a href="https://github.com/mdm-code/termcols/actions?query=workflow%3ACI">
        <img alt="Build status" src="https://github.com/mdm-code/termcols/workflows/CI/badge.svg">
    </a>
    <a href="https://app.codecov.io/gh/mdm-code/termcols">
        <img alt="Code coverage" src="https://codecov.io/gh/mdm-code/termcols/branch/main/graphs/badge.svg?branch=main">
    </a>
    <a href="https://opensource.org/licenses/MIT" rel="nofollow">
        <img alt="MIT license" src="https://img.shields.io/github/license/mdm-code/termcols">
    </a>
    <a href="https://goreportcard.com/report/github.com/mdm-code/termcols">
        <img alt="Go report card" src="https://goreportcard.com/badge/github.com/mdm-code/termcols">
    </a>
    <a href="https://pkg.go.dev/github.com/mdm-code/termcols">
        <img alt="Go package docs" src="https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white">
    </a>
</p>
</div>

The `termcols` package implements ANSI color codes that can be used to color
text on the terminal. Different styles and foreground/background colors can be
chained together through an intuitive package API to arrive at some cool visual
effects.

The selection of style and color control sequences implemented by the package
was largely based on an exhaustive list of Select Graphic Rendition (SGR)
control sequences available at [Wikipedia ANSI](https://en.wikipedia.org/wiki/ANSI_escape_code).
It is a great resource in case one or more elements appear not to be supported
in a given terminal.

The Escape sequence for `termcols` is set to `\033`, which means that it should
work without any issues with Bash, Zsh or Dash. Other shells might not support
it.

The same applies to 8-bit and 24-bit colors: there is no guarantee that these
escape sequences are supported will be rendered properly on some terminals.
Results may vary, so it is good practice to test it first for compatibility.

Consult the [package documentation](https://pkg.go.dev/github.com/mdm-code/termcols)
or see [Usage](#usage) section below to check how to use the public API of the
`termcols` package.


## Installation

Use the following command to add the package to an existing project.

```sh
go get github.com/mdm-code/termcols
```

Install the package to use the command-line `tcols` command to colorize
text on the terminal.

```sh
go install github.com/mdm-code/termcols@latest
```


## Usage

Here is an example of how to use the public API of the `termcols` package.

```go
package main

import (
	"fmt"

	"github.com/mdm-code/termcols"
)

func main() {
	s := termcols.Colorize(
		"Colorized text!",
		termcols.RedFg,
		termcols.Underline,
		termcols.Rgb24(termcols.BG, 120, 255, 54),
	)
	fmt.Println(s)
}
```

Aside from using the `termcols` package API that can be used in your Go
project, can use the `tcols` terminal command:

```sh
tcols --style 'redfg underline rgb24=bg:120:255:54' < <(echo -n 'Hello, world!')
```

Type `tcols -h` to get a list of styles and colors to (1) see what is implemented
and (2) what is supported by your terminal.


Alternatively, `tcols` can be run from inside of the Docker container:

```sh
docker run -i ghcr.io/mdm-code/tcols:latest tcols -s 'redfg bluebg' < <(echo -n 'Hello, world!')
```


## Development

Consult [Makefile](Makefile) to see how to format, examine code with `go vet`,
run unit test, run code linter with `golint` get test coverage and check if the
package builds all right.

Remember to install `golint` before you try to run tests and test the build:

```sh
go install golang.org/x/lint/golint@latest
```

In order to run the benchmark test on unsafe pointers in the tercmols package,
fire up the following command:

```sh
go test -bench=.
```

This will give you ns/op value for the setup it's been benchmarked on.


## License

Copyright (c) 2024 Michał Adamczyk.

This project is licensed under the [MIT license](https://opensource.org/licenses/MIT).
See [LICENSE](LICENSE) for more details.
