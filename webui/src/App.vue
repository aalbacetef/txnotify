<script setup lang="ts">
import { inject, ref, computed } from 'vue';
import { WorkerPeer } from '@/lib/peer';
import { getEndpoints, type Endpoint } from '@/lib/tx';

import { useTransactionStore } from '@/stores/transactions';
import { useNotificationsStore } from '@/stores/notifications';


const store = useTransactionStore();
const ns = useNotificationsStore();

import AppNotifications from '@/components/app-notifications.vue';
import NotificationList from '@/components/notification-list.vue';
import ToggleMode from '@/components/toggle-mode.vue';


const rpcEndpoint = ref<string>('');
const ethAddress = ref<string>('0xdAC17F958D2ee523a2206206994597C13D831ec7');
const endpoints = ref<Endpoint[]>([]);
const workerPeer = inject<WorkerPeer>('workerPeer');
const customRPCURL = ref<string>('');
const started = ref<boolean>(false);
const customEndpointSet = ref<boolean>(false);

const rpcEndpointReady = computed(() => {
  if (rpcEndpoint.value === 'custom' && customEndpointSet.value) {
    return true;
  }

  if (rpcEndpoint.value === '' || rpcEndpoint.value === 'custom') {
    return false;
  }

  return true;
});

workerPeer!.setStore(store);


(async function () {
  endpoints.value = await getEndpoints();
})();


function handleEndpointSelected() {
  if (rpcEndpoint.value === "custom") {
    return;
  }

  workerPeer!.updateSettings({ rpcEndpoint: rpcEndpoint.value });
}

function handleSubscribeClicked() {
  if (started.value) {
    return;
  }

  const address = ethAddress.value;
  if (address.trim() === '') {
    return;
  }

  workerPeer!.subscribe(address);
  started.value = true;
  workerPeer!.start();
}

function setCustomRPCEndpoint() {
  if (customRPCURL.value.trim() === '') {
    return;
  }

  customEndpointSet.value = true;
  workerPeer!.updateSettings({ rpcEndpoint: customRPCURL.value });
  ns.pushNotification("custom endpoint has been set!");
}
</script>

<template>
  <main>
    <div class="container">
      <div class="box">
        <div class="header-row">
          <h1 class="title is-4">txnotify</h1>
          <toggle-mode></toggle-mode>
        </div>

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
              <select v-model="rpcEndpoint" @change="handleEndpointSelected">
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
            <input class="input" type="text" placeholder="Enter custom RPC URL" v-model="customRPCURL" />
          </div>
        </div>

        <div class="buttons">
          <button class="button subscribe" :disabled="started" @click="handleSubscribeClicked">Subscribe</button>
          <button class="button" v-if="rpcEndpoint === 'custom'" @click="setCustomRPCEndpoint">
            Set custom RPC endpoint
          </button>

        </div>

        <notification-list :started="started" :eth-address="ethAddress"></notification-list>
      </div>
    </div>

    <app-notifications></app-notifications>
  </main>
</template>


<style scoped>
main {
  position: relative;
}

.header-row {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
}

.container {
  margin: 2rem auto;
  padding: 0 1rem;
}

.box {
  padding: 2rem;
}

.field {
  margin-bottom: 1.5rem;
}

.subtitle {
  margin-bottom: 1rem !important;
}

.button.subscribe {
  margin-right: 5px;
}
</style>
