import Vue from 'vue';
import Vuex from 'vuex';
import { values, orderBy, uniqBy } from 'lodash';
import config from '@/libs/config';

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    isConnected: false,
    favoriteStops: values(JSON.parse(localStorage.getItem('favoriteStops')) || {}),
    title: null,
    stopQuery: '',
    gpsStops: [],
    stops: [],
    isTester: localStorage.getItem('tester') === 'true' || false,
  },
  getters: {
    allStops(state) {
      const stops = [
        ...state.favoriteStops,
        ...state.gpsStops,
        ...state.stops,
      ];

      return orderBy(uniqBy(stops, 'id'), ['favorite', 'gps'], ['desc', 'desc']);
    },
  },
  mutations: {
    connect(state) {
      state.isConnected = true;
    },
    disconnect(state) {
      state.isConnected = false;
    },
    addFavoriteStop(state, { id, name }) {
      state.favoriteStops.push({ id, name, favorite: true });
      localStorage.setItem('favoriteStops', JSON.stringify(state.favoriteStops));
    },
    removeFavoriteStop(state, id) {
      const stops = state.favoriteStops.filter((i) => i.id !== id);
      localStorage.setItem('favoriteStops', JSON.stringify(stops));
      state.favoriteStops = stops;
    },
    setTitle(state, title) {
      state.title = title;
      document.title = `${title ? `${title} - ` : ''}${config('title', 'OPNV Live')}`;
    },
    setStopQuery(state, query) {
      state.stopQuery = query;
    },
    setGPSStops(state, stops) {
      state.gpsStops = stops;
    },
    setStops(state, stops) {
      state.stops = stops;
    },
    setTester(state, isTester) {
      state.isTester = isTester;
      localStorage.setItem('tester', isTester);
    },
  },
  actions: {
  },
  modules: {
  },
});
