<script setup lang="ts">
import { inject, ref, computed } from 'vue';
import { WorkerPeer } from '@/lib/peer';

import { useTransactionStore } from '@/stores/transactions';

const store = useTransactionStore();

import ToggleMode from '@/components/toggle-mode.vue';

type Endpoint = {
  url: string;
}

const rpcEndpoint = ref<string>('');
const ethAddress = ref<string>('0xdAC17F958D2ee523a2206206994597C13D831ec7');
const endpoints = ref<Endpoint[]>([]);
const workerPeer = inject<WorkerPeer>('workerPeer');
workerPeer.setStore(store);

const transactions = computed(() => store.transactions.toReversed());

// abbreviated version
type ChainlistDataEntry = {
  name: string;
  chain: string;
  rpc: Endpoint[];
}

(async function () {
  const response = await fetch(
    'https://chainlist.org/rpcs.json',
    { method: 'GET', mode: 'cors' },
  );
  const data = await response.json<ChainlistDataEntry>();

  const mainnet = data.find(
    elem => elem.name === "Ethereum Mainnet" && elem.chain === "ETH"
  );
  if (mainnet === null) {
    console.error("could not find RPC endpoints");
    return;
  }

  endpoints.value = mainnet.rpc.filter(elem => elem.url.startsWith("https://"));
})();


function handleSubscribeClicked() {
  const address = ethAddress.value;
  console.log('address: ', address);
  if (address.trim() === '') {
    return;
  }

  workerPeer.subscribe(address);
}
</script>

<template>
  <header>
    <toggle-mode></toggle-mode>
  </header>

  <main>
    <div class="container">
      <div class="box">
        <h1 class="title is-4">txnotify</h1>

        <div class="field">
          <label class="label">Ethereum Address</label>
          <div class="control">
            <input class="input" type="text" placeholder="Enter Ethereum address (e.g., 0x...)" v-model="ethAddress" />
          </div>
        </div>

        <div class="field">
          <label class="label">RPC Endpoint</label>
          <div class="control">
            <div class="select is-fullwidth">
              <select v-model="rpcEndpoint">
                <option value="" disabled>Select an RPC endpoint</option>
                <option value="custom">Custom</option>
                <option v-for="endpoint in endpoints" :key="endpoint.url">
                  {{ endpoint.url }}
                </option>
              </select>
            </div>
          </div>
        </div>

        <div class="field" v-if="rpcEndpoint === 'custom'">
          <label class="label">Custom RPC URL</label>
          <div class="control">
            <input class="input" type="text" placeholder="Enter custom RPC URL" v-model="customRpcUrl" />
          </div>
        </div>

        <div class="button">
          <button @click="handleSubscribeClicked">Subscribe</button>
        </div>

        <div class="transaction-list">
          <h2 class="subtitle is-5">Transactions</h2>
          <div class="notification is-info" v-if="!transactions.length && ethAddress">
            <p>Listening for transactions...</p>
          </div>
          <div class="notification is-warning" v-else-if="!ethAddress">
            <p>Please enter a valid Ethereum address.</p>
          </div>
          <article class="message is-dark" v-for="tx in transactions" :key="tx.hash">
            <div class="message-header">
              <p>Transaction</p>
            </div>
            <div class="message-body">
              <p><strong>Hash:</strong> {{ tx.hash }}</p>
              <p><strong>From:</strong> {{ tx.from }}</p>
              <p><strong>To:</strong> {{ tx.to }}</p>
              <p><strong>Value:</strong> {{ tx.value }} ETH</p>
              <p><strong>Block:</strong> {{ tx.blockNumber }}</p>
            </div>
          </article>
        </div>
      </div>
    </div>
  </main>
</template>


<style scoped>
.container {
  max-width: 600px;
  margin: 2rem auto;
  padding: 0 1rem;
}

.box {
  padding: 2rem;
}

.field {
  margin-bottom: 1.5rem;
}

.transaction-list {
  margin-top: 2rem;
}

.message {
  margin-bottom: 1rem;
}

.message-body {
  word-break: break-all;
}

.subtitle {
  margin-bottom: 1rem !important;
}
</style>
