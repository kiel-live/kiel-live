<template>
  <select
    v-model="model"
    class="h-12 w-full max-w-64 rounded-md border border-gray-300 bg-transparent p-2 dark:border-gray-400"
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
