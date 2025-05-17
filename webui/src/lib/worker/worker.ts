import type {
  GenericMessage,
  LoadWASM,
  WorkerEvent,
  Subscribe,
  UpdateSettings,
  StartWatcher,
} from '@/lib/messages';
import { MessageType, Event } from '@/lib/messages';
import { loadWASM, serializeStr, type WASMLoadResult } from '@/lib/wasm';

// ensure wasm glue code for Go is registered in the worker.
import '@/lib/worker/wasm_exec.js';

export function initialize() {
  registerHandlers();
  notifyStateChange(Event.WorkerLoaded);
}


// WorkerInstance handles incoming messages from the main thread and
// invokes methods of the WASM binary.
class WorkerInstance {
  result?: WASMLoadResult;

  constructor() {
    self.WASM_listenNotification = this.handleNotification;
  }

  handleMessage(msg: GenericMessage) {
    switch (msg.type) {
      case MessageType.LoadWASM:
        return this.handleLoadWASM(msg as LoadWASM);

      case MessageType.Subscribe:
        return this.handleSubscribe(msg as Subscribe);

      case MessageType.UpdateSettings:
        return this.handleUpdateSettings(msg as UpdateSettings);

      case MessageType.StartWatcher:
        return this.handleStartWatcher(msg as StartWatcher);

      default:
        console.log('unhandled message: ', msg);
    }
  }

  handleLoadWASM(msg: LoadWASM): void {
    loadWASM(msg.data.filename).then((result) => {
      this.result = result;
      notifyStateChange(Event.WASMLoaded);
    });
  }

  handleSubscribe(msg: Subscribe): void {
    const { buf, n } = serializeStr(msg.data.address);
    console.log('[wasm subscribe] code: ', self.WASM_subscribe(buf, n));
  }

  handleNotification(buf: Uint8Array): void {
    const decoder = new TextDecoder();
    const s = decoder.decode(buf);
    const data = JSON.parse(s);

    self.postMessage({
      type: MessageType.TxNotification,
      data: { tx: data },
    });
  }

  handleUpdateSettings(msg: UpdateSettings): void {
    const stringified = JSON.stringify(msg.data);
    const { buf, n } = serializeStr(stringified);

    console.log('[wasm us]: ', self.WASM_updateSettings(buf, n));
  }

  handleStartWatcher(msg: StartWatcher): void {
    console.log('[wasm start]: ', self.WASM_start());
  }
}

function registerHandlers() {
  console.log('registering handlers...');
  const worker = new WorkerInstance();

  self.addEventListener('message', (event) => {
    worker.handleMessage(event.data as GenericMessage);
  });

  notifyStateChange(Event.WorkerLoaded, 100);
}

function notifyStateChange(state: Event, delay?: number): void {
  const message: WorkerEvent = {
    type: MessageType.WorkerEvent,
    data: {
      state,
    },
  };

  if (typeof delay !== 'undefined') {
    setTimeout(() => self.postMessage(message), delay);
    return;
  }

  self.postMessage(message);
}
