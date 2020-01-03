const looper = {};

function loop(channel) {
  return async () => {
    const data = await looper[channel].load();
    looper[channel].data = data;
    looper[channel].cb(data);
  };
}

function open({ channel, load, cb, timeout }) {
  // if already existing
  if (looper[channel]) {
    return;
  }

  looper[channel] = {
    connected: [],
    loop: setInterval(loop(channel), timeout),
    timeout,
    load,
    cb,
  };

  // first run
  loop(channel)();
}

function close(channel) {
  if (!looper[channel]) {
    return;
  }

  clearInterval(looper[channel].loop);

  delete looper[channel];
}

function join({ channel, load, cb, timeout, clientId }) {
  open({ channel, load, cb, timeout });

  looper[channel].connected.push(clientId);

  // send last fetched data
  const { data } = looper[channel];
  if (data) {
    cb(data);
  }
}

function leave({ channel, clientId }) {
  if (looper[channel]) {
    // remove client from list
    const { connected } = looper[channel];
    looper[channel].connected = connected.filter((c) => c !== clientId);

    if (looper[channel].connected.length < 1) {
      close(channel);
    }
  }
}

function channels() {
  return Object.keys(looper).map((channel) => ({
    name: channel,
    connected: looper[channel].connected.length,
    timeout: looper[channel].timeout,
  }));
}

module.exports = {
  join,
  leave,
  channels,
};
