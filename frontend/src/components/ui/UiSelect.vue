<script setup lang="ts">
type Option = {
  value: string;
  label: string;
};

withDefaults(
  defineProps<{
    modelValue: string;
    options: Option[];
    ariaLabel?: string;
  }>(),
  {
    ariaLabel: "select",
  },
);

const emit = defineEmits<{
  "update:modelValue": [value: string];
}>();

function onChange(event: Event): void {
  emit("update:modelValue", (event.target as HTMLSelectElement).value);
}
</script>

<template>
  <select
    class="native-select rounded-xl border border-slate-300 bg-white px-3 py-2 text-sm text-slate-800 dark:border-slate-600 dark:bg-slate-800 dark:text-slate-100"
    :value="modelValue"
    :aria-label="ariaLabel"
    @change="onChange"
  >
    <option v-for="option in options" :key="option.value" :value="option.value">
      {{ option.label }}
    </option>
  </select>
</template>

