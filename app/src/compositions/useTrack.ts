interface TrackData {
  hostname: string;
  language: string;
  referrer: string;
  screen: string;
  title: string;
  url: string;
  website: string;
  data: Record<string, string | number | boolean>;
  [key: string]: unknown;
}

interface TrackFunction {
  // Track pageview
  (): void;

  // Track custom event
  (eventName: string, eventData?: Record<string, unknown>): void;

  // Modify default tracking properties
  (callback: (props: TrackData) => TrackData): void;
}

export function useTrack() {
  const { umami } = globalThis as {
    umami?: {
      track: TrackFunction;
    };
  };

  const noop = () => {};

  return { track: (umami?.track ?? noop) as TrackFunction };
}
