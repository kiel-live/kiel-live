<template>
  <div
    v-show="isOpen"
    class="md:shadow-right absolute right-0 bottom-0 left-0 z-10 flex w-full flex-col bg-white px-4 pt-2 pb-0 transition md:top-0 md:h-auto md:w-80 md:rounded-none dark:border-neutral-950 dark:bg-neutral-800 dark:text-gray-300"
    :class="{
      'max-h-0 overflow-hidden': actualSize === 'closed',
      'max-md:max-h-[calc(100%-var(--safe-area-top)-var(--app-bar-space))]': actualSize !== 'closed',
      'h-full shadow-none md:mx-auto md:w-200': actualSize === 'full',
      'h-1/2': size === '1/2' && actualSize === 'default',
      'h-3/4': size === '3/4' && actualSize === 'default',
      'shadow-top rounded-t-2xl': actualSize !== 'full' || !disableResize,
      'rounded-none': actualSize === 'full' && disableResize,
      'opacity-80': actualSize === 'closing',
      fade: !dragging,
    }"
    :style="{ height: isOpen ? (height === undefined ? undefined : `${height}px`) : 0 }"
  >
    <button
      v-if="!disableResize"
      type="button"
      class="-mt-4 w-full touch-none pt-4 pb-4 md:hidden"
      :title="$t('drag_to_resize')"
      @pointerdown="drag"
    >
      <div class="mx-auto h-1.5 w-12 shrink-0 rounded-full bg-gray-500" />
    </button>
    <slot />
  </div>
</template>

<script lang="ts" setup>
import { computed, onUnmounted, ref, toRef } from 'vue';

const props = defineProps<{
  isOpen: boolean;
  size: '3/4' | '1/2' | '1';
  disableResize?: boolean;
}>();

const emit = defineEmits<{
  (event: 'close'): void;
}>();

const dragging = ref(false);
const height = ref<number>();
const isOpen = toRef(props, 'isOpen');
const size = toRef(props, 'size');
const disableResize = toRef(props, 'disableResize');

const actualSize = computed(() => {
  if (disableResize.value && size.value === '1') {
    return 'full';
  }

  if (!isOpen.value) {
    return 'closed';
  }

  if (dragging.value) {
    if (height.value === undefined) {
      return 'closed';
    }

    const percentage = height.value / window.innerHeight;
    if ((size.value === '1/2' && percentage > 0.6) || (size.value === '3/4' && percentage > 0.85)) {
      return 'maximizing';
    }

    if ((size.value === '1/2' && percentage < 0.4) || (size.value === '3/4' && percentage < 0.65)) {
      return 'closing';
    }

    return 'defaulting';
  }

  if (height.value === 0) {
    return 'closed';
  }

  if (height.value === window.innerHeight) {
    return 'full';
  }

  return 'default';
});

function drag(e: PointerEvent) {
  if (disableResize.value) {
    return;
  }

  dragging.value = true;
  height.value = window.innerHeight - e.clientY;

  window.addEventListener('pointermove', move);
  window.addEventListener('pointerup', drop);
  window.addEventListener('pointercancel', drop);
}

function move(e: PointerEvent) {
  if (!dragging.value) {
    return;
  }
  height.value = window.innerHeight - e.clientY;
}

function removeEventListeners() {
  window.removeEventListener('pointermove', move);
  window.removeEventListener('pointerup', drop);
  window.removeEventListener('pointercancel', drop);
}

function drop() {
  if (!dragging.value) {
    return;
  }

  if (actualSize.value === 'maximizing') {
    height.value = window.innerHeight;
  } else if (actualSize.value === 'closing') {
    height.value = undefined;
    emit('close');
  } else if (actualSize.value === 'defaulting') {
    height.value = undefined;
  }

  dragging.value = false;

  removeEventListeners();
}

onUnmounted(() => {
  removeEventListeners();
});
</script>

<style scoped>
.fade {
  transition: height 0.15s ease;
}
</style>
