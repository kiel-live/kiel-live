import Vue from 'vue';
import Vuex from 'vuex';
import config from '@/libs/config';

import Map from './map';
import Trips from './trips';
import Stops from './stops';

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    isConnected: false,
    title: null,
    isTester: localStorage.getItem('tester') === 'true' || false,
  },
  mutations: {
    connect(state) {
      state.isConnected = true;
    },
    disconnect(state) {
      state.isConnected = false;
    },
    setTitle(state, title) {
      state.title = title;
      document.title = `${title ? `${title} - ` : ''}${config('title', 'OPNV Live')}`;
    },
    setTester(state, isTester) {
      state.isTester = isTester;
      localStorage.setItem('tester', isTester);
    },
  },
  actions: {
  },
  modules: {
    map: Map,
    stops: Stops,
    trips: Trips,
  },
});
