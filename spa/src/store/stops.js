import { values, orderBy, uniqBy } from 'lodash';

export default {
  namespaced: true,
  state: {
    favoriteStops: values(JSON.parse(localStorage.getItem('favoriteStops')) || {}),
    gpsStops: [],
    stops: [],
    stopQuery: '',
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
    addFavoriteStop(state, { id, name }) {
      state.favoriteStops.push({ id, name, favorite: true });
      localStorage.setItem('favoriteStops', JSON.stringify(state.favoriteStops));
    },
    removeFavoriteStop(state, id) {
      const stops = state.favoriteStops.filter((i) => i.id !== id);
      localStorage.setItem('favoriteStops', JSON.stringify(stops));
      state.favoriteStops = stops;
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
};
