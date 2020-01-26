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
yarn install
```

### Compiles and hot-reloads for development
```
yarn run serve
```

### Compiles and minifies for production
```
yarn run build
```

### Run your tests
```
yarn run test
```

### Lints and fixes files
```
yarn run lint
```

### Customize configuration
See [Configuration Reference](https://cli.vuejs.org/config/).


## Roadmap
- [ ] add ferry departures
- [ ] ! add about and imprint (link to external page?)
- [x] add stops to map
- [x] add map link: home => map city overview
- [x] add map link: bus trip (map icon) <=> live bus position on map
- [x] add map link: stop (map icon) <=> stop on map
- [x] add map with vehicle live locations
- [ ] map smooth vehicle locations (https://github.com/LouisMazel/vue2-leaflet-movingmarker)
- [ ] add routes tab to stop view (and current view to depatures tab)
- [ ] add location image to stop view (none clickable map)
- [ ] add translations
- [ ] ~~fetch live updates directly from client~~ (api doesn't provide correct CORS header)
- [x] add backend caching
- [x] new icon
- [x] new title
- [x] fancier back button on stop page
- [x] "connection to server lost" screen
- [x] sort orders based on status
- [x] add favorite stops
- [x] show alerts
- [x] add route details on bus click
- [x] ! show nearby stops based on html5 geolocation
- [x] add spa tracking
- [x] store last stop queries
