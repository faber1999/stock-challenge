<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import { useStocksStore } from "./stores/stocks";
import type { AppLocale } from "./i18n";
import type { SortBy, SortOrder } from "./types/stocks";
import UiAlert from "./components/ui/UiAlert.vue";
import UiButton from "./components/ui/UiButton.vue";
import UiSelect from "./components/ui/UiSelect.vue";
import RecommendationsGrid from "./components/stocks/RecommendationsGrid.vue";
import StocksFilters from "./components/stocks/StocksFilters.vue";
import StocksTable from "./components/stocks/StocksTable.vue";
import UiPagination from "./components/ui/UiPagination.vue";
import StockDetailsModal from "./components/stocks/StockDetailsModal.vue";

type Theme = "light" | "dark";

const store = useStocksStore();
const { t, locale } = useI18n();

const theme = ref<Theme>("light");
const localeOptions = [
  { value: "es", label: "Español" },
  { value: "en", label: "English" },
];
const sortOptions: Array<{ value: SortBy; label: string }> = [
  { value: "recommend_score", label: "recommend_score" },
  { value: "ticker", label: "ticker" },
  { value: "company", label: "company" },
  { value: "brokerage", label: "brokerage" },
  { value: "target_to", label: "target_to" },
  { value: "synced_at", label: "synced_at" },
];
const orderOptions: Array<{ value: SortOrder; label: string }> = [
  { value: "desc", label: "desc" },
  { value: "asc", label: "asc" },
];

const themeLabel = computed(() =>
  theme.value === "dark" ? t("app.themeDark") : t("app.themeLight"),
);

onMounted(async () => {
  const savedTheme = localStorage.getItem("theme");
  const savedLocale = localStorage.getItem("locale");

  if (savedTheme === "dark" || savedTheme === "light") {
    theme.value = savedTheme;
  } else if (window.matchMedia("(prefers-color-scheme: dark)").matches) {
    theme.value = "dark";
  }

  if (savedLocale === "en" || savedLocale === "es") {
    locale.value = savedLocale;
  }

  applyTheme(theme.value);
  await Promise.all([store.loadStocks(), store.loadRecommendations()]);
});

watch(theme, (value) => {
  applyTheme(value);
  localStorage.setItem("theme", value);
});

watch(locale, (value) => {
  localStorage.setItem("locale", value);
});

function applyTheme(value: Theme): void {
  document.documentElement.classList.toggle("dark", value === "dark");
}

function toggleTheme(): void {
  theme.value = theme.value === "dark" ? "light" : "dark";
}

function changeLocale(value: string): void {
  locale.value = value as AppLocale;
}

const currentLocale = computed(() => (locale.value === "es" ? "es-CO" : "en-US"));
</script>

<template>
  <main class="mx-auto max-w-7xl p-4 md:p-8">
    <section class="mb-6 rounded-3xl border border-slate-200/70 bg-white/85 p-6 shadow-lg shadow-slate-200/60 dark:border-slate-700 dark:bg-slate-900/90 dark:shadow-slate-950/60">
      <div class="mb-4 flex flex-wrap items-center justify-end gap-2">
        <label class="text-xs font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-300">
          {{ t("app.language") }}
        </label>
        <UiSelect :model-value="locale" :options="localeOptions" aria-label="language" @update:model-value="changeLocale" />
        <UiButton variant="secondary" @click="toggleTheme()">
          {{ themeLabel }}
        </UiButton>
      </div>

      <div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <p class="text-sm font-semibold uppercase tracking-[0.14em] text-slate-500 dark:text-slate-400">{{ t("app.badge") }}</p>
          <h1 class="mt-1 text-3xl font-extrabold text-slate-900 dark:text-slate-100">{{ t("app.title") }}</h1>
          <p class="mt-2 text-sm text-slate-600 dark:text-slate-300">{{ t("app.subtitle") }}</p>
        </div>
        <div class="flex flex-wrap gap-3">
          <UiButton :disabled="store.syncing" @click="store.runSync(10)">
            {{ store.syncing ? t("app.syncing") : t("app.sync") }}
          </UiButton>
          <UiButton variant="secondary" @click="store.loadRecommendations(5)">
            {{ t("app.refreshTop") }}
          </UiButton>
        </div>
      </div>
      <UiAlert v-if="store.syncInfo">
        {{ t("app.syncComplete", { stocks: store.syncInfo.stocks, pages: store.syncInfo.pages }) }}
      </UiAlert>
      <UiAlert v-if="store.error" tone="error">
        {{ store.error }}
      </UiAlert>
    </section>

    <RecommendationsGrid
      :items="store.recommendations"
      :title="t('app.topRecommendations')"
      :score-label="t('app.score')"
      :upside-label="t('app.upside')"
    />

    <section class="rounded-3xl border border-slate-200/70 bg-white/90 p-6 shadow-md shadow-slate-200/40 dark:border-slate-700 dark:bg-slate-900/90 dark:shadow-slate-950/60">
      <StocksFilters
        :q="store.filters.q"
        :action="store.filters.action"
        :sort-by="store.filters.sortBy"
        :order="store.filters.order"
        :sort-options="sortOptions"
        :order-options="orderOptions"
        :search-placeholder="t('app.searchPlaceholder')"
        :action-placeholder="t('app.actionPlaceholder')"
        :apply-label="t('app.apply')"
        @update:q="store.filters.q = $event"
        @update:action="store.filters.action = $event"
        @update:sort-by="store.filters.sortBy = $event"
        @update:order="store.filters.order = $event"
        @apply="store.applyFilters()"
      />

      <StocksTable
        :items="store.items"
        :loading="store.loading"
        :empty-label="t('app.emptyStocks')"
        :headers="{
          ticker: t('app.ticker'),
          company: t('app.company'),
          action: t('app.action'),
          target: t('app.target'),
          score: t('app.score'),
          updated: t('app.updated'),
        }"
        :locale="currentLocale"
        @select="store.openDetails($event)"
      />

      <UiPagination
        :page="store.page"
        :total-pages="store.totalPages"
        :total-items="store.total"
        :page-label="t('app.page')"
        :of-label="t('app.of')"
        :total-label="t('app.total')"
        :items-label="t('app.items')"
        :prev-label="t('app.prev')"
        :next-label="t('app.next')"
        @prev="store.previousPage()"
        @next="store.nextPage()"
      />
    </section>
  </main>

  <StockDetailsModal
    :stock="store.selected"
    :close-label="t('app.close')"
    :labels="{
      action: t('app.detailsAction'),
      rating: t('app.detailsRating'),
      target: t('app.detailsTarget'),
      recommendScore: t('app.detailsRecommendScore'),
    }"
    :locale="currentLocale"
    @close="store.closeDetails()"
  />
</template>
