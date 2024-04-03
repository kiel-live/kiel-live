<template>
  <DrawerRoot should-scale-background :open="isOpen" @update:open="(isOpen) => !isOpen && $emit('close')">
    <!-- <DrawerTrigger
      class="rounded-full bg-white px-4 py-2.5 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50"
    >
      Open Drawer
    </DrawerTrigger> -->
    <DrawerPortal>
      <DrawerOverlay class="fixed bg-black/40 inset-0" />
      <DrawerContent
        class="bg-gray-100 dark:bg-dark-400 dark:text-gray-300 dark:border-dark-800 flex flex-col rounded-t-[10px] h-full max-h-[96%] fixed bottom-0 left-0 right-0 p-4 pb-0 pt-2"
        :class="{
          'max-h-1/2': size === '1/2',
          'max-h-3/4': size === '3/4',
        }"
      >
        <slot />
      </DrawerContent>
    </DrawerPortal>
  </DrawerRoot>
</template>

<script lang="ts" setup>
import { DrawerContent, DrawerOverlay, DrawerPortal, DrawerRoot, DrawerTrigger } from 'vaul-vue';

// class="absolute bottom-0 left-0 right-0 flex flex-col w-full z-10 bg-white shadow-top md:shadow-right md:rounded-none md:w-80 md:top-0 md:h-auto transition dark:bg-dark-400 dark:text-gray-300 dark:border-dark-800"
// :class="{
//   'overflow-hidden max-h-0': actualSize === 'closed',
//   'h-full md:mx-auto md:w-200 md:shadow-none': actualSize === 'full',
//   'h-1/2': size === '1/2' && actualSize === 'default',
//   'h-3/4': size === '3/4' && actualSize === 'default',
//   'p-4 pb-0 pt-2': actualSize !== 'closed' && actualSize !== 'full',
//   'rounded-t-2xl': actualSize !== 'full',
//   'rounded-none p-4 pt-16': actualSize === 'full',
//   'opacity-80': actualSize === 'closing',
//   fade: !dragging,
// }"
// :style="{ height: isOpen ? (height === undefined ? undefined : `${height}px`) : 0 }"

defineProps<{
  isOpen: boolean;
  size: '3/4' | '1/2' | '1';
  // disableResize?: boolean;
}>();

defineEmits<{
  (event: 'close'): void;
}>();
</script>
