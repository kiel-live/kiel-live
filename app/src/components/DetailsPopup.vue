<template>
  <div
    v-if="liteMode"
    v-show="isOpen"
    class="z-10 mx-auto mt-[calc(var(--safe-area-top)+var(--app-bar-space))] flex h-[calc(100%-var(--safe-area-top)-var(--app-bar-space))] w-full max-w-4xl flex-col bg-white px-4 py-2 dark:bg-neutral-800 dark:text-gray-300"
  >
    <slot />
  </div>

  <Transition v-else-if="isDesktop" name="fade">
    <div
      v-if="isOpen"
      class="absolute top-0 bottom-0 left-0 z-10 flex w-[400px] flex-col overflow-y-auto border-r border-gray-100 bg-white px-6 py-4 md:left-20 dark:border-neutral-800 dark:bg-neutral-950 dark:text-gray-300"
    >
      <slot />
    </div>
  </Transition>

  <BottomSheet v-else :is-open="isOpen" :size="size" @close="$emit('close')">
    <slot />
  </BottomSheet>
</template>

<script lang="ts" setup>
import { breakpointsTailwind, useBreakpoints } from '@vueuse/core';
import BottomSheet from '~/components/layout/BottomSheet.vue';
import { useUserSettings } from '~/compositions/useUserSettings';

defineProps<{
  isOpen: boolean;
  size: '3/4' | '1/2' | '1';
}>();

defineEmits<{
  (event: 'close'): void;
}>();

const { liteMode } = useUserSettings();

const breakpoints = useBreakpoints(breakpointsTailwind);
const isDesktop = breakpoints.greater('md');
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
