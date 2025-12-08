export const feedbackMail = atob('YW5kcm9pZEBqdTYwLmRl'); // email as base64
export const buildDate = import.meta.env.VITE_BUILD_DATE;
export const analyticsUrl = 'https://boomerang.ju60.de/share/z8KAHmGY/Kiel%20Live';
const tileServer = 'https://tiles.immich.cloud/v1/style';
export const darkMapStyle = `${tileServer}/dark.json`;
export const lightMapStyle = `${tileServer}/light.json`;
export const natsServerUrl = import.meta.env.VITE_NATS_URL;
export const DEBUG = (globalThis?.window as { DEBUG?: boolean })?.DEBUG || import.meta.env.DEV;
