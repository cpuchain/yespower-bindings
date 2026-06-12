// TypeScript bindings for emscripten-generated code.  Automatically generated at compile time.
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
interface WasmModule {
  _yespower_hash(_0: number, _1: number, _2: number, _3: number, _4: number, _5: number, _6: number, _7: number): number;
  _yespower_wasm(_0: number, _1: number, _2: number, _3: number, _4: number, _5: number, _6: number): number;
  _malloc(_0: number): number;
  _free(_0: number): void;
}

export type MainModule = WasmModule & typeof RuntimeExports;
export default function MainModuleFactory (options?: unknown): Promise<MainModule>;
