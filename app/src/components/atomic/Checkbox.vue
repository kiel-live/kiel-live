<template>
  <input
    type="checkbox"
    class="w-5 h-5 flex-shrink-0 rounded"
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

const modelValue = toRef(props, 'modelValue');

defineEmits<{
  (e: 'update:modelValue', v: boolean): void;
}>();
</script>
