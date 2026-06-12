Yespower for Node.js
============================

[![NPM Version](https://img.shields.io/npm/v/yespower)](https://www.npmjs.com/package/yespower)

## Prerequisites

* Node.js LTS

## Install

```bash
$ yarn add yespower
```

## Build

```bash
$ yarn
```

## API

```ts
yespower(
    input: Buffer | Uint8Array,
    n?: number,                          // memory cost N (default 2048)
    r?: number,                          // block size r (default 32)
    pers?: string | Buffer | Uint8Array, // personality (default none)
    version?: 5 | 10,                    // 5 = YESPOWER_0_5, 10 = YESPOWER_1_0 (default)
): Buffer                                // 32-byte hash
```

* `pers` accepts a `string` (UTF-8 encoded) or a `Buffer` / `Uint8Array`
  (used verbatim, so it is binary-safe and may contain NUL bytes).
* `version` selects the algorithm variant; it defaults to `10` (YESPOWER_1_0)
  to preserve backward compatibility. Pass `5` for YESPOWER_0_5.

### Asynchronous

`yespower_async` takes the same arguments and returns a `Promise<Buffer>`. It
runs the (CPU-heavy) hash on the libuv threadpool so the Node.js event loop is
not blocked, which lets you hash many inputs concurrently (e.g. with
`Promise.all`). Inputs are copied, so the source `Buffer` may be reused or
mutated immediately after the call. The Promise rejects on invalid parameters.

```ts
yespower_async(
    input: Buffer | Uint8Array,
    n?: number,
    r?: number,
    pers?: string | Buffer | Uint8Array,
    version?: 5 | 10,
): Promise<Buffer>
```

## Example Code

```js
import { yespower, yespower_async } from 'yespower';

const data = Buffer.from("7000000001e980924e4e1109230383e66d62945ff8e749903bea4336755c00000000000051928aff1b4d72416173a8c3948159a09a73ac3bb556aa6bfbcad1a85da7f4c1d13350531e24031b939b9e2b", "hex");

// YESPOWER_1_0, N=2048, r=32 (defaults)
console.log(yespower(data).toString('hex'));

// YESPOWER_0_5 with a personality (e.g. "Client Key")
console.log(yespower(data, 2048, 8, 'Client Key', 5).toString('hex'));

// Asynchronous (non-blocking) — hash several inputs concurrently
const hashes = await Promise.all([
    yespower_async(data),
    yespower_async(data, 2048, 8, 'Client Key', 5),
]);
console.log(hashes.map((h) => h.toString('hex')));
```
