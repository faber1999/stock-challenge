<script setup lang="ts">
import { computed } from "vue";

type Tone = "success" | "error";

const props = withDefaults(
  defineProps<{
    tone?: Tone;
    class?: string;
  }>(),
  {
    tone: "success",
    class: "",
  },
);

const classes = computed(() => {
  const base = "mt-3 rounded-lg px-3 py-2 text-sm";
  const tone =
    props.tone === "error"
      ? "bg-red-50 text-red-700 dark:bg-red-900/40 dark:text-red-100"
      : "bg-emerald-50 text-emerald-800 dark:bg-emerald-900/40 dark:text-emerald-100";
  return `${base} ${tone} ${props.class}`.trim();
});
</script>

<template>
  <p :class="classes">
    <slot />
  </p>
</template>

