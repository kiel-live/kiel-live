# kiel-live

This wep app allows you to view live updates of bus arrivals.

[![Docker Build](https://github.com/kiel-live/kiel-live/workflows/Docker%20Build/badge.svg)](https://github.com/kiel-live/kiel-live/actions?query=workflow%3A%22Docker+Build%22)
[![Linting](https://github.com/kiel-live/kiel-live/workflows/Linting/badge.svg)](https://github.com/kiel-live/kiel-live/actions?query=workflow%3ALinting)
[![Docker layers](https://images.microbadger.com/badges/image/anbraten/opnv-live.svg)](https://microbadger.com/images/anbraten/opnv-live)

[![Docker stats](https://dockeri.co/image/anbraten/opnv-live)](https://hub.docker.com/r/anbraten/opnv-live)

## Features

* Show realtime information (route, direction, eta) of bus arrvials on a specific stop
* Add stops to favorites for everyday use
* View trip stops of currently driving busses
* Show nearby stops and distance by using gps location
* View stops and realtime bus locations on a map

## Screenshots

![Screenshot](screenshot.jpg)

## Install on server

```bash
docker pull anbraten/opnv-live
```

## Development

### Project setup

```bash
yarn
```

### Compiles and hot-reloads for development

```bash
yarn dev
```

### Compiles and minifies for production

```bash
yarn build
```

### Lints and fixes files

```bash
yarn lint

yarn lint:fix
```

## Roadmap

See [Roadmap](https://github.com/anbraten/opnv-live/projects/1)

## API

The backend is accessable via socket.io-websocket (`/api/socket/`).

### Datatypes

#### ID clashing

To prevent id clashing by data from multiple provides always use `<provider><id>` as the actual id.

#### Stop

A `stop` is a fixed point their for example a bus stop or a car-sharing parking spot is located.

[Stop type](packages/commons/src/stop.class.ts)

#### Stop arrival

A `stop-arrival` is a list of vehicles arriving and departing at a stop.

[Stop arrival type](packages/commons/src/stop-arrival.class.ts)

#### Vehicle

A `vehicle` can be of a specific type (exp. bus, bike).

[Vehicle type](packages/commons/src/vehicle.class.ts)

#### Trip

A `trip` is a tour represented by a list of `stops` executed by a `vehicle` (exp. bus) on a specific `route`.

[Trip type](packages/commons/src/trip.class.ts)

#### Route

A `route` is a fixed list of waypoints a vehicle could drive along. A one-time tour on a specific `route` is called a `trip`.

[Route type](packages/commons/src/route.class.ts)

### Topics

* `stop`: List of all stop
* `stop-arrivals/:stopId` Arrivals of a single stop
* `vehicle`: List of all vehicles
* `trip/:id` Unsubscribe from a specific trip
* `route/:id` List of all routes
