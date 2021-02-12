import Vue from 'vue';
import router from '@/router';
import store from '@/store';
import '@/libs/api';
import '@/libs/sentry';
import '@/registerServiceWorker';
import App from './App.vue';

Vue.config.productionTip = false;

new Vue({
  router,
  store,
  render: (h) => h(App),
}).$mount('#app');
