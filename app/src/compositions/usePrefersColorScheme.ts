import { ref } from 'vue';

function setDarkMode(enabled: boolean) {
  if (enabled) {
    document.documentElement.classList.add('dark');
  } else {
    document.documentElement.classList.remove('dark');
  }
}

const prefersColorSchemeDark = ref<boolean>();

export const usePrefersColorSchemeDark = () => {
  if (prefersColorSchemeDark.value !== undefined) {
    return prefersColorSchemeDark;
  }

  if (!window.matchMedia) {
    prefersColorSchemeDark.value = false;
    return prefersColorSchemeDark;
  }

  prefersColorSchemeDark.value = window.matchMedia('(prefers-color-scheme: dark)').matches;
  setDarkMode(prefersColorSchemeDark.value);
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
    prefersColorSchemeDark.value = e.matches;
    setDarkMode(prefersColorSchemeDark.value);
  });

  return prefersColorSchemeDark;
};
