#! /bin/sh

# check initial setup
if [[ ! -f /app/setup_done ]] || [[ -n "${NODE_ENV}" ]] || [[ "${NODE_ENV}" == "development" ]]; then
  cd /app/spa
  NODE_ENV= npm install --silent
  npm run build

  echo "" > /app/setup_done
fi

if [[ -n "${NODE_ENV}" ]] && [[ "${NODE_ENV}" == "development" ]]; then
  cd /app && npm run dev
else
  cd /app && npm run start
fi