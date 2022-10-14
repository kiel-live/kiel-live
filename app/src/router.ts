import { Component } from 'vue';
import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';

import { useUserSettings } from '~/compositions/useUserSettings';

let firstStartOfApp = true;

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'home',
    component: (): Component => import('~/views/Home.vue'),
  },
  {
    path: '/map/:markerType/:markerId',
    name: 'map-marker',
    component: (): Component => import('~/views/Home.vue'),
    props: true,
  },
  {
    path: '/search',
    name: 'search',
    component: (): Component => import('~/views/Home.vue'),
  },
  {
    path: '/favorites',
    name: 'favorites',
    component: (): Component => import('~/views/Home.vue'),
  },
  {
    path: '/settings/about',
    name: 'settings-about',
    component: (): Component => import('~/views/settings/About.vue'),
    meta: { settings: true },
  },
  {
    path: '/settings/changelog',
    name: 'settings-changelog',
    component: (): Component => import('~/views/settings/Changelog.vue'),
    meta: { settings: true },
  },
  {
    path: '/settings/settings',
    name: 'settings-settings',
    component: (): Component => import('~/views/settings/Settings.vue'),
    meta: { settings: true },
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: (): Component => import('~/views/NotFound.vue'),
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

const { liteMode } = useUserSettings();

router.beforeEach((to, from, next): void => {
  if (to.name === 'home' && (firstStartOfApp || liteMode.value)) {
    firstStartOfApp = false;
    next({ name: 'favorites' });
    return;
  }

  firstStartOfApp = false;
  next();
});

export default router;
