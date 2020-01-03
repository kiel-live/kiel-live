const createClient = require('hafas-client');
const nahshProfile = require('hafas-client/p/nahsh');
const { post } = require('./cachedRequest');
const { join, leave, channels } = require('./autoUpdater');

const STOP_REFRESH_RATE = 10000;
const TRIP_REFRESH_RATE = 10000;
const STOP_DATA_URL = 'https://www.kvg-kiel.de/internetservice/services/passageInfo/stopPassages/stop';
const STOP_LOOKUP_URL = 'https://www.kvg-kiel.de/internetservice/services/lookup/autocomplete/json';
const TRIP_INFO_URL = 'https://www.kvg-kiel.de/internetservice/services/tripInfo/tripPassages';
// const STOP_URL = 'https://www.kvg-kiel.de/internetservice/services/stopInfo/stop';
// const STOP_ROUTE_URL = 'https://www.kvg-kiel.de/internetservice/services/routeInfo/route';

const hafas = createClient(nahshProfile, 'NAHSHPROD');
const joinedChannels = {};

/*
async function getStop(stopId) {
  const data = {
    stop: stopId,
    language: 'de',
  };

  return post(STOP_URL, data);
}
*/

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

  return res.filter((i) => i.id);
}

async function nearby({ longitude, latitude }) {
  const res = [];
  let tmp;

  try {
    tmp = await hafas.nearby({
      type: 'location',
      latitude,
      longitude,
    }, { distance: 400 });
  } catch (error) {
    console.log(error);
    return null;
  }

  for (let i = 0; i < tmp.length; i += 1) {
    const item = tmp[i];
    const name = item.name.replace(/Kiel\s/, '');
    const lookup = await lookupStops(name);

    if (lookup && lookup.length === 1) {
      res.push({
        ...item,
        ...lookup[0],
        gps: true,
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

function leaveChannels(clientId) {
  if (joinedChannels[clientId]) {
    leave({
      channel: joinedChannels[clientId],
      clientId,
    });
    joinedChannels[clientId] = null;
  }
}

function joinStop({ stopId, clientId }, cb) {
  const channel = `stop:${stopId}`;

  // leave old channels first
  leaveChannels(clientId);

  join({
    channel,
    clientId,
    timeout: STOP_REFRESH_RATE,
    load: async () => getStopData(stopId),
    cb,
  });
  joinedChannels[clientId] = channel;
}

function leaveStop({ clientId }) {
  leaveChannels(clientId);
}

function joinTrip({ tripId, vehicleId, clientId }, cb) {
  const channel = `trip:${tripId}:${vehicleId}`;

  // leave old channels first
  leaveChannels(clientId);

  join({
    channel,
    clientId,
    timeout: TRIP_REFRESH_RATE,
    load: async () => trip({ tripId, vehicleId }),
    cb,
  });
  joinedChannels[clientId] = channel;
}

function leaveTrip({ clientId }) {
  leaveChannels(clientId);
}

function status() {
  const c = channels();
  return {
    loadedChannels: c.length,
    channels: c,
  };
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
  leaveChannels,
};
