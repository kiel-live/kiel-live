export const feedbackMail = atob('YW5kcm9pZEBqdTYwLmRl'); // email as base64
export const buildDate = import.meta.env.VITE_BUILD_DATE;
export const analyticsUrl = 'https://boomerang.ju60.de/share/z8KAHmGY/Kiel%20Live';
export const darkMapStyle = '/map-style/dark.json';
export const lightMapStyle = '/map-style/light.json';
export const natsServerUrl = import.meta.env.VITE_NATS_URL;
export const DEBUG = (globalThis?.window as { DEBUG?: boolean })?.DEBUG || import.meta.env.DEV;
