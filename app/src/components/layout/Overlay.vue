<template>
  <transition name="fade">
    <div v-if="isOpen" class="absolute top-0 left-0 flex h-full w-full min-h-0 bg-white dark:bg-dark-400 z-10">
      <slot />
    </div>
  </transition>
</template>

<script lang="ts">
import { defineComponent } from 'vue';

export default defineComponent({
  name: 'Overlay',

  props: {
    isOpen: {
      type: Boolean,
    },
  },
});
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition:
    opacity 0.3s ease,
    transform 0.3s ease,
    height 0.3s ease;
}

.fade-enter-to,
.fade-leave-from {
  transform: translateY(0%);
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-5%);
  height: 105% !important;
}
</style>
