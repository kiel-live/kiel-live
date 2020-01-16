import Api from '@/libs/api';

export default {
  namespaced: true,

  state: {
    loadedTripId: null,
    trip: null,
  },

  mutations: {
    setLoadedStopId(state, loadedTripId) {
      state.loadedTripId = loadedTripId;
    },
    setTrip(state, trip) {
      state.trip = trip;
    },
  },

  actions: {
    load({ commit }, ids) {
      commit('setTrip', null);
      commit('setLoadedStopId', ids);
      Api.emit('trip:join', ids);
    },
    unload({ commit, state }) {
      if (state.loadedTripId) {
        Api.emit('trip:leave', state.loadedTripId);
        commit('setLoadedStopId', null);
      }
    },
  },
};
