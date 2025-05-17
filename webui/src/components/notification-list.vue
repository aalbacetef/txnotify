<script setup lang="ts">
import { computed } from 'vue';
import { useTransactionStore } from '@/stores/transactions';

defineProps<{ ethAddress: string }>();

const store = useTransactionStore();
const transactions = computed(() => store.transactions.toReversed());
</script>



<template>
  <div class="transaction-list">
    <h2 class="subtitle is-5">
      Transactions
      <span v-if="transactions.length > 0">
        ({{ transactions.length }})
      </span>
    </h2>
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
</style>
