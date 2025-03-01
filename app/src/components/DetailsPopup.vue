<template>
  <div
    v-show="isOpen"
    ref="popup"
    class="absolute bottom-0 left-0 right-0 flex flex-col w-full z-10 bg-white shadow-top md:shadow-right md:rounded-none md:w-80 md:top-0 md:h-auto transition dark:bg-dark-400 dark:text-gray-300 dark:border-dark-800"
    :class="{
      'overflow-hidden max-h-0': !isOpen,
      'md:mx-auto md:w-200 md:shadow-none': size === 'full',
      'p-4 pb-0 pt-2': !disableResize,
      'rounded-t-2xl': !disableResize,
      'rounded-none p-4 pt-16': disableResize,
      fade: !dragging,
    }"
    :style="{ height: isOpen ? `${height}px` : 0 }"
    @touchmove="move"
    @touchend="drop"
  >
    <div v-if="!disableResize" class="w-full -mt-4 pt-4 pb-4 md:hidden" @touchstart="drag">
      <div class="flex-shrink-0 bg-gray-500 w-12 h-1.5 rounded-full mx-auto" />
    </div>
    <slot />
  </div>
</template>

<script lang="ts" setup>
import { computed, ref, toRef, useTemplateRef } from 'vue';

const props = defineProps<{
  isOpen: boolean;
  disableResize?: boolean;
}>();

const size = defineModel<'full' | 'half' | 'minimized'>('size', {
  required: true,
});

const popup = useTemplateRef<HTMLDivElement>('popup');

const dragging = ref(false);
const isOpen = toRef(props, 'isOpen');
const disableResize = toRef(props, 'disableResize');

const dragStartPositon = ref(0);
const dragCurrentPosition = ref(0);

/**
 * The top offset given by the status bar in the android app
 */
const topOffset = window.topOffset || 0;

const dragDistance = computed(() => dragCurrentPosition.value - dragStartPositon.value);

const height = computed(() => {
  if (!isOpen.value) {
    return 0;
  }

  if (dragging.value) {
    return window.innerHeight - dragCurrentPosition.value;
  }

  if (size.value === 'full') {
    return window.innerHeight - topOffset;
  }

  if (size.value === 'half') {
    return (window.innerHeight - topOffset) / 2;
  }

  if (size.value === 'minimized') {
    return 65;
  }

  // TODO
  return window.innerHeight;
});

/**
 * Which state is next when dragging stops
 */
const dragState = computed<'full' | 'half' | 'minimized' | null>(() => {
  if (!dragging.value) {
    return null;
  }

  if (size.value === 'full') {
    if (dragDistance.value > 40) {
      return 'half';
    }
  }

  if (size.value === 'half') {
    if (dragDistance.value > 40) {
      return 'minimized';
    }

    if (dragDistance.value < 40) {
      return 'full';
    }
  }

  if (size.value === 'minimized') {
    if (dragDistance.value < 40) {
      return 'half';
    }
  }

  return null;
});

function drag(e: TouchEvent) {
  if (disableResize.value) {
    return;
  }

  dragging.value = true;
  dragStartPositon.value = e.touches[0].clientY;
  dragCurrentPosition.value = e.touches[0].clientY;
}

function move(e: TouchEvent) {
  if (!dragging.value) {
    return;
  }

  dragCurrentPosition.value = e.touches[0].clientY;
}

function drop() {
  if (!dragging.value) {
    return;
  }

  if (dragState.value !== null) {
    size.value = dragState.value;
  }

  dragging.value = false;
}
</script>

<style scoped>
.fade {
  transition: height 0.15s ease;
}
</style>
