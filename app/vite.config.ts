import path from 'node:path';
import VueI18nPlugin from '@intlify/unplugin-vue-i18n/vite';
import tailwindcss from '@tailwindcss/vite';
import vue from '@vitejs/plugin-vue';
import IconsResolver from 'unplugin-icons/resolver';
import Icons from 'unplugin-icons/vite';
import Components from 'unplugin-vue-components/vite';
import { VitePWA } from 'vite-plugin-pwa';
import { defineConfig } from 'vitest/config';

// https://vitejs.dev/config/
export default defineConfig({
  envDir: '../',
  build: {
    sourcemap: true,
    rollupOptions: {
      output: {
        // maplibre-gl's worker script references its shared chunk via a
        // hardcoded relative import (`./maplibre-gl-shared.mjs`), so both
        // files must keep their original names to stay resolvable at runtime.
        assetFileNames: (assetInfo) => {
          const names = assetInfo.names ?? (assetInfo.name ? [assetInfo.name] : []);
          return names.includes('maplibre-gl-worker.mjs') || names.includes('maplibre-gl-shared.mjs')
            ? 'assets/[name][extname]'
            : 'assets/[name]-[hash][extname]';
        },
      },
    },
  },
  plugins: [
    vue(),
    VueI18nPlugin({
      include: path.resolve(__dirname, 'src/locales/**'),
    }),
    tailwindcss(),
    Icons({ compiler: 'vue3' }),
    Components({
      resolvers: [IconsResolver()],
    }),
    VitePWA({
      workbox: {
        cleanupOutdatedCaches: true,
      },
      includeAssets: ['favicon.ico'],
      manifest: {
        name: 'Kiel Live',
        short_name: 'Kiel Live',
        description: 'Wo bleibt mein Bus?',
        start_url: './',
        display: 'standalone',
        theme_color: '#2c3e50',
        background_color: '#FFFFFF',
        icons: [
          {
            src: './img/icons/android-chrome-192x192.png',
            sizes: '192x192',
            type: 'image/png',
          },
          {
            src: 'img/icons/android-chrome-512x512.png',
            sizes: '512x512',
            type: 'image/png',
          },
          {
            src: 'img/icons/android-chrome-512x512.png',
            sizes: '512x512',
            type: 'image/png',
            purpose: 'any maskable',
          },
        ],
      },
    }),
  ],
  resolve: {
    alias: {
      '~/': `${path.resolve(__dirname, 'src')}/`,
    },
  },
  test: {
    dir: path.resolve(__dirname, 'src'),
  },
});
