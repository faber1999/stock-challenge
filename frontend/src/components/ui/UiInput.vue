<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    modelValue: string;
    placeholder?: string;
    ariaLabel?: string;
    class?: string;
  }>(),
  {
    placeholder: "",
    ariaLabel: "input",
    class: "",
  },
);

const emit = defineEmits<{
  "update:modelValue": [value: string];
}>();

function onInput(event: Event): void {
  emit("update:modelValue", (event.target as HTMLInputElement).value);
}
</script>

<template>
  <input
    type="text"
    :value="props.modelValue"
    :placeholder="placeholder"
    :aria-label="ariaLabel"
    class="rounded-xl border border-slate-300 bg-white px-3 py-2 text-sm text-slate-800 dark:border-slate-600 dark:bg-slate-800 dark:text-slate-100"
    :class="props.class"
    @input="onInput"
  />
</template>

