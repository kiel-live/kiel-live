<template>
  <select
    v-model="model"
    class="bg-transparent p-2 w-full max-w-64 h-12 border border-transparent rounded-md border-gray-300 border-opacity-50 focus-visible:(outline-none border-blue-700 border-opacity-100 dark:border-blue-400)"
  >
    <option v-for="option in options" :key="option.value" :value="option.value">{{ option.label }}</option>
  </select>
</template>

<script lang="ts" setup generic="T extends string | number">

defineProps<{
  options: { value: T; label: string }[];
}>();

const model = defineModel<T>({required: true});
</script>
