import Vue from 'vue';
import Vuex from 'vuex';
import config from '@/libs/config';
import Api from '@/libs/api';
import subscribe from '@/libs/subscriptions';

import Map from './map';
import Trip from './trip';
import Stop from './stop';

Vue.use(Vuex);

const store = new Vuex.Store({
  modules: {
    map: Map,
    stop: Stop,
    trip: Trip,
  },

  state: {
    isConnected: false,
    joinedChannels: [],
    title: null,
    isTester: localStorage.getItem('tester') === 'true' || false,
  },

  mutations: {
    connect(state) {
      state.isConnected = true;
    },
    disconnect(state) {
      state.isConnected = false;
    },
    setTitle(state, title) {
      state.title = title;
      document.title = `${title ? `${title} - ` : ''}${config('title', 'OPNV Live')}`;
    },
    setTester(state, isTester) {
      state.isTester = isTester;
      localStorage.setItem('tester', isTester);
    },
    addJoinedChannel(state, { name, data = null }) {
      state.joinedChannels.push({ name, data });
    },
    removeJoinedChannel(state, channelName) {
      state.joinedChannels = state.joinedChannels.filter((c) => c.name !== channelName);
    },
  },

  actions: {
    joinChannel({ commit }, { name, data = null }) {
      commit('addJoinedChannel', { name, data });
      Api.emit(`join:${name}`, data);
    },
    leaveChannel({ commit, state }, channelName) {
      state.joinedChannels.forEach((c) => {
        if (c.name === channelName) {
          Api.emit(`leave:${c.name}`, c.data);
        }
      });
      commit('removeJoinedChannel', channelName);
    },
    reJoin({ state }) {
      state.joinedChannels.forEach((c) => {
        Api.emit(`join:${c.name}`, c.data);
      });
    },
  },
});

// register socket.io subscriptions
subscribe(store);

export default store;
