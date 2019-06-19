#! /bin/sh

# check initial setup
if [ ! -f /app/setup_done ]; then
  cd /app/spa
  npm install --silent
  npm run build

  echo "" > /app/setup_done
fi

if [[ -n "${NODE_ENV}" ]] && [[ "${NODE_ENV}" == "development" ]]; then
  cd /app && npm run dev
  cd /app/spa && npm run build # rebuild on each start
else
  cd /app && npm run start
fi