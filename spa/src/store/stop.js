import { values, orderBy, uniqBy } from 'lodash';
import Api from '@/libs/api';

export default {
  namespaced: true,

  state: {
    favoriteStops: values(JSON.parse(localStorage.getItem('favoriteStops')) || {}),

    // Home.vue
    GPSstops: [],
    searchStops: [],
    stopQuery: null,
    gpsLoading: false,
    gpsSupport: !!navigator.geolocation,

    // Stop.vue
    loadedStopId: null,
    stop: null,
  },

  getters: {
    allStops(state) {
      const stops = [
        ...state.favoriteStops,
        ...state.GPSstops,
        ...state.searchStops,
      ];

      return orderBy(uniqBy(stops, 'id'), ['favorite', 'gps'], ['desc', 'desc']);
    },
    isFavoriteStop(state) {
      if (!state.loadedStopId) {
        return false;
      }

      const favoriteStop = state.favoriteStops.filter((i) => i.id === state.loadedStopId);
      return favoriteStop.length === 1;
    },
  },

  mutations: {
    addFavoriteStop(state, { id, name }) {
      state.favoriteStops.push({ id, name, favorite: true });
      localStorage.setItem('favoriteStops', JSON.stringify(state.favoriteStops));
    },
    removeFavoriteStop(state, id) {
      const stops = state.favoriteStops.filter((i) => i.id !== id);
      localStorage.setItem('favoriteStops', JSON.stringify(stops));
      state.favoriteStops = stops;
    },

    // Home.vue
    setStopQuery(state, query) {
      state.stopQuery = query;
    },
    setGPSStops(state, stops) {
      state.gpsStops = stops;
      state.gpsLoading = false;
    },
    setSearchStops(state, searchStops) {
      state.searchStops = searchStops;
    },
    setGPSLoading(state, gpsLoading) {
      state.gpsLoading = gpsLoading;
    },
    setGPSSupport(state, gpsSupport) {
      state.gpsSupport = gpsSupport;
    },

    // Stop.vue
    setLoadedStopId(state, loadedStopId) {
      state.loadedStopId = loadedStopId;
    },
    setStop(state, stop) {
      state.stop = stop;
    },
  },

  actions: {
    // Home.vue
    searchStops({ commit }, query) {
      commit('setStopQuery', query);
      Api.emit('stop:search', query);
    },
    searchGPSStops({ commit, state }) {
      // prevent simultaneous calls
      if (!state.gpsSupport || state.gpsLoading) { return; }

      commit('setGPSLoading', true);
      navigator.geolocation.getCurrentPosition((position) => {
        Api.emit('stop:nearby', {
          latitude: position.coords.latitude,
          longitude: position.coords.longitude,
        });
      }, (error) => {
        // eslint-disable-next-line no-console
        console.error(error);
        commit('setGPSSupport', false);
      });
    },

    // Stop.vue
    load({ commit, dispatch }, stopId) {
      commit('setStop', null);
      dispatch('joinChannel', { name: 'stop', data: stopId }, { root: true });
      commit('setLoadedStopId', stopId);
    },
    unload({ state, commit, dispatch }) {
      if (state.loadedStopId) {
        dispatch('leaveChannel', 'stop', { root: true });
        commit('setLoadedStopId', null);
      }
    },
  },
};
