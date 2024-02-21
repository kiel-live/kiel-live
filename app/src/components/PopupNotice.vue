<template>
  <div
    v-if="show"
    class="flex fixed top-0 left-0 right-0 bottom-0 bg-gray-900 bg-opacity-80 z-1000 justify-center items-center"
    @click="$emit('close')"
  >
    <div
      class="m-2 flex flex-col rounded-md p-4 bg-white border-1 border-gray-200 shadow-xl z-20 md:w-104 dark:bg-dark-400 dark:text-gray-300 dark:border-dark-800"
    >
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  show: boolean;
}>();

defineEmits<{
  (event: 'close'): void;
}>();
</script>
