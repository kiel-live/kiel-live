<template>
  <div
    v-if="liteMode"
    v-show="isOpen"
    class="absolute top-0 bottom-0 left-0 right-0 flex flex-col z-10 bg-white dark:bg-dark-400 dark:text-gray-300 dark:border-dark-800 p-4 pt-16 mx-auto max-w-4xl"
  >
    <slot />
  </div>

  <Transition v-else-if="isDesktop" name="fade">
    <div
      v-if="isOpen"
      class="absolute bottom-0 left-0 right-0 flex flex-col w-full z-10 bg-white shadow-top md:shadow-right md:rounded-none md:w-80 md:top-0 md:h-auto transition dark:bg-dark-400 dark:text-gray-300 dark:border-dark-800 h-1/2 p-4 pb-0 pt-2 rounded-t-2xl fade"
    >
      <slot />
    </div>
  </Transition>

  <BottomSheet
    v-else
    :current-snap-point="currentSnapPoint"
    :is-open="isOpen"
    :snap-points="snapPoints"
    :show-backdrop="false"
    @update:current-snap-point="$emit('update:current-snap-point', $event)"
    @close="$emit('close')"
  >
    <slot />
  </BottomSheet>
</template>

<script lang="ts" setup>
import { breakpointsTailwind, useBreakpoints } from '@vueuse/core';
import BottomSheet from '~/components/layout/BottomSheet.vue';
import { useUserSettings } from '~/compositions/useUserSettings';

type SnapPoint = number | string;

defineProps<{
  isOpen: boolean;
  disableResize?: boolean;
  snapPoints?: SnapPoint[];
  currentSnapPoint?: SnapPoint;
}>();

defineEmits<{
  (event: 'close'): void;
  (event: 'update:current-snap-point', snapPoint: SnapPoint): void;
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
