import { ref } from 'vue';

export const usePrefersColorSchemeDark = () => {
  if (!window.matchMedia) {
    return ref(false);
  }

  const prefersColorSchemeDark = ref(window.matchMedia('(prefers-color-scheme: dark)').matches);
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
    prefersColorSchemeDark.value = e.matches;
  });

  return prefersColorSchemeDark;
};
