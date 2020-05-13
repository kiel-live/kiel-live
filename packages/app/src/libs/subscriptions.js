import Api from '@/libs/api';

export default (store) => {
  // general socket listeners
  Api.on('connect', () => {
    store.commit('connect');
    store.dispatch('reJoin');
  });

  Api.on('disconnect', () => {
    store.commit('disconnect');
  });

  // map listeners
  Api.on('geo:vehicles', ({ vehicles }) => {
    store.commit('map/setVehicles', vehicles);
  });

  Api.on('geo:stops', ({ stops }) => {
    store.commit('map/setStops', stops);
  });

  // stop listeners
  Api.on('stop', (stop) => {
    store.commit('stop/setStop', stop);
  });

  Api.on('stop:search', (stops) => {
    store.commit('stop/setSearchStops', stops);
  });

  Api.on('stop:nearby', (stops) => {
    store.commit('stop/setGPSStops', stops);
  });

  // trip listeners
  Api.on('trip', (trip) => {
    store.commit('trip/setTrip', trip);
  });
};
