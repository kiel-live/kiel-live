import Vue from 'vue';
import Router from 'vue-router';
import Store from '@/store';

Vue.use(Router);

// route level code-splitting
// this generates a separate chunk (name.[hash].js) for this route
// which is lazy-loaded when the route is visited.
function loadView(...view) {
  return () => import(/* webpackChunkName: "view-[request]" */ `@/views/${view.join('/')}.vue`);
}

const router = new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      name: 'home',
      component: loadView('Home'),
    },
    {
      path: '/about',
      name: 'about',
      component: loadView('About'),
      meta: {
        title: 'About',
      },
    },
    {
      path: '/stop/:stop',
      name: 'stop',
      component: loadView('Stop'),
    },
    {
      path: '*',
      name: 'notFound',
      component: loadView('NotFound'),
    },
  ],
});

router.beforeEach((to, from, next) => {
  Store.commit('setTitle', to.meta.title);
  next();
});

export default router;
