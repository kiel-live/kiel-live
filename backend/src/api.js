const createClient = require('hafas-client');
const nahshProfile = require('hafas-client/p/nahsh');
const { post } = require('./cachedRequest');
const { join, leave, channels } = require('./autoUpdater');

const STOP_REFRESH_RATE = 10000;
const TRIP_REFRESH_RATE = 10000;
const STOP_URL = 'https://www.kvg-kiel.de/internetservice/services/stopInfo/stop';
const STOP_DATA_URL = 'https://www.kvg-kiel.de/internetservice/services/passageInfo/stopPassages/stop';
const STOP_LOOKUP_URL = 'https://www.kvg-kiel.de/internetservice/services/lookup/autocomplete/json';
const STOP_ROUTE_URL = 'https://www.kvg-kiel.de/internetservice/services/routeInfo/route';
const TRIP_INFO_URL = 'https://www.kvg-kiel.de/internetservice/services/tripInfo/tripPassages';

const hafas = createClient(nahshProfile, 'NAHSHPROD');

async function getStop(stopId) {
  const data = {
    stop: stopId,
    language: 'de',
  };

  return post(STOP_URL, data);
}

async function getStopData(stopId) {
  const data = {
    cacheBuster: new Date().getTime(),
    mode: 'departure',
    language: 'de',
    stop: stopId,
  };

  return post(STOP_DATA_URL, data);
}

async function lookupStops(query) {
  const data = {
    query,
    language: 'de',
  };

  const res = await post(STOP_LOOKUP_URL, data);

  return res.filter(i => i.id);
}

async function nearby({ longitude, latitude }) {
  const res = [];
  let tmp;

  try {
    tmp = await hafas.nearby({
      type: 'location',
      latitude,
      longitude,
    }, {distance: 400});
  } catch (error) {
    console.log(error);
    return;
  }

  for (let i = 0; i < tmp.length; i++) {
    const item = tmp[i];
    const name = item.name.replace(/Kiel\s/, '');
    const lookup = await lookupStops(name);

    if (lookup && lookup.length === 1) {
      res.push({
        ...item,
        ...lookup[0],
        'gps': true,
      });
    }
  }

  return res;
}

async function trip({ tripId, vehicleId }) {
  const data = {
    cacheBuster: new Date().getTime(),
    tripId,
    vehicleId,
    mode: 'departure',
    language: 'de',
  };

  return post(TRIP_INFO_URL, data);
}

function joinStop(id, cb) {
  join({
    channel: `stop:${id}`,
    timeout: STOP_REFRESH_RATE,
    load: async () => getStopData(id),
    cb,
  });
}

function leaveStop(id) {
  leave(`stop:${id}`);
}

function joinTrip({ tripId, vehicleId }, cb) {
  join({
    channel: `trip:${tripId}:${vehicleId}`,
    timeout: TRIP_REFRESH_RATE,
    load: async () => trip({ tripId, vehicleId }),
    cb,
  });
}

function leaveTrip({ tripId, vehicleId }) {
  leave(`trip:${tripId}:${vehicleId}`);
}

function status() {
  return {
    channels: channels(),
  }
}

module.exports = {
  joinStop,
  leaveStop,
  joinTrip,
  leaveTrip,
  nearby,
  trip,
  lookupStops,
  status,
};
