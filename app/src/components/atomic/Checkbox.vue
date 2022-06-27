<template>
  <input
    type="checkbox"
    class="w-6 h-6 rounded"
    :checked="modelValue"
    @change="$emit('update:modelValue', !modelValue)"
  />
</template>

<srcript lang="ts" setup>
import { toRef } from 'vue';

const props = defineProps({
  modelValue: {
    type: Boolean,
    required: true,
  },
});

const modelValue = toRef('props');

defineEmits({
  'update:modelValue': (e: boolean) => true,
});
</srcript>
