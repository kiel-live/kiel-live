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

defineEmits({
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  'update:modelValue': (_e: boolean) => true,
});
</script>
