<template>
  <input
    type="checkbox"
    class="h-5 w-5 shrink-0 rounded"
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
