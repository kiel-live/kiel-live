<template>
  <div
    v-if="show"
    class="bg-opacity-80 fixed top-0 right-0 bottom-0 left-0 z-1000 flex items-center justify-center bg-neutral-800"
    @click="$emit('close')"
  >
    <div
      class="z-20 m-2 flex flex-col rounded-md border border-gray-200 bg-white p-4 shadow-xl md:w-104 dark:border-neutral-950 dark:bg-neutral-800 dark:text-gray-300"
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
