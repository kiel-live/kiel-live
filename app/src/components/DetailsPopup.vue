<template>
  <DrawerRoot
    v-model:open="open"
    :modal="false"
    :should-scale-background="false"
    :no-body-styles="true"
    :snap-points="snapPoints"
    :active-snap-point="snap"
    @update:open="(isOpen) => !isOpen && $emit('close')"
  >
    <DrawerPortal>
      <DrawerContent
        class="bg-gray-100 dark:bg-dark-400 dark:text-gray-300 dark:border-dark-800 flex flex-col rounded-t-[10px] h-full max-h-[96%] fixed bottom-0 left-0 right-0 p-4 pb-0 pt-2 z-10"
        :class="{
          'max-h-1/4': size === '1/4',
          'max-h-1/2': size !== '1/4',
        }"
      >
        <div v-if="!disableResize" class="w-full -mt-4 pt-4 pb-4 md:hidden">
          <div class="flex-shrink-0 bg-gray-500 w-12 h-1.5 rounded-full mx-auto" />
        </div>
        <slot />
      </DrawerContent>
    </DrawerPortal>
  </DrawerRoot>
</template>

<script lang="ts" setup>
import { DrawerContent, DrawerPortal, DrawerRoot } from 'vaul-vue';
import { ref } from 'vue';

const open = defineModel<boolean>('open', { required: true });

defineProps<{
  size: '1/4' | '1/2' | '3/4' | '1';
  disableResize?: boolean;
}>();

const snapPoints = [0, 0.2, 0.5, 0.75, 1];
const snap = ref<number | string | null>(snapPoints[1]);

defineEmits<{
  (event: 'close'): void;
}>();
</script>
