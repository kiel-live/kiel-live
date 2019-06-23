# opnv-live

[![Build Status](https://travis-ci.org/Garogat/opnv-live.svg?branch=master)](https://travis-ci.org/Garogat/opnv-live)
[![](https://images.microbadger.com/badges/image/anbraten/opnv-live.svg)](https://microbadger.com/images/anbraten/opnv-live "Get your own image badge on microbadger.com")

This wep app allows you to view live updates of bus arrivals.

[![dockeri.co](https://dockeri.co/image/anbraten/opnv-live)](https://hub.docker.com/r/anbraten/opnv-live)

## Project setup
```
docker pull anbraten/opnv-live
```

### docker-compose.yml
```
version: '3'

services:
  transport:
    build: .
    environment:
      PORT: 8080
    restart: always
    ports:
      - 8080:8080
```

## Development

### Project setup
```
npm install
```

### Compiles and hot-reloads for development
```
npm run serve
```

### Compiles and minifies for production
```
npm run build
```

### Run your tests
```
npm run test
```

### Lints and fixes files
```
npm run lint
```

### Customize configuration
See [Configuration Reference](https://cli.vuejs.org/config/).


## Roadmap
- [ ] add ferry departures
- [ ] new icon
- [ ] ! add about and imprint (included from file)
- [ ] add openstreetmaps routes with estimated bus position
- [ ] add backend caching
- [ ] fetch live updates directly from client
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
