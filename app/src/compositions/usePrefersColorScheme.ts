import { ref, watch } from 'vue';

function setDarkMode(enabled: boolean) {
  if (enabled) {
    document.documentElement.classList.add('dark');
  } else {
    document.documentElement.classList.remove('dark');
  }
}

export const usePrefersColorSchemeDark = () => {
  if (!window.matchMedia) {
    return ref(false);
  }

  const prefersColorSchemeDark = ref(window.matchMedia('(prefers-color-scheme: dark)').matches);
  setDarkMode(prefersColorSchemeDark.value);
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
    prefersColorSchemeDark.value = e.matches;
    setDarkMode(prefersColorSchemeDark.value);
  });

  return prefersColorSchemeDark;
};
