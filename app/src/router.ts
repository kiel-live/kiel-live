import type { Component } from 'vue';
import type { RouteRecordRaw } from 'vue-router';
import { createRouter, createWebHistory } from 'vue-router';

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
    path: '/settings/contact',
    name: 'settings-contact',
    component: (): Component => import('~/views/settings/Contact.vue'),
    meta: { settings: true },
  },
  {
    path: '/contact',
    name: 'contact',
    redirect: { name: 'settings-contact' },
  },
  {
    path: '/dev',
    name: 'dev',
    component: (): Component => import('~/views/Development.vue'),
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
