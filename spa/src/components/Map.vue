<template>
  <div class="map-container">
    <div id="map"></div>
  </div>
</template>

<script>
import L from 'leaflet';
import 'leaflet.locatecontrol';
import '@/libs/LVectorMarker';

import Api from '@/api';

export default {
  name: 'osmap',
  data() {
    return {
      vehicles: {},
      stops: null,
      osmap: null,
      vehicleLayer: null,
      stopLayer: null,
    };
  },
  props: {
    focusVehicle: {
      type: String,
      default: null,
    },
    focusStop: {
      type: String,
      default: null,
    },
  },
  methods: {
    loadMap() {
      const CustomCanvas = L.Canvas.extend({
        _updateCustomPath(layer) {
          if (!this._drawing || layer._empty()) { return; }

          const ctx = this._ctx;
          layer._customDraw(ctx);
        },
      });

      this.osmap = L.map('map', {
        preferCanvas: true,
        renderer: new CustomCanvas(),
        minZoom: 12,
        // maxZoom: 18,
        zoomControl: false,
        maxBounds: [
          [54.52, 9.9],
          [54.21, 10.44],
        ], // kiel city: left=9.9 bottom=54.21 right=10.44 top=54.52
      });

      this.osmap.on('zoomend', () => {
        this.updateLayer();
      });

      // const tileUrl = '/api/osm-tiles/{z}/{x}/{y}.png';
      const tileUrl = 'http://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png';
      L.tileLayer(tileUrl, {
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
      }).addTo(this.osmap);

      L.control.zoom({
        position: 'bottomright',
      }).addTo(this.osmap);

      L.control.locate({
        position: 'bottomright',
      }).addTo(this.osmap);

      const savedView = this.$store.state.map.view || null;
      if (!this.focusVehicle && !this.focusStop) {
        if (savedView) {
          this.osmap.setView(savedView.center, savedView.zoom); // center last location
        } else {
          this.osmap.setView([54.321, 10.131], 13); // center kiel city
        }
      }

      this.vehicleLayer = L.layerGroup();
      this.vehicleLayer.addTo(this.osmap);
    },
    updateLayer() {
      // add stops if zoom is at least 14, else remove it
      if (this.stopLayer) {
        if (this.osmap.getZoom() < 14) {
          this.osmap.removeLayer(this.stopLayer);
        } else if (!this.osmap.hasLayer(this.stopLayer)) {
          this.osmap.addLayer(this.stopLayer);

          // re add vehicles to keep them in front
          this.osmap.removeLayer(this.vehicleLayer);
          this.osmap.addLayer(this.vehicleLayer);
        }
      }
    },
    updateVehicles({ vehicles }) {
      const vehicleUpdates = [];

      vehicles.forEach((v) => {
        if (!v.id || !v.name || !v.longitude || !v.latitude) {
          return;
        }

        vehicleUpdates.push(v.id);

        // vehicle already exists
        if (this.vehicles[v.id]) {
          // this.vehicles[v.id].slideTo([v.latitude / 3600000, v.longitude / 3600000], { duration: 5000, /* keepAtCenter: true, */ });
          this.vehicles[v.id].setLatLng([v.latitude / 3600000, v.longitude / 3600000]);
        } else {
          const marker = L.vehicleMarker([v.latitude / 3600000, v.longitude / 3600000], {
            id: v.id,
            radius: 12,
            color: '#A00',
            fillOpacity: 1,
            label: v.name.split(' ')[0],
          }).addTo(this.vehicleLayer);

          marker.on('click', () => {
            this.osmap.setView([v.latitude / 3600000, v.longitude / 3600000], 17);
            if (this.focusVehicle !== v.id) {
              this.$router.push({ name: 'mapTrip', params: { vehicle: v.id, trip: v.tripId } });
            }
          });

          this.vehicles[v.id] = marker;
        }

        if (this.focusVehicle === v.id) {
          this.osmap.setView([v.latitude / 3600000, v.longitude / 3600000], 17);
        }
      });

      // remove not updated vehicles
      Object.keys(this.vehicles).forEach((vid) => {
        if (!vehicleUpdates.includes(vid)) {
          this.vehicles[vid].remove();
          delete this.vehicles[vid];
        }
      });
    },
    updateStops({ stops }) {
      const layer = L.layerGroup();
      this.stops = [];

      stops.forEach((s) => {
        const marker = L.stopMarker([s.latitude / 3600000, s.longitude / 3600000], {
          radius: 7,
          color: '#3388ff',
          fillColor: '#3388ff',
          fillOpacity: 1,
          stroke: false,
          data: s,
        }).addTo(layer);

        marker.on('click', () => {
          this.osmap.setView([s.latitude / 3600000, s.longitude / 3600000], 17);
          if (this.focusStop !== s.shortName) {
            this.$router.push({ name: 'mapStop', params: { stop: s.shortName } });
          }
        });

        this.stops.push(marker);

        if (this.focusStop === s.shortName) {
          this.osmap.setView([s.latitude / 3600000, s.longitude / 3600000], 17);
        }
      });

      if (this.stopLayer) {
        this.osmap.removeLayer(this.stopLayer);
      }
      this.stopLayer = layer;
      this.updateLayer();
    },
    join() {
      Api.emit('geo:vehicles:join');
      Api.emit('geo:stops');
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
  },
  mounted() {
    this.load();

    // center requested vehicle / stop or gps location if available
    if (this.showVehicle || this.showStop) {
      this.zoom = 17;
    }
  },
  beforeDestroy() {
    this.unload();
    this.$store.commit('map/setView', {
      center: this.osmap.getCenter(),
      zoom: this.osmap.getZoom(),
    });
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

  .leaflet-tile {
    filter: grayscale(1);
  }
</style>
