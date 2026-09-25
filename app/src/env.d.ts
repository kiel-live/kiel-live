/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_BUILD_DATE: string;
  readonly VITE_NATS_URL: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}

// unplugin-icons ?raw imports (raw SVG markup) vs. the project-wide Vue
// component typing for `~icons/*` (see tsconfig "types")
declare module '~icons/*?raw' {
  const svg: string;
  export default svg;
}
