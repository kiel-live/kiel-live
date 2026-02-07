<template>
  <select
    v-model="model"
    class="bg-transparent p-2 w-full max-w-64 h-12 border border-gray-300 dark:border-gray-400 rounded-md"
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
