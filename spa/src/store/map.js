export default {
  namespaced: true,
  state: {
    view: null,
  },
  mutations: {
    setView(state, view) {
      state.view = view;
    },
  },
  actions: {},
};
