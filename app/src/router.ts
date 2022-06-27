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
    props: true,
  },
  {
    path: '/favorites',
    name: 'favorites',
    component: (): Component => import('~/views/Home.vue'),
    props: true,
  },
  {
    path: '/settings',
    name: 'settings',
    component: (): Component => import('~/views/Settings.vue'),
    props: true,
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

const { useLiteMode } = useUserSettings();

router.beforeEach((to, from, next): void => {
  if (to.name === 'home' && (firstStartOfApp || useLiteMode.value)) {
    firstStartOfApp = false;
    next({ name: 'favorites' });
    return;
  }

  firstStartOfApp = false;
  next();
});

export default router;
