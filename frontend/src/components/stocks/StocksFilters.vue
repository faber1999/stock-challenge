<script setup lang="ts">
import UiButton from "../ui/UiButton.vue";
import UiInput from "../ui/UiInput.vue";
import UiSelect from "../ui/UiSelect.vue";
import type { SortBy, SortOrder } from "../../types/stocks";

defineProps<{
  q: string;
  action: string;
  sortBy: SortBy;
  order: SortOrder;
  sortOptions: Array<{ value: SortBy; label: string }>;
  orderOptions: Array<{ value: SortOrder; label: string }>;
  searchPlaceholder: string;
  actionPlaceholder: string;
  applyLabel: string;
}>();

const emit = defineEmits<{
  "update:q": [value: string];
  "update:action": [value: string];
  "update:sortBy": [value: SortBy];
  "update:order": [value: SortOrder];
  apply: [];
}>();

function onSortBy(value: string): void {
  emit("update:sortBy", value as SortBy);
}

function onOrder(value: string): void {
  emit("update:order", value as SortOrder);
}
</script>

<template>
  <div class="mb-4 grid gap-3 md:grid-cols-6">
    <UiInput
      :model-value="q"
      :placeholder="searchPlaceholder"
      class="md:col-span-2"
      aria-label="search"
      @update:model-value="emit('update:q', $event)"
    />
    <UiInput
      :model-value="action"
      :placeholder="actionPlaceholder"
      aria-label="action filter"
      @update:model-value="emit('update:action', $event)"
    />
    <UiSelect
      :model-value="sortBy"
      :options="sortOptions"
      aria-label="sort by"
      @update:model-value="onSortBy"
    />
    <UiSelect
      :model-value="order"
      :options="orderOptions"
      aria-label="sort order"
      @update:model-value="onOrder"
    />
    <UiButton @click="emit('apply')">
      {{ applyLabel }}
    </UiButton>
  </div>
</template>
