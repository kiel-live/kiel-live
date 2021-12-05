import 'windi.css';

import { createApp } from 'vue';
import { createPinia } from 'pinia';

import App from '~/App.vue';
import router from '~/router';
import { loadApi } from './api';

const app = createApp(App);

app.use(router);
app.use(createPinia());
app.mount('#app');

void loadApi();
