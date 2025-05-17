
type Settings = Record<string, number | string | boolean | object>;

declare global {
  class Go {
    constructor();
    run(instance: WebAssembly.Instance): void;

    exited: boolean;
    mem: DataView;
    importObject: WebAssembly.Imports;
  }

  function WASM_subscribe(buf: Uint8Array, n: number);
  function WASM_listenNotification(buf: Uint8Array);
  function WASM_settingsUpdated(settings: Settings);
  function WASM_start();
}
