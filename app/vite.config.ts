import path from 'node:path';
import VueI18nPlugin from '@intlify/unplugin-vue-i18n/vite';
import vue from '@vitejs/plugin-vue';
import IconsResolver from 'unplugin-icons/resolver';
import Icons from 'unplugin-icons/vite';
import Components from 'unplugin-vue-components/vite';
import { defineConfig } from 'vite';
import { VitePWA } from 'vite-plugin-pwa';
import { viteStaticCopy } from 'vite-plugin-static-copy';
import WindiCSS from 'vite-plugin-windicss';

// https://vitejs.dev/config/
export default defineConfig({
  envDir: '../',
  build: {
    sourcemap: true,
  },
  plugins: [
    vue(),
    VueI18nPlugin({
      include: path.resolve(__dirname, 'src/locales/**'),
    }),
    WindiCSS(),
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
    viteStaticCopy({
      targets: [
        {
          src: 'public/map-styles/*.json',
          dest: 'map-styles',
          transform: (contents) => JSON.stringify(JSON.parse(contents.toString())),
        },
      ],
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
