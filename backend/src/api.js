const createClient = require('hafas-client');
const nahshProfile = require('hafas-client/p/nahsh');
const { post } = require('./cachedRequest');
const { join, leave, channels } = require('./autoUpdater');

const STOP_REFRESH_RATE = 10000;
const TRIP_REFRESH_RATE = 10000;
const GEO_VEHICLES_REFRESH_RATE = 1000;

const BASE_URL = 'https://www.kvg-kiel.de';
const STOP_DATA_URL = `${BASE_URL}/internetservice/services/passageInfo/stopPassages/stop`;
const STOP_LOOKUP_URL = `${BASE_URL}/internetservice/services/lookup/autocomplete/json`;
const TRIP_INFO_URL = `${BASE_URL}/internetservice/services/tripInfo/tripPassages`;
// const STOP_URL = `${BASE_URL}/internetservice/services/stopInfo/stop`;
// const STOP_ROUTE_URL = `${BASE_URL}/internetservice/services/routeInfo/route`;
const GEO_VEHICLES_URL = `${BASE_URL}/internetservice/geoserviceDispatcher/services/vehicleinfo/vehicles`;
const GEO_STOPS_URL = `${BASE_URL}/internetservice/geoserviceDispatcher/services/stopinfo/stops`;

const hafas = createClient(nahshProfile, 'NAHSHPROD');
const joinedChannels = {};

function unEscapeHtml(unsafe) {
  return unsafe
    .replace('&amp;', /&/g)
    .replace('&lt;', /</g)
    .replace('&gt;', />/g)
    .replace('&quot;', /"/g)
    .replace('&#039;', /'/g)
    .replace('&auml;', 'ä')
    .replace('&Auml;', 'Ä')
    .replace('&ouml;', 'ö')
    .replace('&Ouml;', 'Ö')
    .replace('&uuml;', 'ü')
    .replace('&Uuml;', 'Ü')
    .replace('&szlig;', 'ß');
}

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

  res.forEach((stop) => {
    stop.name = unEscapeHtml(stop.name);
  });

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

async function geoVehiclesLocation() {
  const data = {
    cacheBuster: new Date().getTime(),
    colorType: 'ROUTE_BASED',
    // lastUpdate: new Date().getTime(),
    positionType: 'RAW', // 'CORRECTED',
  };

  return post(GEO_VEHICLES_URL, data);
}

function joinGeoVehicles({ clientId }, cb) {
  const channel = 'geo:vehicles';

  // leave old channels first
  leaveChannels(clientId);

  join({
    channel,
    clientId,
    timeout: GEO_VEHICLES_REFRESH_RATE,
    load: async () => geoVehiclesLocation(),
    cb,
  });
  joinedChannels[clientId] = channel;
}

async function geoStops() {
  const data = {
    top: 324000000,
    bottom: -324000000,
    left: -648000000,
    right: 648000000,
  };

  return post(GEO_STOPS_URL, data);
}

function leaveGeoVehicles({ clientId }) {
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
  joinGeoVehicles,
  leaveGeoVehicles,
  geoStops,
};
