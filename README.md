# opnv-live

[![Build Status](https://travis-ci.org/anbraten/opnv-live.svg?branch=master)](https://travis-ci.org/anbraten/opnv-live)
[![](https://images.microbadger.com/badges/image/anbraten/opnv-live.svg)](https://microbadger.com/images/anbraten/opnv-live "Get your own image badge on microbadger.com")

This wep app allows you to view live updates of bus arrivals.

[![dockeri.co](https://dockeri.co/image/anbraten/opnv-live)](https://hub.docker.com/r/anbraten/opnv-live)

## Features
* Show realtime information (route, direction, eta) of bus arrvials on a specific stop
* Add stops to favorites for everyday use
* View trip stops of currently driving busses
* Show nearby stops and distance by using gps location
* View stops and realtime bus locations on a map

## Screenshots
![Screenshot](screenshot.jpg)

## Project setup
```
docker pull anbraten/opnv-live
```

## Development

### Project setup
```
yarn
```

### Compiles and hot-reloads for development
```
yarn dev
```

### Compiles and minifies for production
```
yarn build
```

### Lints and fixes files
```
yarn lint

yarn lint:fix
```

## Roadmap
See [Roadmap](https://github.com/anbraten/opnv-live/projects/1)
