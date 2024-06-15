<template>
  <DrawerRoot
    :open="isOpen"
    :modal="false"
    :should-scale-background="false"
    :no-body-styles="true"
    @update:open="(isOpen) => !isOpen && $emit('close')"
  >
    <DrawerPortal>
      <DrawerContent
        class="bg-gray-100 dark:bg-dark-400 dark:text-gray-300 dark:border-dark-800 flex flex-col rounded-t-[10px] h-full max-h-[96%] fixed bottom-0 left-0 right-0 p-4 pb-0 pt-2 z-10"
        :class="{
          'max-h-1/4': size === '1/4',
          'max-h-1/2': size === '1/2',
          'max-h-3/4': size === '3/4',
          'max-h-full': size === '1',
        }"
      >
        <slot />
      </DrawerContent>
    </DrawerPortal>
  </DrawerRoot>
</template>

<script lang="ts" setup>
import { DrawerContent, DrawerPortal, DrawerRoot } from 'vaul-vue';

defineProps<{
  isOpen: boolean;
  size: '1/4' | '1/2' | '3/4' | '1';
  // disableResize?: boolean;
}>();

defineEmits<{
  (event: 'close'): void;
}>();
</script>
