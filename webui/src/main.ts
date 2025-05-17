import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import { WorkerPeer } from '@/lib/peer'
import { Event } from '@/lib/messages'

const app = createApp(App)


const worker = new Worker(new URL('@/lib/worker/index.ts', import.meta.url), { type: 'module' });
const workerPeer = new WorkerPeer(worker);

workerPeer.on(Event.WASMLoaded, () => {
  app.use(createPinia());

  app.provide<WorkerPeer>('workerPeer', workerPeer);

  app.mount('#app');
});

workerPeer.loadWASM(new URL('txnotify.wasm', document.baseURI).toString());


