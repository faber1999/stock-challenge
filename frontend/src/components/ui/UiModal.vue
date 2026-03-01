<script setup lang="ts">
import UiButton from "./UiButton.vue";

defineProps<{
  open: boolean;
  closeLabel: string;
}>();

const emit = defineEmits<{
  close: [];
}>();
</script>

<template>
  <section
    v-if="open"
    class="fixed inset-0 z-10 bg-slate-900/40 p-4 backdrop-blur-sm"
    @click.self="emit('close')"
  >
    <article class="mx-auto mt-16 max-w-xl rounded-3xl border border-slate-200 bg-white p-6 shadow-2xl dark:border-slate-700 dark:bg-slate-900">
      <div class="mb-4 flex items-start justify-between gap-3">
        <div>
          <slot name="header" />
        </div>
        <UiButton variant="outline" size="sm" @click="emit('close')">
          {{ closeLabel }}
        </UiButton>
      </div>
      <slot />
    </article>
  </section>
</template>

