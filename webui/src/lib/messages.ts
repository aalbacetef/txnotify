import type { Transaction } from "./tx";
import type { Settings } from "./wasm";

export type GenericMessage = {
  type: MessageType;
  data: any;
};

export enum Event {
  WorkerLoaded = 'worker-loaded',
  WASMLoaded = 'wasm-loaded',
}

export enum MessageType {
  WorkerEvent = 'worker-event',
  LoadWASM = 'load-wasm',
  Subscribe = "subscribe",
  TxNotification = 'tx-received',
  UpdateSettings = 'update-settings',
  StartWatcher = 'start-watcher',
};

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

export type TxNotification = {
  type: MessageType.TxNotification;
  data: {
    tx: Transaction;
  }
}

export type UpdateSettings = {
  type: MessageType.UpdateSettings;
  data: Settings;
}

export type StartWatcher = {
  type: MessageType.StartWatcher;
  data: null;
}
