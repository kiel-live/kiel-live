<template>
  <DrawerRoot
    v-model:open="open"
    v-model:active-snap-point="snap"
    :snap-points="snapPoints"
    :modal="false"
    :should-scale-background="false"
    :no-body-styles="true"
    @update:open="(isOpen) => !isOpen && $emit('close')"
  >
    <DrawerPortal>
      <DrawerContent
        class="fixed flex flex-col bg-gray-100 dark:bg-dark-400 dark:text-gray-300 dark:border-dark-800 rounded-t-[10px] max-h-[97%] bottom-0 left-0 right-0 h-full mx-[-1px] z-10"
      >
        <div class="p-4 pb-0 pt-2">
          <div v-if="!disableResize" class="w-full -mt-4 pt-4 pb-4 md:hidden">
            <div class="flex-shrink-0 bg-gray-500 w-12 h-1.5 rounded-full mx-auto" />
          </div>
          <slot />
        </div>
      </DrawerContent>
    </DrawerPortal>
  </DrawerRoot>
</template>

<script lang="ts" setup>
import { DrawerContent, DrawerPortal, DrawerRoot } from 'vaul-vue';
import { computed, ref, watch } from 'vue';

const props = defineProps<{
  size: '1/2' | '1';
  disableResize?: boolean;
}>();

defineEmits<{
  (event: 'close'): void;
}>();

const open = defineModel<boolean>('open', { required: true });

const snapPoints = computed(() => {
  if (props.size === '1/2') return [0.2, 0.5];
  return [0.2, 0.5, 1];
});

const snap = ref<number | string | null>(0.5);
watch([() => props.size, open], () => {
  snap.value = 0.5;
});
</script>
