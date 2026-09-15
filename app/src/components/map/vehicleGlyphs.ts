import type { VehicleType } from '~/api/types';

// Inline SVG bodies for the same icon set already used elsewhere in the app
// (see VehiclePopup.vue), so map markers match the icons shown in popups.
// Extracted once from the @iconify-json/{mdi,ph,carbon} packages to avoid
// bundling those (multi-MB) icon sets at runtime.
interface VehicleGlyph {
  viewBox: number;
  body: string;
}

const glyphs: Record<VehicleType, VehicleGlyph> = {
  bus: {
    viewBox: 24,
    body: '<path fill="currentColor" d="M18 11H6V6h12m-1.5 11a1.5 1.5 0 0 1-1.5-1.5a1.5 1.5 0 0 1 1.5-1.5a1.5 1.5 0 0 1 1.5 1.5a1.5 1.5 0 0 1-1.5 1.5m-9 0A1.5 1.5 0 0 1 6 15.5A1.5 1.5 0 0 1 7.5 14A1.5 1.5 0 0 1 9 15.5A1.5 1.5 0 0 1 7.5 17M4 16c0 .88.39 1.67 1 2.22V20a1 1 0 0 0 1 1h1a1 1 0 0 0 1-1v-1h8v1a1 1 0 0 0 1 1h1a1 1 0 0 0 1-1v-1.78c.61-.55 1-1.34 1-2.22V6c0-3.5-3.58-4-8-4s-8 .5-8 4z"/>',
  },
  bike: {
    viewBox: 32,
    body: '<path fill="currentColor" d="M26 16c-.088 0-.173.01-.26.013L24.237 9H28V7h-5a1 1 0 0 0-.978 1.21L22.62 11H12.387l-1-3H14V6H7v2h2.28l1.041 3.123l-2.57 5.14A6 6 0 1 0 11.91 23h2.61a2 2 0 0 0 1.562-.75l7.058-8.824l.644 3.004A5.99 5.99 0 1 0 26 16M6 26a4 4 0 1 1 .836-7.91l-1.73 3.463l.009.004A1 1 0 0 0 5 22a.993.993 0 0 0 1.885.443l.01.004L8.618 19A3.984 3.984 0 0 1 6 26m5.91-5a6 6 0 0 0-2.373-3.836l1.678-3.358L13.613 21Zm3.458-1.06L13.054 13h7.865ZM26 26a3.988 3.988 0 0 1-1.786-7.56l.808 3.77l.02-.004A.986.986 0 0 0 26 23a1 1 0 0 0 1-1a1 1 0 0 0-.041-.206l.02-.004l-.81-3.773A3.993 3.993 0 0 1 26 26"/>',
  },
  car: {
    viewBox: 256,
    body: '<path fill="currentColor" d="M240 104h-10.8l-27.78-62.5A16 16 0 0 0 186.8 32H69.2a16 16 0 0 0-14.62 9.5L26.8 104H16a8 8 0 0 0 0 16h8v80a16 16 0 0 0 16 16h24a16 16 0 0 0 16-16v-16h96v16a16 16 0 0 0 16 16h24a16 16 0 0 0 16-16v-80h8a8 8 0 0 0 0-16M69.2 48h117.6l24.89 56H44.31ZM64 200H40v-16h24Zm128 0v-16h24v16Zm24-32H40v-48h176ZM56 144a8 8 0 0 1 8-8h16a8 8 0 0 1 0 16H64a8 8 0 0 1-8-8m112 0a8 8 0 0 1 8-8h16a8 8 0 0 1 0 16h-16a8 8 0 0 1-8-8"/>',
  },
  'e-scooter': {
    viewBox: 256,
    body: '<path fill="currentColor" d="M212 136c-1.18 0-2.35.06-3.51.17l-32.9-98.7A8 8 0 0 0 168 32h-32a8 8 0 0 0 0 16h26.23l17.44 52.31L124.21 168H79.77a36 36 0 1 0-1.83 16H128a8 8 0 0 0 6.19-2.93l51.46-62.81l7.66 23A36 36 0 1 0 212 136M44 192a20 20 0 1 1 20-20a20 20 0 0 1-20 20m168 0a20 20 0 1 1 20-20a20 20 0 0 1-20 20"/>',
  },
  ferry: {
    viewBox: 24,
    body: '<path fill="currentColor" d="M6 6h12v3.96L12 8L6 9.96M3.94 19H4c1.6 0 3-.88 4-2c1 1.12 2.4 2 4 2s3-.88 4-2c1 1.12 2.4 2 4 2h.05l1.9-6.69c.08-.25.05-.53-.06-.77c-.13-.24-.34-.42-.6-.5L20 10.62V6a2 2 0 0 0-2-2h-3V1H9v3H6a2 2 0 0 0-2 2v4.62l-1.29.42c-.26.08-.47.26-.6.5c-.11.24-.14.52-.06.77M20 21c-1.39 0-2.78-.47-4-1.33c-2.44 1.71-5.56 1.71-8 0C6.78 20.53 5.39 21 4 21H2v2h2c1.37 0 2.74-.35 4-1c2.5 1.3 5.5 1.3 8 0c1.26.65 2.62 1 4 1h2v-2z"/>',
  },
  train: {
    viewBox: 32,
    body: '<path fill="currentColor" d="M30 25H2v2h2v2h2v-2h5v2h2v-2h5v2h2v-2h5v2h2v-2h3zM8 16H2v-2h6v-2H2v-2h6a2 2 0 0 1 2 2v2a2 2 0 0 1-2 2"/><path fill="currentColor" d="m28.55 14.23l-8.58-7.864A8.98 8.98 0 0 0 13.888 4H2v2h10v4a2 2 0 0 0 2 2h9.157l4.041 3.705A2.472 2.472 0 0 1 25.528 20H2v2h23.527a4.473 4.473 0 0 0 3.023-7.77M14 10V6.005a6.98 6.98 0 0 1 4.618 1.835L20.975 10Z"/>',
  },
  subway: {
    viewBox: 256,
    body: '<path fill="currentColor" d="M224 96v112a8 8 0 0 1-16 0V96a56.06 56.06 0 0 0-56-56h-48a56.06 56.06 0 0 0-56 56v112a8 8 0 0 1-16 0V96a72.08 72.08 0 0 1 72-72h48a72.08 72.08 0 0 1 72 72m-40 0v72a24 24 0 0 1-19.29 23.53l2.45 4.89a8 8 0 0 1-14.32 7.16L147.06 192h-38.12l-5.78 11.58a8 8 0 0 1-14.32-7.16l2.45-4.89A24 24 0 0 1 72 168V96a24 24 0 0 1 24-24h64a24 24 0 0 1 24 24m-96 0v48h80V96a8 8 0 0 0-8-8H96a8 8 0 0 0-8 8m32 64v16h16v-16Zm-24 16h8v-16H88v8a8 8 0 0 0 8 8m72-8v-8h-16v16h8a8 8 0 0 0 8-8"/>',
  },
  tram: {
    viewBox: 256,
    body: '<path fill="currentColor" d="M184 48h-48V24h32a8 8 0 0 0 0-16H88a8 8 0 0 0 0 16h32v24H72a32 32 0 0 0-32 32v104a32 32 0 0 0 32 32h8l-14.4 19.2a8 8 0 1 0 12.8 9.6L100 216h56l21.6 28.8a8 8 0 1 0 12.8-9.6L176 216h8a32 32 0 0 0 32-32V80a32 32 0 0 0-32-32M72 64h112a16 16 0 0 1 16 16v40H56V80a16 16 0 0 1 16-16m112 136H72a16 16 0 0 1-16-16v-48h144v48a16 16 0 0 1-16 16m-88-28a12 12 0 1 1-12-12a12 12 0 0 1 12 12m88 0a12 12 0 1 1-12-12a12 12 0 0 1 12 12"/>',
  },
  moped: {
    viewBox: 24,
    body: '<path fill="currentColor" d="M19 15c.55 0 1 .45 1 1s-.45 1-1 1s-1-.45-1-1s.45-1 1-1m0-2c-1.66 0-3 1.34-3 3s1.34 3 3 3s3-1.34 3-3s-1.34-3-3-3m-9-7H5v2h5zm7-1h-3v2h3v2.65L13.5 14H10V9H6c-2.21 0-4 1.79-4 4v3h2c0 1.66 1.34 3 3 3s3-1.34 3-3h4.5l4.5-5.65V7a2 2 0 0 0-2-2M7 17c-.55 0-1-.45-1-1h2c0 .55-.45 1-1 1"/>',
  },
  'e-moped': {
    viewBox: 24,
    body: '<path fill="currentColor" d="M19 5c0-1.1-.9-2-2-2h-3v2h3v2.65L13.5 12H10V7H6c-2.21 0-4 1.79-4 4v3h2c0 1.66 1.34 3 3 3s3-1.34 3-3h4.5L19 8.35zM7 15c-.55 0-1-.45-1-1h2c0 .55-.45 1-1 1M5 4h5v2H5zm14 7c-1.66 0-3 1.34-3 3s1.34 3 3 3s3-1.34 3-3s-1.34-3-3-3m0 4c-.55 0-1-.45-1-1s.45-1 1-1s1 .45 1 1s-.45 1-1 1M7 20h4v-2l6 3h-4v2z"/>',
  },
};

export function vehicleGlyphDataUrl(type: VehicleType, color: string): string {
  const glyph = glyphs[type];
  const svg = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 ${glyph.viewBox} ${glyph.viewBox}" style="color:${color}">${glyph.body}</svg>`;
  return `data:image/svg+xml;charset=utf-8,${encodeURIComponent(svg)}`;
}
