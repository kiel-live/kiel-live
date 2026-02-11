<template>
  <select
    v-model="model"
    class="h-10 w-full max-w-36 rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm transition-colors focus:border-gray-300 focus:outline-none dark:border-neutral-800 dark:bg-neutral-950 dark:focus:border-neutral-700"
  >
    <option v-for="option in options" :key="option.value" :value="option.value">{{ option.label }}</option>
  </select>
</template>

<script lang="ts" setup generic="T extends string | number">
defineProps<{
  options: { value: T; label: string }[];
}>();

const model = defineModel<T>({ required: true });
</script>
