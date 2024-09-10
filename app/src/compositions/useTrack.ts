export function useTrack() {
  const { umami } = globalThis as {
    umami?: {
      track: (event: string, data?: any) => void;
    };
  };

  async function track(event: string, data?: any) {
    umami?.track(event, data);
  }

  return { track };
}
