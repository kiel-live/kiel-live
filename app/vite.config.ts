import path from 'node:path';
import VueI18nPlugin from '@intlify/unplugin-vue-i18n/vite';
import tailwindcss from '@tailwindcss/vite';
import vue from '@vitejs/plugin-vue';
import IconsResolver from 'unplugin-icons/resolver';
import Icons from 'unplugin-icons/vite';
import Components from 'unplugin-vue-components/vite';
import { loadEnv } from 'vite';
import csp from 'vite-plugin-csp-guard';
import { VitePWA } from 'vite-plugin-pwa';
import { defineConfig } from 'vitest/config';

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, '../');
  return {
    envDir: '../',
    build: {
      sourcemap: true,
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
          globPatterns: ['**/*.{js,css,html,png}'],
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
      csp({
        dev: {
          run: true, // Run the plugin in dev mode
          outlierSupport: ['vue', 'tailwind'],
        },
        policy: {
          'script-src-elem': ["'self'", 'https://boomerang.ju60.de'],
          'style-src-elem': ["'self'", "'unsafe-inline'"],
          'connect-src': ["'self'", env.VITE_NATS_URL, 'https://tiles.ju60.de'],
          'worker-src': ["'self'", 'blob:'],
        },
        build: {
          sri: true,
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
  };
});
