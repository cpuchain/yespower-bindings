# go-yespower

Go bindings for the [yespower](https://www.openwall.com/yespower) proof-of-work
hashing algorithm.

Two independent implementations live in the same package:

* **cgo** — `Hash` / `HashVersion` compile and call the reference C sources in
  `../yespower-c` directly. Fast; requires a C toolchain and `CGO_ENABLED=1`.
* **pure Go** — `YespowerNative` / `YescryptNative` are a dependency-light port
  (no cgo). Roughly 5–10× slower, but build anywhere without a C compiler.

Both produce identical 32-byte digests for the same parameters.

## Install

```bash
$ go get github.com/cpuchain/go-yespower
```

The cgo functions need a C compiler (e.g. `gcc`/`clang`) and `CGO_ENABLED=1`
(the default when a compiler is present). The pure-Go functions work even with
`CGO_ENABLED=0`.

### Windows

On Windows the cgo path needs MSYS2's UCRT64 gcc. The `build.ps1` helper
locates it automatically (checking `$env:UCRT64`, `$env:MSYS2_ROOT`,
`C:\msys64`, then the registry) and runs the tests and/or example:

```powershell
.\build.ps1            # go test -v ./...  (default)
.\build.ps1 -Example   # go run ./example
.\build.ps1 -Test -Example
```

Install the toolchain first with `pacman -S mingw-w64-ucrt-x86_64-gcc`.

## API

### cgo (reference C implementation)

```go
// Version constants for HashVersion (match yespower_version_t).
const (
    Version05 uint32 = 5  // YESPOWER_0_5
    Version10 uint32 = 10 // YESPOWER_1_0
)

// HashVersion computes the 32-byte yespower hash with an explicit version.
func HashVersion(input []byte, N uint32, r uint32, per string, version uint32) []byte

// Hash computes the YESPOWER_1_0 hash (version 10).
func Hash(input []byte, N uint32, r uint32, per string) []byte
```

`per` is the personalization token, passed as a length-delimited buffer
(binary-safe; an empty string means no personalization).

### Pure Go (no cgo)

```go
// YespowerNative computes the YESPOWER_1_0 hash.
func YespowerNative(in []byte, N int, r int, persToken string) []byte

// YescryptNative computes the YESPOWER_0_5 hash.
func YescryptNative(in []byte, N int, r int, persToken string) []byte
```

## Example

```go
package main

import (
    "encoding/hex"
    "fmt"

    "github.com/cpuchain/go-yespower"
)

func main() {
    in, _ := hex.DecodeString("eebb7bf9a8c813b5e0a03ce627bd1a0c836e0a89793743666dc82b83e28e8f00")

    // cgo, YESPOWER_1_0, N=2048, r=32, no personalization
    fmt.Println(hex.EncodeToString(yespower.Hash(in, 2048, 32, "")))

    // cgo, YESPOWER_0_5 with a personalization token
    fmt.Println(hex.EncodeToString(yespower.HashVersion(in, 2048, 8, "Client Key", yespower.Version05)))

    // pure Go (no cgo), YESPOWER_1_0
    fmt.Println(hex.EncodeToString(yespower.YespowerNative(in, 2048, 32, "")))
}
```

Run the bundled example:

```bash
$ go run ./example
```

## Testing

```bash
$ go test -v ./...

# Benchmarks: ns/op -> s/op -> h/s
$ go test -bench=Yes -benchtime 4s
```

The cgo path is checked against every vector in `../yespower-c/TESTS-OK`
(including the binary-`pers` `BSTY` case and the XOR-aggregate lines), and the
pure-Go path has its own matching tests.

## LICENSE

BSD 2-Clause, as per the yespower source files.
