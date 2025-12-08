import { onMounted, onUnmounted, readonly, ref } from 'vue';

export function useDomSize(el: HTMLElement | null) {
  const width = ref(0);
  const height = ref(0);

  function updateSize() {
    if (el) {
      width.value = el.offsetWidth;
      height.value = el.offsetHeight;
    } else {
      width.value = window.innerWidth;
      height.value = window.innerHeight;
    }
  }

  onMounted(() => {
    updateSize();
    window.addEventListener('resize', updateSize);
  });

  onUnmounted(() => {
    window.removeEventListener('resize', updateSize);
  });

  return {
    width: readonly(width),
    height: readonly(height),
  };
}
