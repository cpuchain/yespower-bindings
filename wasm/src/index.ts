import { bundled } from './bundled.js';
import { base64ToBytes } from './utils.js';
import yespowerWasm, { MainModule } from './yespower_wasm.js';

export * from './utils.js';

type yespower_hash = (
    input: number,
    inputLen: number,
    N: number,
    r: number,
    pers: number,
    persLen: number,
    version: number,
    output: number,
) => number;

const textEncoder = new TextEncoder();

export class Yespower {
    nByte: number;
    Module: MainModule;
    yespower_hash: yespower_hash;
    // Reusable 32-byte output buffer on the wasm heap (allocated once, reused
    // across Hash() calls instead of malloc/free per call).
    outPtr: number;

    constructor(Module: MainModule) {
        this.nByte = 1;
        this.Module = Module;
        // pers is passed as a heap pointer + length ('number'), NOT 'string',
        // so personalities may contain arbitrary / NUL bytes (binary-safe).
        // yespower_hash writes 32 bytes into the caller-provided output pointer
        // and returns 0 on success / -1 on error.
        this.yespower_hash = this.Module.cwrap('yespower_hash', 'number', [
            'number',
            'number',
            'number',
            'number',
            'number',
            'number',
            'number',
            'number',
        ]) as yespower_hash;
        this.outPtr = this.Module._malloc(32);
    }

    static async init() {
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        if (typeof (globalThis as any).WebAssembly === 'undefined') {
            throw new Error('WebAssembly is not enabled with this browser');
        }

        const wasmBinary = base64ToBytes(bundled);

        const module = await yespowerWasm({
            wasmBinary,
            locateFile: (file: string) => file,
        });

        return new Yespower(module);
    }

    // https://stackoverflow.com/questions/41875728/pass-a-javascript-array-as-argument-to-a-webassembly-function
    // Takes an Uint8Array, copies it to the heap and returns a pointer
    arrayToPtr(array: Uint8Array): number {
        const ptr = this.Module._malloc(array.length * this.nByte);
        this.Module.HEAPU8.set(array, ptr / this.nByte);
        return ptr;
    }

    // Takes a pointer and  array length, and returns a Uint8Array from the heap
    ptrToArray(ptr: number, length: number): Uint8Array {
        const array = new Uint8Array(length);
        const pos = ptr / this.nByte;
        array.set(this.Module.HEAPU8.subarray(pos, pos + length));
        return array;
    }

    freePtr(ptr: number) {
        this.Module._free(ptr);
    }

    /**
     * Compute a 32-byte yespower hash.
     *
     * @param input   Data to hash.
     * @param N       Memory cost parameter N (default 2048).
     * @param r       Block size parameter r (default 32).
     * @param pers    Optional personality. A string is UTF-8 encoded; a
     *                Uint8Array is used verbatim (binary-safe).
     * @param version yespower version: 5 (YESPOWER_0_5) or 10 (default).
     */
    Hash(input: Uint8Array, N = 2048, r = 32, pers: string | Uint8Array = '', version = 10): Uint8Array {
        const persBytes = typeof pers === 'string' ? textEncoder.encode(pers) : pers;

        const inputPtr = this.arrayToPtr(input);
        const persPtr = persBytes.length ? this.arrayToPtr(persBytes) : 0;

        // Writes 32 bytes into the reusable output buffer (no per-call malloc).
        const rc = this.yespower_hash(
            inputPtr,
            input.length,
            N,
            r,
            persPtr,
            persBytes.length,
            version,
            this.outPtr,
        );

        this.freePtr(inputPtr);
        if (persPtr) {
            this.freePtr(persPtr);
        }

        if (rc !== 0) {
            throw new Error('yespower failed: invalid parameters (N, r or version)');
        }

        return this.ptrToArray(this.outPtr, 32);
    }
}
