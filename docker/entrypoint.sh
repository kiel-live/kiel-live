#!/bin/sh
set -e

./set-env-vars.sh /app/dist/env-config.js

exec "$@"