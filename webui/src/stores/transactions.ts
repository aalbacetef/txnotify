import { ref } from 'vue'
import { defineStore } from 'pinia'
import type { Transaction } from '@/lib/tx';


export const useTransactionStore = defineStore('transactions', () => {
  const transactions = ref<Transaction[]>([]);

  function addTx(tx: Transaction) {
    transactions.value.push(tx);
  }

  return { transactions, addTx };
})
