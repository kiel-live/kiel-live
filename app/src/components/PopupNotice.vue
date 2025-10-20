<template>
  <div
    v-if="show"
    class="flex fixed top-0 left-0 right-0 bottom-0 bg-zinc-900 bg-opacity-80 z-1000 justify-center items-center"
    @click="$emit('close')"
  >
    <div
      class="m-2 flex flex-col rounded-md p-4 bg-white border border-gray-200 shadow-xl z-20 md:w-104 dark:bg-zinc-900 dark:text-gray-300 dark:border-zinc-950"
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
