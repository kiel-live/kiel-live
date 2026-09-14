import { ref } from 'vue';

export function useShare() {
  const copied = ref(false);
  let resetTimeout: ReturnType<typeof setTimeout> | undefined;

  async function share(data: { title?: string; url: string }) {
    if (navigator.share) {
      try {
        await navigator.share(data);
        return;
      } catch (error) {
        // user cancelled the native share sheet
        if (error instanceof Error && error.name === 'AbortError') {
          return;
        }
      }
    }

    await navigator.clipboard.writeText(data.url);
    copied.value = true;
    clearTimeout(resetTimeout);
    resetTimeout = setTimeout(() => {
      copied.value = false;
    }, 2000);
  }

  return { share, copied };
}
