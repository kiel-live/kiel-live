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
    connected: 0,
    loop: setInterval(loop(channel), timeout),
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

function join({ channel, load, cb, timeout }) {
  open({ channel, load, cb, timeout });

  looper[channel].connected += 1;
 
  // send last fetched data
  const { data } = looper[channel];
  if (data) {
    cb(data);
  }
}

function leave(channel) {
  if (looper[channel]) {
    looper[channel].connected -= 1;

    if (looper[channel].connected < 1) {
      close(channel);
    }
  }
}

module.exports = {
  join,
  leave,
};
