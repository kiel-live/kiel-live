export function useTrack() {
  const umami = (
    globalThis as {
      umami?: {
        track: (event: string, data?: any) => void;
      };
    }
  ).umami;

  async function track(event: string, data?: any) {
    umami?.track(event, data);
    console.log('Tracked event:', event, data);
  }

  return { track };
}
