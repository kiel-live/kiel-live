<template>
  <select
    v-model="innerValue"
    class="bg-transparent p-2 w-full max-w-64 h-12 border border-transparent rounded-md border-gray-300 border-opacity-50 focus-visible:(outline-none border-blue-700 border-opacity-100 dark:border-blue-400)"
  >
    <option v-for="option in options" :key="option.value" :value="option.value">{{ option.label }}</option>
  </select>
</template>

<script lang="ts" setup>
import { computed, toRef } from 'vue';

type OptionValue = string | number;

const props = defineProps<{
  modelValue: OptionValue;
  options: { value: OptionValue; label: string }[];
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', v: OptionValue): void;
}>();

const modelValue = toRef(props, 'modelValue');
const innerValue = computed({
  get: () => modelValue.value,
  set: (value) => {
    emit('update:modelValue', value);
  },
});
</script>
