/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_BUILD_DATE: string;
  readonly VITE_NATS_URL: string;
  readonly VITE_HUB_URL: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
