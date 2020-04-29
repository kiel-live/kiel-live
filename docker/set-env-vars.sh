#!/bin/sh
FILE_PATH=${1:-./spa/env-config.js}
FILE_DIR=$(dirname "${FILE_PATH}")
PREFIX=${2:-APP_}

# Recreate config file
rm -rf $FILE_PATH
mkdir -p $FILE_DIR && touch $FILE_PATH

# load env vars from file if one exists
if [ -f .env.local ]; then
  set -a
  . ./.env.local
  set +a
fi

# Add assignment
echo "window._env_ = {" >> $FILE_PATH

# Read each line in .env file
# Each line represents key=value pairs
printenv | while read line; do
  # skip env vars not starting with $PREFIX
  case "$line" in
    $PREFIX*)
      ;;
    *)
      continue
      ;;
  esac

  # Split env variables by character `=`
  if printf '%s\n' "$line" | grep -q -e '='; then
    varname=$(printf '%s\n' "$line" | sed -e 's/=.*//')
    value=$(printf '%s\n' "$line" | sed -e 's/^[^=]*=//')
  fi

  echo "Setting $varname"

  # Append configuration property to JS file
  echo "  $varname: '$value'," >> $FILE_PATH
done

echo "};" >> $FILE_PATH