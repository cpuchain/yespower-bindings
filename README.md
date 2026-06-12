# Node.js / WASM / Go bindings for Yespower

[![NPM Version](https://img.shields.io/npm/v/yespower)](https://www.npmjs.com/package/yespower)
[![NPM Version](https://img.shields.io/npm/v/yespower-wasm)](https://www.npmjs.com/package/yespower-wasm)

For node.js bindings refer to `./node` directory (NAPI addon),

For WebAssembly bindings refer to `./wasm` directory (Emscripten),

For Go bindings refer to `./go` directory (cgo + a pure-Go fallback)

## API parity with the reference C implementation

The Node.js and WASM bindings wrap the reference `yespower_tls()` /
`yespower()` API in `./yespower-c` (`yespower.h`). The thin C shim in
`./yespower-c/yespower.c` (`yespower_hash` / `yespower_wasm`) exposes the full
parameter set. The Go binding compiles the same `./yespower-c` sources directly
via cgo (`Hash` / `HashVersion`) and additionally ships a dependency-light pure
-Go port (`YespowerNative` / `YescryptNative`) for builds without a C compiler.

| C param (`yespower_params_t`) | binding argument |
| ----------------------------- | ---------------- |
| `version` (5 / 10)            | `version` (default `10`) |
| `N`                           | `N` (default `2048`) |
| `r`                           | `r` (default `32`) |
| `pers` + `perslen`            | `pers` (string or bytes, length-delimited) |

Every binding's test suite reproduces the vectors in `./yespower-c/TESTS-OK` —
the YESPOWER_0_5 and YESPOWER_1_0 single hashes (including the binary-`pers`
`BSTY` case) and the two XOR-aggregate lines — by replicating the input
construction of the reference driver `./yespower-c/tests.c`.

Run the reference C test yourself with:

```bash
cd yespower-c && make check   # builds tests and diffs against TESTS-OK
```

## Wrapper improvements

Implemented:

* **Full `version` support** — previously the wrapper hard-coded
  `YESPOWER_1_0`, so the YESPOWER_0_5 vectors were unreachable.
* **Error propagation** — `yespower()`'s return code is now checked; the
  bindings throw on failure instead of returning garbage.
* **Binary-safe personality** — `pers` is passed as a length-delimited byte
  buffer (a heap pointer + length in WASM, a `Buffer`/`Uint8Array` or UTF-8
  string in NAPI) instead of a NUL-terminated C string, so personalities with
  embedded NUL or non-ASCII bytes work.
* **No per-call leaks (NAPI)** — the addon writes into a stack buffer and lets
  Node own the result `Buffer`; the previous `malloc`'d output and `strcpy`'d
  personality were never freed.
* **Reused buffers (WASM + C)** — the WASM `Hash()` reuses a single 32-byte
  heap output buffer across calls (via the `yespower_hash` out-param export),
  and the C shim keeps a thread-local `yespower_local_t` so the memory arena is
  initialised once and reused (mirrors the reference `yespower()` usage).
* **Async / non-blocking NAPI** — `yespower_async` runs the hash on the libuv
  threadpool via `Napi::AsyncWorker` and returns a `Promise`, so the event loop
  stays free during heavy `N`/`r` hashing and inputs can be hashed concurrently.

Future work (not yet implemented):

* **Typed-array (zero-copy) returns** instead of copying into a `Buffer`.
* **Explicit `yespower_init_local` / `yespower_free_local` lifecycle** exposed
  to callers managing their own per-thread arenas.

## References

- [openwall/yespower](https://github.com/openwall/yespower) Reference implementation itself

- [bellcoin-electrum/node-bellcoin-yespower](https://github.com/bellcoin-electrum/node-bellcoin-yespower) Another Node.js bindings

- [bellcoin-electrum/bell_yespower_python3](https://github.com/bellcoin-electrum/bell_yespower_python3) Python3 bindings

- [mraksoll4/yespower_go](https://github.com/mraksoll4/yespower_go) Golang bindings
