import startPubServer from './pub-server';
import startSubServer from './sub-server';

const SUB_PORT = process.env.BACKEND_SUB_PORT ? parseInt(process.env.BACKEND_SUB_PORT) : 3000;
const PUB_PORT = process.env.BACKEND_PUB_PORT ? parseInt(process.env.BACKEND_PUB_PORT) : 3030;

async function start() {
  await startPubServer(PUB_PORT);
  await startSubServer(SUB_PORT);

  console.log('Server started.');
}

void start();
