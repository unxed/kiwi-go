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
## Discrete Cassowary & TUI Grid Hinting

`kiwi-go` includes a **Discrete Cassowary Engine** designed for Terminal User Interfaces (TUI), adapting concepts from font hinting (TrueType & FreeType):

- **FreeType-style Autohinting (`ApportionSum`)**: Distributes rounding remainders using Hamilton / Hare-Niemeyer apportionment so that integer component sizes sum to container totals with zero character gaps or overflows, preserving layout symmetry.
- **TrueType-style Rule Directives (`RuleHinter`)**: Apply explicit layout instructions such as `SnapToGrid` (for double-width CJK or grid boundaries), `ClampMinMax`, `EqualizeGroup`, and `AlignEdges`.
- **`DiscreteSolver`**: Integrates continuous Cassowary constraint solving, autohinting apportionment, and rule directives in a unified TUI layout solver.

For full mathematical design and TUI layout examples, see [docs/DISCRETE_CASSOWARY.md](docs/DISCRETE_CASSOWARY.md).

## Console Demo

Run the interactive TUI console demo:
```sh
go run ./cmd/kiwi-demo
```

## Testing & Benchmarks

Run unit tests:
```sh
go test -v ./...
```

Run benchmarks:
```sh
go test -bench=. ./...
```

### Test in NodeJS

First, ensure you have generated `kiwi.wasm` and `wasm_exec.js` using `./build.sh`.

To run the NodeJS integration test, execute:

```sh
node test.mjs
```

You can also run the internal Go bridge tests using standard Go tools:
```sh
GOOS=js GOARCH=wasm go test -v .
```

### Web Demo

To see the engine solving UI layouts in your browser:
1. Start a local web server in this directory (browsers do not allow fetching WASM from local `file://` URLs due to CORS). For example:
   ```sh
   python3 -m http.server 8080
   ```
2. Open [http://localhost:8080/index.html](http://localhost:8080/index.html) in your browser.
3. Resize the gray container box to see Cassowary dynamically calculate and apply layout constraints in real-time.
```

## License

Modified BSD License.
