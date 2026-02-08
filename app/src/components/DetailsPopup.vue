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
      class="shadow-right absolute top-0 bottom-0 left-0 z-10 flex w-80 flex-col bg-white px-4 py-2 dark:bg-neutral-800 dark:text-gray-300"
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
    prevent-close
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
