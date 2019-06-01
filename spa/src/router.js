import Vue from 'vue';
import Router from 'vue-router';

Vue.use(Router);

// route level code-splitting
// this generates a separate chunk (name.[hash].js) for this route
// which is lazy-loaded when the route is visited.
function loadView(...view) {
  return () => import(/* webpackChunkName: "view-[request]" */ `@/views/${view.join('/')}.vue`);
}

export default new Router({
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
