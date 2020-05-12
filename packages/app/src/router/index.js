import Vue from 'vue';
import VueRouter from 'vue-router';
import Store from '@/store';

Vue.use(VueRouter);

// route level code-splitting
// this generates a separate chunk (name.[hash].js) for this route
// which is lazy-loaded when the route is visited.
// function loadView(...view) {
//   return () => import(/* webpackChunkName: "view-[request]" */ `@/views/${view.join('/')}.vue`);
// }

const routes = [
  {
    path: '/',
    name: 'home',
    component: () => import(/* webpackChunkName: "home" */ '../views/Home.vue'),
  },
  {
    path: '/about',
    name: 'about',
    component: () => import(/* webpackChunkName: "about" */ '../views/About.vue'),
    meta: {
      title: 'About',
    },
  },
  {
    path: '/stop/:stop',
    name: 'stop',
    component: () => import(/* webpackChunkName: "stop" */ '../views/Stop.vue'),
  },
  {
    path: '/trip/:trip/:vehicle',
    name: 'trip',
    component: () => import(/* webpackChunkName: "trip" */ '../views/Trip.vue'),
  },
  {
    path: '/map',
    name: 'map',
    meta: { disableFooter: true },
    component: () => import(/* webpackChunkName: "map" */ '../views/Map.vue'),
  },
  {
    path: '/map/vehicle/:vehicle',
    name: 'mapVehicle',
    meta: { disableFooter: true },
    component: () => import(/* webpackChunkName: "map" */ '../views/Map.vue'),
  },
  {
    path: '/map/stop/:stop',
    name: 'mapStop',
    meta: { disableFooter: true },
    component: () => import(/* webpackChunkName: "map" */ '../views/Map.vue'),
  },
  {
    path: '*',
    name: 'notFound',
    component: () => import(/* webpackChunkName: "notFound" */ '../views/NotFound.vue'),
  },
];

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes,
});

router.beforeEach((to, from, next) => {
  Store.commit('setTitle', to.meta.title);
  next();
});

export default router;
