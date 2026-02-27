import type { Ref } from 'vue';
import { onMounted, onUnmounted, readonly, ref, watch } from 'vue';

export function useDomSize(el: Ref<HTMLElement | null | undefined>) {
  const width = ref(0);
  const height = ref(0);

  function updateSize() {
    if (!el.value) {
      return;
    }

    width.value = el.value.offsetWidth;
    height.value = el.value.offsetHeight;
  }

  onMounted(() => {
    updateSize();
    window.addEventListener('resize', updateSize);
  });

  onUnmounted(() => {
    window.removeEventListener('resize', updateSize);
  });

  watch(el, updateSize);

  return {
    width: readonly(width),
    height: readonly(height),
  };
}
