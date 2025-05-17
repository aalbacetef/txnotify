import type {
  GenericMessage,
  LoadWASM,
  WorkerEvent,
  Subscribe,
  TxReceived,
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

  handleNotification(buf: Uint8Array) {
    const decoder = new TextDecoder();
    const s = decoder.decode(buf);
    const data = JSON.parse(s);

    console.log('data: ', data);

    self.postMessage({
      type: MessageType.TxReceived,
      data: { tx: data },
    });
  }

  handleMessage(msg: GenericMessage) {
    console.log('got message: ', msg);

    switch (msg.type) {
      case MessageType.LoadWASM:
        return this.handleLoadWASM(msg as LoadWASM);

      case MessageType.Subscribe:
        return this.handleSubscribe(msg as Subscribe);

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
    self.WASM_subscribe(buf, n);
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
