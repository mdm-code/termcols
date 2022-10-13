Package termcols implements ANSI color codes that can be used to color text on
the terminal. Different styles and foreground/background colors can be chained
together through an intuitive package API to arrive at some cool visual
effects.

The selection of style and color control sequences implemented by the package
was largely based on an exhaustive list of Select Graphic Rendition (SGR)
control sequences available at [Wikipedia ANSI](https://en.wikipedia.org/wiki/ANSI_escape_code). It is a great resource
in case one or more elements appear not to be supported in a given terminal.

The Escape sequence for termcols is set to \033, which means that it should
work without any issues with Bash, Zsh or Dash. Other shells might not support
it.

The same applies to 8-bit and 24-bit colors: there is no guarantee that these
escape sequences are supported will be rendered properly on some terminals.
Results may vary, so it is good practice to test it first for compatibility.

See [Usage](#usage) section below to check how to use the public API of the
package.


## Installation

```sh
go get github.com/mdm-code/termcols
```


## Usage

Here is an example of how to use the public API of the `termcols` package.

```go
package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, world!")
}
```


## Development

Consult [Makefile](Makefile) to see how to format, examine code with `go vet`,
run unit test, run code linter with `golint` get test coverage and check if the
package builds all right.

Remember to install `golint` before you try to run tests and test the build:

```sh
go install golang.org/x/lint/golint@latest
```


## License

Copyright (c) 2022 Michał Adamczyk.

This project is licensed under the [MIT license](https://opensource.org/licenses/MIT).
See [LICENSE](LICENSE) for more details.
