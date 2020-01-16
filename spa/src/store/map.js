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
    load({ dispatch }) {
      dispatch('joinChannel', { name: 'geo:vehicles' }, { root: true });
      Api.emit('geo:stops');
    },
    unload({ commit, dispatch }, savedView) {
      dispatch('leaveChannel', 'geo:vehicles', { root: true });
      commit('setSavedView', savedView);
    },
  },
};
