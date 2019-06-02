import Vue from 'vue';
import Vuex from 'vuex';

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    isConnected: false,
    favoriteStops: JSON.parse(localStorage.getItem('favoriteStops')) || {},
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
  },
  actions: {

  },
});
