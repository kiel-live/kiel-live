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
    load({ commit, dispatch }, ids) {
      commit('setTrip', null);
      dispatch('joinChannel', { name: 'trip', data: ids }, { root: true });
      commit('setLoadedStopId', ids);
    },
    unload({ commit, state, dispatch }) {
      if (state.loadedTripId) {
        dispatch('leaveChannel', 'trip', { root: true });
        commit('setLoadedStopId', null);
      }
    },
  },
};
