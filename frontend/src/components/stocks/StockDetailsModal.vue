<script setup lang="ts">
import UiModal from "../ui/UiModal.vue";
import type { Stock } from "../../types/stocks";

const props = defineProps<{
  stock: Stock | null;
  closeLabel: string;
  labels: {
    action: string;
    rating: string;
    target: string;
    recommendScore: string;
  };
  locale: string;
}>();

const emit = defineEmits<{
  close: [];
}>();

function formatCurrency(value: number, currency = "USD"): string {
  return new Intl.NumberFormat(props.locale, {
    style: "currency",
    currency,
    maximumFractionDigits: 2,
  }).format(value);
}
</script>

<template>
  <UiModal :open="Boolean(stock)" :close-label="closeLabel" @close="emit('close')">
    <template #header>
      <template v-if="stock">
        <p class="text-xs uppercase tracking-widest text-indigo-600 dark:text-indigo-300">{{ stock.ticker }}</p>
        <h3 class="text-2xl font-bold text-slate-900 dark:text-slate-100">{{ stock.company }}</h3>
        <p class="mt-1 text-sm text-slate-600 dark:text-slate-300">{{ stock.brokerage }}</p>
      </template>
    </template>

    <dl v-if="stock" class="grid grid-cols-2 gap-3 text-sm">
      <div class="rounded-xl bg-slate-50 p-3 dark:bg-slate-800">
        <dt class="text-slate-500 dark:text-slate-300">{{ labels.action }}</dt>
        <dd class="font-semibold text-slate-900 dark:text-slate-100">{{ stock.action }}</dd>
      </div>
      <div class="rounded-xl bg-slate-50 p-3 dark:bg-slate-800">
        <dt class="text-slate-500 dark:text-slate-300">{{ labels.rating }}</dt>
        <dd class="font-semibold text-slate-900 dark:text-slate-100">{{ stock.rating_from }} → {{ stock.rating_to }}</dd>
      </div>
      <div class="rounded-xl bg-slate-50 p-3 dark:bg-slate-800">
        <dt class="text-slate-500 dark:text-slate-300">{{ labels.target }}</dt>
        <dd class="font-semibold text-slate-900 dark:text-slate-100">
          {{ formatCurrency(stock.target_from, stock.currency) }} →
          {{ formatCurrency(stock.target_to, stock.currency) }}
        </dd>
      </div>
      <div class="rounded-xl bg-slate-50 p-3 dark:bg-slate-800">
        <dt class="text-slate-500 dark:text-slate-300">{{ labels.recommendScore }}</dt>
        <dd class="font-semibold text-slate-900 dark:text-slate-100">{{ stock.recommend_score.toFixed(2) }}</dd>
      </div>
    </dl>
  </UiModal>
</template>

