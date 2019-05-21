const axios = require('axios');
const qs = require('querystring');

const REFRESH_RATE = 5000;
const STOP_URL = 'https://www.kvg-kiel.de/internetservice/services/stopInfo/stop';
const STOP_DATA_URL = 'https://www.kvg-kiel.de/internetservice/services/passageInfo/stopPassages/stop';
const STOP_LOOKUP_URL = 'https://www.kvg-kiel.de/internetservice/services/lookup/autocomplete';
const STOP_ROUTE_URL = 'https://www.kvg-kiel.de/internetservice/services/routeInfo/route';

const stops = {};
let io;

function setIO(_io) {
  io = _io;
}

async function post(url, data) {
  const options = {
    method: 'POST',
    headers: { 'content-type': 'application/x-www-form-urlencoded' },
    data: qs.stringify(data),
    url,
  };

  try {
    const repsonse = await axios(options);
    return repsonse.data;
  } catch (e) {
    console.log('HTTP-ERROR', url, data);
    return null;
  }
}

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

  return post(STOP_LOOKUP_URL, data);
}

function loop(stopId) {
  return async () => {
    const stop = await getStopData(stopId);
    stops[stopId].data = stop;
    io.to(`stop:${stopId}`).emit('stop', stop);
  };
}

function open(stopId) {
  // if already existing
  if (stops[stopId]) {
    return;
  }

  stops[stopId] = {
    connected: 0,
    loop: setInterval(loop(stopId), REFRESH_RATE),
  };

  loop(stopId)();
}

function close(stopId) {
  if (!stops[stopId]) {
    return;
  }

  clearInterval(stops[stopId].loop);

  delete stops[stopId];
}

function join(stopId, socket) {
  open(stopId);

  stops[stopId].connected += 1;

  // send last fetched data
  const { data } = stops[stopId];
  if (data) {
    socket.emit('stop', data);
  }
}

function leave(stopId) {
  if (stops[stopId]) {
    stops[stopId].connected -= 1;

    if (stops[stopId].connected < 1) {
      close(stopId);
    }
  }
}

module.exports = {
  setIO,
  join,
  leave,
  lookupStops,
};
