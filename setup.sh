#!/bin/sh

# Set these environment variables to use this script
# - NATS_URL
# - NATS_USER
# - NATS_PASSWORD

nats str rm --force data
nats str add data --config /stream.json
