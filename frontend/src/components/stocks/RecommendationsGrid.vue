<script setup lang="ts">
import type { RecommendationsResponse } from "../../types/stocks";

defineProps<{
  items: RecommendationsResponse["items"];
  title: string;
  scoreLabel: string;
  upsideLabel: string;
}>();
</script>

<template>
  <section class="mb-6 rounded-3xl border border-slate-200/70 bg-white/85 p-6 shadow-md shadow-slate-200/40 dark:border-slate-700 dark:bg-slate-900/90 dark:shadow-slate-950/60">
    <h2 class="mb-4 text-xl font-bold text-slate-900 dark:text-slate-100">{{ title }}</h2>
    <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-5">
      <article
        v-for="item in items"
        :key="item.id"
        class="rounded-2xl border border-emerald-100 bg-gradient-to-br from-emerald-50 to-cyan-50 p-4 dark:border-emerald-800 dark:from-emerald-950 dark:to-cyan-950"
      >
        <p class="text-xs font-semibold uppercase tracking-wider text-emerald-800 dark:text-emerald-300">{{ item.ticker }}</p>
        <h3 class="mt-1 truncate text-sm font-bold text-slate-900 dark:text-slate-100">{{ item.company }}</h3>
        <p class="mt-2 text-xs text-slate-600 dark:text-slate-300">{{ item.action }}</p>
        <div class="mt-3 flex items-center justify-between">
          <span class="text-xs text-slate-500 dark:text-slate-400">{{ scoreLabel }}</span>
          <strong class="text-sm text-emerald-800 dark:text-emerald-300">{{ item.score.toFixed(2) }}</strong>
        </div>
        <div class="mt-1 flex items-center justify-between">
          <span class="text-xs text-slate-500 dark:text-slate-400">{{ upsideLabel }}</span>
          <strong class="text-sm text-cyan-800 dark:text-cyan-300">{{ item.upside_pct.toFixed(2) }}%</strong>
        </div>
      </article>
    </div>
  </section>
</template>

