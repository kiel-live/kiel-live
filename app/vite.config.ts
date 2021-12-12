/* eslint-disable import/no-extraneous-dependencies */
import vue from '@vitejs/plugin-vue';
import path from 'path';
import IconsResolver from 'unplugin-icons/resolver';
import Icons from 'unplugin-icons/vite';
import Components from 'unplugin-vue-components/vite';
import { defineConfig } from 'vite';
import { VitePWA } from 'vite-plugin-pwa';
import WindiCSS from 'vite-plugin-windicss';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    WindiCSS(),
    Icons({ compiler: 'vue3' }),
    Components({
      resolvers: [IconsResolver()],
    }),
    VitePWA({
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
});
