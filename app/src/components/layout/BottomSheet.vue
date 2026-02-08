<template>
  <div
    v-show="isOpen"
    class="absolute bottom-0 left-0 right-0 flex flex-col w-full px-4 pb-0 pt-2 z-10 transition bg-white dark:bg-neutral-800 dark:text-gray-300 dark:border-neutral-950"
    :class="{
      'overflow-hidden max-h-0': actualSize === 'closed',
      'max-h-[calc(100%-var(--safe-area-top)-var(--app-bar-space))]': actualSize !== 'closed',
      'h-full shadow-none': actualSize === 'full',
      'h-1/2': size === '1/2' && actualSize === 'default',
      'h-3/4': size === '3/4' && actualSize === 'default',
      'rounded-t-2xl shadow-top': actualSize !== 'full',
      'rounded-none': actualSize === 'full',
      'opacity-80': actualSize === 'closing',
      fade: !dragging,
    }"
    :style="{ height: isOpen ? (height === undefined ? undefined : `${height}px`) : 0 }"
  >
    <button type="button" class="w-full -mt-4 pt-4 pb-4 touch-none" :title="$t('drag_to_resize')" @pointerdown="drag">
      <div class="shrink-0 bg-gray-500 w-12 h-1.5 rounded-full mx-auto" />
    </button>
    <slot />
  </div>
</template>

<script lang="ts" setup>
import { computed, onUnmounted, ref, toRef } from 'vue';

const props = defineProps<{
  isOpen: boolean;
  size: '3/4' | '1/2' | '1';
}>();

const emit = defineEmits<{
  (event: 'close'): void;
}>();

const dragging = ref(false);
const height = ref<number>();
const isOpen = toRef(props, 'isOpen');
const size = toRef(props, 'size');

const actualSize = computed(() => {
  if (size.value === '1') {
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
