import {
  MessageType,
  Event,
  type Subscribe,
  type TxNotification,
  type UpdateSettings,
  type StartWatcher,
} from '@/lib/messages';

import {
  type GenericMessage,
  type LoadWASM,
  type WorkerEvent,
} from '@/lib/messages';
import type { Transaction } from './tx';


type StateChangeCB = (state: Event) => void;
type Settings = {
  rpcEndpoint: string;
}

type Store = {
  addTx(tx: Transaction): void;
}

// WorkerPeer provides a set of methods to interact with the Worker.
export class WorkerPeer {
  store: Store | null = null;
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

  setStore(store: Store) {
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

  start(): void {
    this.postMessage<StartWatcher>({
      type: MessageType.StartWatcher,
      data: null,
    });
  }

  updateSettings(settings: Settings): void {
    this.postMessage<UpdateSettings>({
      type: MessageType.UpdateSettings,
      data: settings,
    });
  }


  handleWorkerStateChange(msg: WorkerEvent): void {
    console.log('[handleWorkerStateChange] msg:', msg);
    this.runCallbacks(msg.data.state);
  }

  handleMessage(msg: GenericMessage): void {
    switch (msg.type) {
      case MessageType.WorkerEvent:
        return this.handleWorkerStateChange(msg as WorkerEvent);

      case MessageType.TxNotification:
        return this.handleTxReceived(msg as TxNotification);

      default:
        console.log('unknown msg type:', msg.type);
        console.log(msg);
    }
  }

  handleTxReceived(msg: TxNotification): void {
    const { tx } = msg.data;
    if (this.store === null) {
      return;
    }

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


  postMessage<T>(msg: T): void {
    this.worker.postMessage(msg);
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
}
