import {
  MessageType,
  Event,
  type Subscribe,
  type TxReceived,
} from '@/lib/messages';

import {
  type GenericMessage,
  type LoadWASM,
  type WorkerEvent,
} from '@/lib/messages';


type StateChangeCB = (state: Event) => void;

// WorkerPeer provides a set of methods to interact with the Worker.
export class WorkerPeer {
  store = null;
  worker: Worker;
  callbacks: {
    [key in Event]?: StateChangeCB[];
  } = {};

  oneoffs: {
    [key in Event]?: StateChangeCB[];
  } = {};

  constructor(worker: Worker) {
    this.worker = worker;
    this.worker.addEventListener('message', (msg) => this.handleMessage(msg.data));
  }

  setStore(store) {
    this.store = store;
  }

  loadWASM(filename: string): void {
    this.postMessage<LoadWASM>({
      type: MessageType.LoadWASM,
      data: { filename },
    });
  }


  subscribe(address: string): void {
    this.postMessage<Subscribe>({
      type: MessageType.Subscribe,
      data: { address },
    });
  }

  postMessage<T>(msg: T): void {
    this.worker.postMessage(msg);
  }

  handleWorkerStateChange(msg: WorkerEvent): void {
    console.log('[handleWorkerStateChange] msg:', msg);
    this.runCallbacks(msg.data.state);
  }

  runCallbacks(state: Event) {
    if (typeof this.callbacks[state] === 'undefined') {
      return;
    }

    this.callbacks[state]!.forEach((cb) => cb(state));
  }

  runOneOffs(state: Event) {
    if (typeof this.oneoffs[state] === 'undefined') {
      return;
    }

    this.oneoffs[state]!.forEach((cb) => cb(state));
    this.oneoffs[state] = [];
  }

  handleMessage(msg: GenericMessage): void {
    switch (msg.type) {
      case MessageType.WorkerEvent:
        return this.handleWorkerStateChange(msg as WorkerEvent);

      case MessageType.TxReceived:
        return this.handleTxReceived(msg as TxReceived);

      default:
        console.log('unknown msg type:', msg.type);
        console.log(msg);
    }
  }

  handleTxReceived(msg: TxReceived): void {
    const { tx } = msg.data;
    this.store.addTx(tx);
  }

  on(state: Event, cb: StateChangeCB, once: boolean = false) {
    let target = this.callbacks;
    if (once) {
      target = this.oneoffs;
    }

    if (typeof target[state] === 'undefined') {
      target[state] = [cb];
      return;
    }

    target[state]!.push(cb);
  }
}
