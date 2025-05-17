
import { ReturnCode } from "@/lib/wasm";


declare global {
  class Go {
    constructor();
    run(instance: WebAssembly.Instance): void;

    exited: boolean;
    mem: DataView;
    importObject: WebAssembly.Imports;
  }

  // @TODO: rename this one to something like JS_ or CLIENT_ to distinguish that it is not a WASM
  // fn.
  function WASM_listenNotification(buf: Uint8Array): void;

  function WASM_subscribe(buf: Uint8Array, n: number): ReturnCode;
  function WASM_updateSettings(buf: Uint8Array, n: number): ReturnCode;
  function WASM_start(): ReturnCode;
}
