<template>
  <div
    class="
      absolute
      bottom-0
      left-0
      right-0
      flex flex-col
      w-full
      z-10
      bg-white
      shadow-all
      md:rounded-none md:w-80 md:top-0 md:h-auto
      transition
      dark:bg-dark-400 dark:text-gray-300 dark:border-dark-800
    "
    :class="{
      'overflow-hidden max-h-0': size === 'closed',
      'h-48': size === 'default',
      'p-4 pt-2': size !== 'closed' && size !== 'full',
      'rounded-t-2xl': size !== 'full',
      'rounded-none p-4 pt-16': size === 'full',
      'opacity-80': size === 'closing',
    }"
    :style="{ height: isOpen ? (height === undefined ? undefined : `${height}px`) : 0 }"
    @touchstart="drag"
    @touchmove="move"
    @touchend="drop"
  >
    <div class="bg-gray-500 w-12 h-1.5 mb-4 rounded-full mx-auto md:hidden" v-show="size !== 'full'" />
    <slot />
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, ref, toRef } from 'vue';

export default defineComponent({
  name: 'DetailPopup',

  emits: {
    close: () => true,
  },

  props: {
    isOpen: {
      type: Boolean,
      required: true,
    },
  },

  setup(props, { emit }) {
    const dragging = ref(false);
    const height = ref<number>();
    const isOpen = toRef(props, 'isOpen');
    const size = computed(() => {
      if (!isOpen.value) {
        return 'closed';
      }

      if (dragging.value) {
        if (height.value === undefined) {
          throw new Error('hmm');
        }

        const percentage = height.value / window.innerHeight;
        if (percentage > 0.4) {
          return 'maximizing';
        }

        if (percentage < 0.2) {
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

    function drag(e: TouchEvent) {
      dragging.value = true;
      height.value = window.innerHeight - e.touches[0].clientY;
    }

    function move(e: TouchEvent) {
      if (!dragging.value) {
        return;
      }
      height.value = window.innerHeight - e.touches[0].clientY;
    }

    function drop(e: TouchEvent) {
      if (size.value === 'maximizing') {
        height.value = window.innerHeight;
      } else if (size.value === 'closing') {
        height.value = undefined;
        emit('close');
      } else if (size.value === 'defaulting') {
        height.value = undefined;
      }
      dragging.value = false;
    }

    return { drag, move, drop, size, height };
  },
});
</script>
