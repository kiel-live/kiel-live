# kiel-live

This app allows you to view Kiels public transport (busses, bus-stops) in realtime.

[![Docker Build](https://github.com/kiel-live/kiel-live/workflows/Docker%20Build/badge.svg)](https://github.com/kiel-live/kiel-live/actions?query=workflow%3A%22Docker+Build%22)
[![Linting](https://github.com/kiel-live/kiel-live/workflows/Linting/badge.svg)](https://github.com/kiel-live/kiel-live/actions?query=workflow%3ALinting)
[![Docker layers](https://images.microbadger.com/badges/image/anbraten/opnv-live.svg)](https://microbadger.com/images/anbraten/opnv-live)

[![Docker stats](https://dockeri.co/image/anbraten/opnv-live)](https://hub.docker.com/r/anbraten/opnv-live)

## Features

* View realtime bus positions and bus-stops on map
* Get realtime information (route, direction, eta) of bus arrivals of a specific bus-stop
* Add bus-stops to favorites
* View trip stops of a bus

## Screenshots

![Screenshot](screenshot.jpg)

# Development

## Structure

The project contains following parts:

- `app/`: A PWA written with Vue3
- `android-app/`: A native android app wrapper of the PWA
- `collectors/*`: Multiple agents to scrape data from different apis
- `nats/`: The NATS server used as message broker to stream data from collectors to the PWA clients

## Gitpod

[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/kiel-live/kiel-live)

## PWA development

Copy the `.env.sample` file to `.env`. For the PWA development you only need to set `VITE_NATS_URL`.
You can set it to `wss://api.kiel-live.ju60.de/` to use the production server so you don't need to start your own backend (nats & collectors).

```bash
cd app/
pnpm install # install dependencies
pnpm start # start the PWA
```

## Nats & collectors development

Nats is the message broker used to bring data from the collectors to the PWA clients.

To start Nats simply copy the `.env.sample` file to `.env` adjust as needed and run `docker-compose up -d`.
