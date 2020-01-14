<template>
  <div class="map-container">
    <div id="map"></div>
  </div>
</template>

<script>
import L from 'leaflet';

import Api from '@/api';

export default {
  name: 'osmap2',
  data() {
    return {
      gpsSupport: true,
      ownLocation: null,
      vehicles: null,
      stops: null,
      osmap: null,
      vehicleLayer: null,
      stopLayer: null,
    };
  },
  computed: {
    visibleStops() {
      if (this.zoom < 15) {
        return []; // don't show stops
      }

      return [];
    },
  },
  props: {
    showVehicle: {
      type: String,
      default: null,
    },
    showStop: {
      type: String,
      default: null,
    },
  },
  methods: {
    loadMap() {
      this.osmap = L.map('map', {
        preferCanvas: true,
        zoom: 13,
        minZoom: 11, // zoom: 11-16
        maxZoom: 16,
        center: [54.321, 10.131],
        maxBounds: [
          [54.52, 9.9],
          [54.21, 10.44],
        ], // kiel city: left=9.9 bottom=54.21 right=10.44 top=54.52
      });

      if (process.env.NODE_ENV === 'development') {
        this.osmap.on('click', (ev) => {
          const latlng = this.osmap.mouseEventToLatLng(ev.originalEvent);
          console.log(latlng);
        });
      }

      this.osmap.on('zoomend', () => {
        // add stops if zoom is at least 14 else remove it
        if (this.stopLayer) {
          if (this.osmap.getZoom() < 13) {
            this.osmap.removeLayer(this.stopLayer);
          } else if (!this.osmap.hasLayer(this.stopLayer)) {
            this.osmap.addLayer(this.stopLayer);
          }
        }
      });

      L.tileLayer('/api/osm-tiles/{z}/{x}/{y}.png', {
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
      }).addTo(this.osmap);
    },
    updateVehicles({ vehicles }) {
      if (process.env.NODE_ENV === 'development') {
        // eslint-disable-next-line
        vehicles = [
          {
            id: 10,
            tripId: 'ttiu',
            name: '60s Ziegel',
            latitude: 54.327335997647666 * 3600000,
            longitude: 10.089225769042969 * 3600000,
          },
          {
            id: 11,
            tripId: 'tt',
            name: '1 HBF',
            latitude: 54.309513453509375 * 3600000,
            longitude: 10.088024139404299 * 3600000,
          },
        ];
      }

      const layer = L.layerGroup();
      this.vehicles = [];

      vehicles.forEach((v) => {
        if (!v.id || !v.name || !v.longitude || !v.latitude) {
          return;
        }

        const marker = L.circleMarker([v.latitude / 3600000, v.longitude / 3600000], {
          radius: 8,
          color: '#A00',
          fillColor: '#A00',
          fillOpacity: 0.5,
          stroke: true,
        }).addTo(layer);

        marker.on('click', () => {
          this.$router.push({ name: 'trip', params: { vehicle: v.id, trip: v.tripId } });
        });

        this.vehicles.push(marker);
      });

      if (this.vehicleLayer) {
        this.osmap.removeLayer(this.vehicleLayer);
      }
      this.vehicleLayer = layer;
      this.osmap.addLayer(this.vehicleLayer);
    },
    updateStops({ stops }) {
      const layer = L.layerGroup();
      this.stops = [];

      stops.forEach((s) => {
        const marker = L.circleMarker([s.latitude / 3600000, s.longitude / 3600000], {
          radius: 5,
          color: '#3388ff',
          fillColor: '#3388ff',
          fillOpacity: 1,
          stroke: false,
          data: s,
        }).addTo(layer);

        marker.on('click', () => {
          this.$router.push({ name: 'stop', params: { stop: s.shortName } });
        });

        this.stops.push(marker);
      });

      if (this.stopLayer) {
        this.osmap.removeLayer(this.stopLayer);
      }
      this.stopLayer = layer;
      this.osmap.addLayer(this.stopLayer);
    },
    centerGPS() {
      if (!this.gpsSupport || this.gpsLoading) { return; }

      this.gpsLoading = true;
      navigator.geolocation.getCurrentPosition((position) => {
        this.ownLocation = {
          lat: position.coords.latitude,
          lng: position.coords.longitude,
        };
        this.center = this.ownLocation;
      }, () => {
        this.gpsSupport = false;
      });
    },
    load() {
      this.join();
      Api.on('connect', this.join);
      // wait for vehicle updates
      Api.on('geo:vehicles', this.updateVehicles);
      Api.on('geo:stops', this.updateStops);
      this.loadMap();
    },
    unload() {
      Api.removeListener('connect', this.join);
      Api.removeListener('geo:vehicles', this.updateVehicles);
      Api.removeListener('geo:stops', this.updateStops);

      if (this.vehicles) {
        Api.emit('geo:vehicles:leave');
      }
    },
    reload() {
      this.unload();
      this.load();
    },
    join() {
      Api.emit('geo:vehicles:join');
      Api.emit('geo:stops');
    },
  },
  mounted() {
    this.gpsSupport = !!navigator.geolocation;
    this.load();

    // center requested vehicle / stop or gps location if available
    if (this.showVehicle || this.showStop) {
      this.zoom = 17;
    } else if (this.gpsSupport) {
      this.centerGPS();
    }
  },
  beforeDestroy() {
    this.unload();
  },
  beforeRouteUpdate(to, from, next) {
    next();
    this.reload();
  },
};
</script>

<style lang="scss" scoped>
  .map-container {
    position: relative;
    display: flex;
    width: 100%;
    flex-flow: column;
    flex-grow: 1;

    #map {
      width: 100%;
      height: 100%;
    }
  }
</style>

<style lang="scss">
  %vehiclemarker-common {
    font-size: 12px;
    color: white;
    padding: 2px;
    display: flex;
    background-image: url('/img/vehicle-icon.svg');
    background-size: 100% auto;
    background-repeat: no-repeat;
    // transition: transform 1s linear;
  }

  %vehiclemarker-common-text{
    display: block;
    text-align: center;
    width: 66%;
    margin: auto 0px auto 0px;
  }

  .vehiclemarker {
    @extend %vehiclemarker-common;

    span {
      @extend %vehiclemarker-common-text;
    }
  }

  .vehiclemarker-rotated {
    @extend %vehiclemarker-common;
    span {
      @extend %vehiclemarker-common-text;
      transform: scale(-1, -1);
      transform-origin: 50% 50% 50%;
      //text-align: right;
    }
  }
</style>
