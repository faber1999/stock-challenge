<script setup lang="ts">
import { computed } from "vue";

type Variant = "primary" | "secondary" | "outline";
type Size = "sm" | "md";

const props = withDefaults(
  defineProps<{
    variant?: Variant;
    size?: Size;
    disabled?: boolean;
    type?: "button" | "submit" | "reset";
    class?: string;
  }>(),
  {
    variant: "primary",
    size: "md",
    disabled: false,
    type: "button",
    class: "",
  },
);

const classes = computed(() => {
  const base =
    "rounded-xl font-semibold transition duration-200 ease-in-out disabled:cursor-not-allowed disabled:opacity-60";
  const size = props.size === "sm" ? "px-3 py-1.5 text-sm" : "px-4 py-2 text-sm";

  const variant =
    props.variant === "secondary"
      ? "border border-slate-300 bg-transparent text-slate-700 hover:bg-slate-100 dark:border-slate-600 dark:bg-transparent dark:text-slate-100 dark:hover:bg-slate-700"
      : props.variant === "outline"
        ? "border border-slate-300 text-slate-700 hover:bg-slate-100 dark:border-slate-600 dark:text-slate-200 dark:hover:bg-slate-800"
        : "bg-indigo-600 text-white hover:bg-indigo-500";

  return `${base} ${size} ${variant} ${props.class}`.trim();
});
</script>

<template>
  <button :type="type" :disabled="disabled" :class="classes">
    <slot />
  </button>
</template>
