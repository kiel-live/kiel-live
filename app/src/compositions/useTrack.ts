import { useUpdate } from '~/compositions/useUpdate';

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

export function setupTracker() {
  const { version } = useUpdate();
  const el = document.createElement('script');
  el.setAttribute('src', 'https://boomerang.ju60.de/main.js');
  el.setAttribute('data-website-id', '69027ced-fbe3-442f-8d15-91d499836034');
  el.setAttribute('data-tag', `v-${version.value}`);
  el.setAttribute('data-do-not-track', 'true');
  el.setAttribute('data-domains', 'kiel-live.github.io,kiel.flott-live.de');
  el.async = true;
  el.defer = true;
  document.head.appendChild(el);
}
