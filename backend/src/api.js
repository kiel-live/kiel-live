const { getBoundingBox, headingDistanceTo } = require('geolocation-utils');
const { post } = require('./cachedRequest');
const { join, leave, channels } = require('./autoUpdater');

const STOP_REFRESH_RATE = 10000;
const TRIP_REFRESH_RATE = 10000;
const GEO_VEHICLES_REFRESH_RATE = 5000;

const BASE_URL = 'https://www.kvg-kiel.de';
const STOP_DATA_URL = `${BASE_URL}/internetservice/services/passageInfo/stopPassages/stop`;
const STOP_LOOKUP_URL = `${BASE_URL}/internetservice/services/lookup/autocomplete/json`;
const TRIP_INFO_URL = `${BASE_URL}/internetservice/services/tripInfo/tripPassages`;
// const STOP_URL = `${BASE_URL}/internetservice/services/stopInfo/stop`;
// const STOP_ROUTE_URL = `${BASE_URL}/internetservice/services/routeInfo/route`;
const GEO_VEHICLES_URL = `${BASE_URL}/internetservice/geoserviceDispatcher/services/vehicleinfo/vehicles`;
const GEO_STOPS_URL = `${BASE_URL}/internetservice/geoserviceDispatcher/services/stopinfo/stops`;

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

async function geoStops(_box) {
  const box = _box || {
    top: 324000000,
    bottom: -324000000,
    left: -648000000,
    right: 648000000,
  };

  return post(GEO_STOPS_URL, box);
}

async function nearby({ longitude, latitude }) {
  const radius = 200; // radius in meters
  const position = {
    lat: latitude,
    lon: longitude,
  };
  const box = getBoundingBox([position], radius);

  const { stops } = await geoStops({
    top: Math.round(box.topLeft.lat * 3600000),
    left: Math.round(box.topLeft.lon * 3600000),
    bottom: Math.round(box.bottomRight.lat * 3600000),
    right: Math.round(box.bottomRight.lon * 3600000),
  });

  return stops.map((s) => {
    const location = {
      lat: s.latitude / 3600000,
      lon: s.longitude / 3600000,
    };

    return {
      ...s,
      distance: Math.round(headingDistanceTo(position, location).distance),
    };
  });
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
