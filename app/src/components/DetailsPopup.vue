<template>
  <div
    v-if="liteMode"
    v-show="isOpen"
    class="h-full w-full flex flex-col z-10 p-4 pt-16 mx-auto max-w-4xl bg-white dark:bg-neutral-800 dark:text-gray-300 dark:border-neutral-950"
  >
    <slot />
  </div>

  <Transition v-else-if="isDesktop" name="fade">
    <div
      v-if="isOpen"
      class="absolute bottom-0 left-0 right-0 flex flex-col z-10 bg-white shadow-right rounded-none w-80 top-0 h-auto transition dark:bg-neutral-800 dark:text-gray-300 dark:border-neutral-950 p-4 pb-0 pt-2 rounded-t-2xl fade"
    >
      <slot />
    </div>
  </Transition>

  <BottomSheet
    v-else
    :is-open="isOpen"
    :size="size"
    @close="$emit('close')"
  >
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
