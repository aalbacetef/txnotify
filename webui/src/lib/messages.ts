import type { Transaction } from "./tx";

export type GenericMessage = {
  type: MessageType;
  data: any;
};

export enum Event {
  WorkerLoaded = 'worker-loaded',
  WASMLoaded = 'wasm-loaded',
  TxReceived = 'tx-received',
}

export enum MessageType {
  WorkerEvent = 'worker-event',
  LoadWASM = 'load-wasm',
  Subscribe = "subscribe",
  TxReceived = 'tx-received',
}

export type WorkerEvent = {
  type: MessageType.WorkerEvent;
  data: {
    state: Event;
  };
};

export type LoadWASM = {
  type: MessageType.LoadWASM;
  data: {
    filename: string;
  };
};

export type Subscribe = {
  type: MessageType.Subscribe;
  data: {
    address: string;
  };
}

export type TxReceived = {
  type: MessageType.TxReceived;
  data: {
    tx: Transaction;
  }
}
