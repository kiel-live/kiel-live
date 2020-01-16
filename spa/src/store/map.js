import Api from '@/libs/api';

export default {
  namespaced: true,

  state: {
    savedView: null,
    vehicles: null,
    stops: null,
  },

  mutations: {
    setSavedView(state, savedView) {
      state.savedView = savedView;
    },
    setVehicles(state, vehicles) {
      state.vehicles = vehicles;
    },
    setStops(state, stops) {
      state.stops = stops;
    },
  },

  actions: {
    load() {
      Api.emit('geo:vehicles:join');
      Api.emit('geo:stops');
    },
    unload({ commit }, savedView) {
      Api.emit('geo:vehicles:leave');
      commit('map/setSavedView', savedView);
    },
  },
};
