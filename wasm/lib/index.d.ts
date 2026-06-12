declare namespace RuntimeExports {
	/**
	 * @param {string|null=} returnType
	 * @param {Array=} argTypes
	 * @param {Array=} args
	 * @param {Object=} opts
	 */
	function ccall(ident: any, returnType?: (string | null) | undefined, argTypes?: any[] | undefined, args?: any[] | undefined, opts?: Object | undefined): any;
	/**
	 * @param {string=} returnType
	 * @param {Array=} argTypes
	 * @param {Object=} opts
	 */
	function cwrap(ident: any, returnType?: string | undefined, argTypes?: any[] | undefined, opts?: Object | undefined): any;
	let HEAPU8: Uint8Array;
}
export interface WasmModule {
	_yespower_hash(_0: number, _1: number, _2: number, _3: number, _4: number, _5: number, _6: number, _7: number): number;
	_yespower_wasm(_0: number, _1: number, _2: number, _3: number, _4: number, _5: number, _6: number): number;
	_malloc(_0: number): number;
	_free(_0: number): void;
}
export type MainModule = WasmModule & typeof RuntimeExports;
export declare function bytesToBase64(bytes: Uint8Array): string;
export declare function base64ToBytes(base64: string): Uint8Array<ArrayBuffer>;
export declare function bytesToHex(bytes: Uint8Array): string;
export declare function hexToBytes(hexStr: string): Uint8Array<ArrayBuffer>;
export type yespower_hash = (input: number, inputLen: number, N: number, r: number, pers: number, persLen: number, version: number, output: number) => number;
export declare class Yespower {
	nByte: number;
	Module: MainModule;
	yespower_hash: yespower_hash;
	outPtr: number;
	constructor(Module: MainModule);
	static init(): Promise<Yespower>;
	arrayToPtr(array: Uint8Array): number;
	ptrToArray(ptr: number, length: number): Uint8Array;
	freePtr(ptr: number): void;
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
	Hash(input: Uint8Array, N?: number, r?: number, pers?: string | Uint8Array, version?: number): Uint8Array;
}

export {};
