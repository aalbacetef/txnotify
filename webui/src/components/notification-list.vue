<script setup lang="ts">
import { computed } from 'vue';
import { useTransactionStore } from '@/stores/transactions';

const props = defineProps<{ ethAddress: string; started: boolean; }>();

const store = useTransactionStore();
const transactions = computed(() => store.transactions.toReversed());

const loading = computed(() => props.started && transactions.value.length === 0);

function etherURL(txHash: string): string {
  return `https://etherscan.io/tx/${txHash}`;
}
</script>

<template>
  <div class="transaction-list">
    <h2 class="subtitle is-5">
      Transactions
      <button v-if="loading" class="button is-loading"></button>
      <span v-if="transactions.length > 0">
        ({{ transactions.length }})
      </span>
    </h2>
    <div class="notification is-info" v-if="started">
      <p>Listening for transactions...</p>
    </div>
    <div class="notification is-primary" v-else-if="!started">
      <p>Please enter a valid Ethereum address, select and endpoint and hit subscribe.</p>
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
        <p class="view-on-etherscan">
          <a class="button" target="_blank" :href="etherURL(tx.hash)">
            View on Etherscan
          </a>
        </p>
      </div>
    </article>
  </div>
</template>

<style scoped>
.transaction-list {
  margin-top: 2rem;
}

.message {
  margin-bottom: 1rem;
}

.message-body {
  word-break: break-all;
}

.button.is-loading {
  width: 30px;
  height: 30px;
  border: 0;
}

.view-on-etherscan {
  text-align: right;
}
</style>
