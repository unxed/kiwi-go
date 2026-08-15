# kiwi-go

Fast Go implementation of the [Cassowary constraint solving algorithm](https://en.wikipedia.org/wiki/Cassowary_(software)), ported directly from [LUME Kiwi](https://github.com/lume/kiwi).

## Install

```sh
go get github.com/unxed/kiwi-go
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/unxed/kiwi-go"
)

func main() {
	solver := kiwi.NewSolver()

	left := kiwi.NewVariable("left")
	width := kiwi.NewVariable("width")
	right := kiwi.NewVariable("right")

	_ = solver.AddEditVariable(left, kiwi.StrengthStrong)
	_ = solver.AddEditVariable(width, kiwi.StrengthStrong)

	_ = solver.SuggestValue(left, 100)
	_ = solver.SuggestValue(width, 400)

	// right == left + width
	cn := kiwi.NewConstraint(right, kiwi.OpEq, left.Plus(width))
	_ = solver.AddConstraint(cn)

	solver.UpdateVariables()

	fmt.Printf("right = %f\n", right.Value()) // 500.0
}
```

## Tests & Benchmarks

Run unit tests:
```sh
go test -v ./...
```

Run benchmarks:
```sh
go test -bench=. ./...
```

## License

Modified BSD License.
