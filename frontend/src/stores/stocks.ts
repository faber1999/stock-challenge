import { defineStore } from "pinia";
import { computed, ref } from "vue";
import { getStock, listRecommendations, listStocks, syncStocks } from "../lib/api";
import type { ListStocksParams, RecommendationsResponse, Stock } from "../types/stocks";

export const useStocksStore = defineStore("stocks", () => {
  const loading = ref(false);
  const syncing = ref(false);
  const error = ref("");
  const syncInfo = ref<{ stocks: number; pages: number } | null>(null);
  const items = ref<Stock[]>([]);
  const total = ref(0);
  const recommendations = ref<RecommendationsResponse["items"]>([]);
  const selected = ref<Stock | null>(null);

  const filters = ref<ListStocksParams>({
    q: "",
    action: "",
    sortBy: "recommend_score",
    order: "desc",
    limit: 20,
    offset: 0,
  });

  const page = computed(() => Math.floor(filters.value.offset / filters.value.limit) + 1);
  const totalPages = computed(() => Math.max(1, Math.ceil(total.value / filters.value.limit)));

  async function loadStocks() {
    loading.value = true;
    error.value = "";
    try {
      const response = await listStocks(filters.value);
      items.value = response.items;
      total.value = response.total;
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to load stocks";
    } finally {
      loading.value = false;
    }
  }

  async function loadRecommendations(limit = 5) {
    try {
      const response = await listRecommendations(limit);
      recommendations.value = response.items;
    } catch {
      recommendations.value = [];
    }
  }

  async function runSync(limit = 10) {
    syncing.value = true;
    error.value = "";
    syncInfo.value = null;
    try {
      const response = await syncStocks(limit);
      syncInfo.value = { stocks: response.stocks_saved, pages: response.pages_processed };
      filters.value.offset = 0;
      await Promise.all([loadStocks(), loadRecommendations()]);
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to sync stocks";
    } finally {
      syncing.value = false;
    }
  }

  async function openDetails(ticker: string) {
    try {
      selected.value = await getStock(ticker);
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to load stock details";
    }
  }

  function closeDetails() {
    selected.value = null;
  }

  function nextPage() {
    const nextOffset = filters.value.offset + filters.value.limit;
    if (nextOffset < total.value) {
      filters.value.offset = nextOffset;
      void loadStocks();
    }
  }

  function previousPage() {
    filters.value.offset = Math.max(0, filters.value.offset - filters.value.limit);
    void loadStocks();
  }

  function applyFilters() {
    filters.value.offset = 0;
    void loadStocks();
  }

  return {
    loading,
    syncing,
    error,
    syncInfo,
    items,
    total,
    recommendations,
    selected,
    filters,
    page,
    totalPages,
    loadStocks,
    loadRecommendations,
    runSync,
    openDetails,
    closeDetails,
    nextPage,
    previousPage,
    applyFilters,
  };
});
