# Yespower WebAssembly Module

[![NPM Version](https://img.shields.io/npm/v/yespower-wasm)](https://www.npmjs.com/package/yespower-wasm)

Yespower WebAssembly Module for Browsers / Node.js

## Prerequisites

* Node.js LTS or Browser where WebAssembly is enabled

## Install

```bash
$ yarn add yespower-wasm
```

or on html header / body

```html
<script src="https://cdn.jsdelivr.net/npm/yespower-wasm/lib/yespower.umd.min.js"></script>
```

WASM file is fully embedded on script so that you need to load nothing (and should work with any bundlers without hassle)

## Build WASM (optional)

Requires latest [Emscripten](https://emscripten.org/docs/getting_started/downloads.html) WASM compiler

```bash
$ yarn && yarn build:wasm
```

Will update .wasm and bundled files

## Build library (optional)

Rebuild node.js and web umd bundle files

```bash
$ yarn && yarn build

```

## API

```ts
const yespower = await Yespower.init();

yespower.Hash(
    input: Uint8Array,
    N?: number,                  // memory cost N (default 2048)
    r?: number,                  // block size r (default 32)
    pers?: string | Uint8Array,  // personality (default '')
    version?: number,            // 5 = YESPOWER_0_5, 10 = YESPOWER_1_0 (default)
): Uint8Array                    // 32-byte hash
```

* `pers` accepts a `string` (UTF-8 encoded) or a `Uint8Array` (used verbatim,
  so it is binary-safe and may contain NUL bytes).
* `version` defaults to `10` (YESPOWER_1_0). Pass `5` for YESPOWER_0_5.

## Example Code

```js
const { Yespower, bytesToHex, hexToBytes } = require('yespower-wasm');

async function test() {
    const yespower = await Yespower.init();

    const input = hexToBytes('2538dad623dab29a3d5387804ab51dea411014fad9c47fb94b2d83e44358064b');

    // YESPOWER_1_0, N=2048, r=32, personality "pers"
    console.log(bytesToHex(yespower.Hash(input, 2048, 32, 'pers')));

    // YESPOWER_0_5 with N=2048, r=8 and a personality
    console.log(bytesToHex(yespower.Hash(input, 2048, 8, 'Client Key', 5)));
}

test();
```

Also refer `yespower.html` for example of using it on browser.