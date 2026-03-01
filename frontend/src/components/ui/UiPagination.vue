<script setup lang="ts">
import UiButton from "./UiButton.vue";

defineProps<{
  page: number;
  totalPages: number;
  totalItems: number;
  pageLabel: string;
  ofLabel: string;
  totalLabel: string;
  itemsLabel: string;
  prevLabel: string;
  nextLabel: string;
}>();

const emit = defineEmits<{
  prev: [];
  next: [];
}>();
</script>

<template>
  <div class="mt-4 flex items-center justify-between">
    <p class="text-sm text-slate-600 dark:text-slate-300">
      {{ pageLabel }} {{ page }} {{ ofLabel }} {{ totalPages }} · {{ totalLabel }} {{ totalItems }} {{ itemsLabel }}
    </p>
    <div class="flex gap-2">
      <UiButton variant="outline" size="sm" :disabled="page <= 1" @click="emit('prev')">
        {{ prevLabel }}
      </UiButton>
      <UiButton variant="outline" size="sm" :disabled="page >= totalPages" @click="emit('next')">
        {{ nextLabel }}
      </UiButton>
    </div>
  </div>
</template>

