<template>
  <input
    type="checkbox"
    class="h-6 w-6 shrink-0 cursor-pointer rounded-md accent-red-600 transition-all dark:accent-red-500"
    :checked="modelValue"
    @change="$emit('update:modelValue', !modelValue)"
  />
</template>

<script lang="ts" setup>
import { toRef } from 'vue';

const props = defineProps({
  modelValue: {
    type: Boolean,
    required: true,
  },
});

defineEmits<{
  (e: 'update:modelValue', v: boolean): void;
}>();

const modelValue = toRef(props, 'modelValue');
</script>
