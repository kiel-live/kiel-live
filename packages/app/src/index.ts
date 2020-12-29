import { topic } from './client';

function start() {
  const z = topic('stop').sub((stop) => {
    console.log(stop);

    if (stop.alerts) {
      z.stop();
    }
  });
}
