import Vue from 'vue';
import Vuex from 'vuex';
import { orderBy } from 'lodash';

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    isConnected: false,
    favoriteStops: JSON.parse(localStorage.getItem('favoriteStops')) || {},
    title: null,
    stopQuery: '',
    gpsStops: {},
    stops: {},
  },
  getters: {
    allStops(state) {
      return orderBy({
        ...state.stops,
        ...state.gpsStops,
        ...state.favoriteStops,
      }, 'favorite', 'desc');
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
      Vue.set(state.favoriteStops, id, { id, name, favorite: true });
      localStorage.setItem('favoriteStops', JSON.stringify(state.favoriteStops));
    },
    removeFavoriteStop(state, id) {
      if (state.favoriteStops[id]) {
        Vue.delete(state.favoriteStops, id);
        localStorage.setItem('favoriteStops', JSON.stringify(state.favoriteStops));
      }
    },
    setTitle(state, title) {
      state.title = title;
      document.title = `${title ? `${title} - ` : ''}${process.env.VUE_APP_TITLE || 'OPNV Live'}`;
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
  },
  actions: {

  },
});
