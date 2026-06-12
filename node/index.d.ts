/**
 * Compute a 32-byte yespower hash.
 *
 * @param input   Data to hash.
 * @param n       Memory cost parameter N (default 2048).
 * @param r       Block size parameter r (default 32).
 * @param pers    Optional personality. A string is UTF-8 encoded; a Buffer /
 *                Uint8Array is used verbatim (binary-safe, may contain NUL bytes).
 * @param version yespower version: 5 (YESPOWER_0_5) or 10 (YESPOWER_1_0, default).
 * @returns       32-byte hash as a Buffer.
 */
export function yespower(
    input: Buffer | Uint8Array,
    n?: number,
    r?: number,
    pers?: string | Buffer | Uint8Array,
    version?: 5 | 10,
): Buffer;

/**
 * Asynchronous variant of {@link yespower}. Runs the CPU-heavy hash on the
 * libuv threadpool and resolves a Promise, so the Node.js event loop is not
 * blocked. The input and personality are copied, so they may be reused or
 * mutated immediately after the call. The returned Promise rejects on invalid
 * parameters (e.g. an unsupported `version`).
 */
export function yespower_async(
    input: Buffer | Uint8Array,
    n?: number,
    r?: number,
    pers?: string | Buffer | Uint8Array,
    version?: 5 | 10,
): Promise<Buffer>;
