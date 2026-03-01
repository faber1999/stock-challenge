<script setup lang="ts">
import type { Stock } from "../../types/stocks";

const props = defineProps<{
  items: Stock[];
  loading: boolean;
  emptyLabel: string;
  headers: {
    ticker: string;
    company: string;
    action: string;
    target: string;
    score: string;
    updated: string;
  };
  locale: string;
}>();

const emit = defineEmits<{
  select: [ticker: string];
}>();

function formatCurrency(value: number, currency = "USD"): string {
  return new Intl.NumberFormat(props.locale, {
    style: "currency",
    currency,
    maximumFractionDigits: 2,
  }).format(value);
}

function formatDate(value: string): string {
  return new Date(value).toLocaleString(props.locale);
}
</script>

<template>
  <div class="overflow-x-auto rounded-2xl border border-slate-200 dark:border-slate-700">
    <table class="min-w-full divide-y divide-slate-200 text-left text-sm dark:divide-slate-700">
      <thead class="bg-slate-50 dark:bg-slate-800">
        <tr>
          <th class="px-4 py-3 font-semibold text-slate-700 dark:text-slate-100">{{ headers.ticker }}</th>
          <th class="px-4 py-3 font-semibold text-slate-700 dark:text-slate-100">{{ headers.company }}</th>
          <th class="px-4 py-3 font-semibold text-slate-700 dark:text-slate-100">{{ headers.action }}</th>
          <th class="px-4 py-3 font-semibold text-slate-700 dark:text-slate-100">{{ headers.target }}</th>
          <th class="px-4 py-3 font-semibold text-slate-700 dark:text-slate-100">{{ headers.score }}</th>
          <th class="px-4 py-3 font-semibold text-slate-700 dark:text-slate-100">{{ headers.updated }}</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-slate-100 bg-white dark:divide-slate-800 dark:bg-slate-900">
        <tr
          v-for="stock in items"
          :key="stock.id"
          class="cursor-pointer transition hover:bg-indigo-50 dark:hover:bg-slate-800"
          @click="emit('select', stock.ticker)"
        >
          <td class="px-4 py-3 font-semibold text-slate-900 dark:text-slate-100">{{ stock.ticker }}</td>
          <td class="px-4 py-3 text-slate-700 dark:text-slate-200">{{ stock.company }}</td>
          <td class="px-4 py-3 text-slate-600 dark:text-slate-300">{{ stock.action }}</td>
          <td class="px-4 py-3 text-slate-700 dark:text-slate-200">{{ formatCurrency(stock.target_to, stock.currency) }}</td>
          <td class="px-4 py-3 text-slate-700 dark:text-slate-200">{{ stock.recommend_score.toFixed(2) }}</td>
          <td class="px-4 py-3 text-slate-500 dark:text-slate-400">{{ formatDate(stock.synced_at) }}</td>
        </tr>
        <tr v-if="!loading && items.length === 0">
          <td colspan="6" class="px-4 py-6 text-center text-slate-500 dark:text-slate-400">{{ emptyLabel }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

