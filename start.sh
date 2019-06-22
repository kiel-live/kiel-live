#! /bin/sh

echo "Starting opnv-live container ..."

# check initial setup
if [[ ! -f /app/setup_done ]] || [[ -n "${NODE_ENV}" ]] || [[ "${NODE_ENV}" == "development" ]]; then
  echo "Building spa ..."
  cd /app/spa
  export VUE_APP_VERSION=$(date)
  NODE_ENV= npm install --silent
  npm run build

  echo "" > /app/setup_done
  echo "Setup done!"
fi

if [[ -n "${NODE_ENV}" ]] && [[ "${NODE_ENV}" == "development" ]]; then
  cd /app && npm run dev
else
  cd /app && npm run start
fi