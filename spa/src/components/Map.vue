<template>
  <div class="map-container">
    <a @click="$router.go(-1)" class="back button"><i class="fas fa-angle-double-left"/></a>
    <div class="map-overlay">
      <div id="map"></div>
    </div>
    <div v-if="(focusVehicle || focusStop) && focusData" class="focus-popup">
      <a v-if="focusVehicle" class="body" @click="$router.push({ name: 'trip', params: { vehicle: focusData.id, trip: focusData.tripId } })">
        <span class="route"><i v-if="focusData.category === 'bus'" class="icon fas fa-bus" />{{ focusData.name.split(' ')[0] }}</span>
        <span class="direction">{{ focusData.name.split(' ').slice(1).join(' ') }}</span>
        <i class="fas fa-external-link-alt"></i>
      </a>
      <a v-if="focusStop" class="body" @click="$router.push({ name: 'stop', params: { stop: focusData.shotName } })">
        <i class="fas fa-sign" />
        <span>{{ focusData.name }}</span>
        <i class="fas fa-external-link-alt"></i>
      </a>
      <a class="close button" @click="leaveFocus"><i class="fas fa-times" /></a>
    </div>
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
      focusData: null,
      isProgramaticViewUpdate: false,
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
          [54.52, 9.8],
          [54.21, 10.44],
        ], // kiel city: left=9.9 bottom=54.21 right=10.44 top=54.52
      });

      // leave possibile vehicle / stop focus if the user is trying to zoom / move the map himself
      /*
      this.osmap.on('zoomstart', (e) => {
        if (!this.isProgramaticViewUpdate) {
          console.log('user', e.type);
          this.leaveFocus();
        }
      });
      */

      // leave possibile vehicle / stop focus to center gps location
      this.osmap.on('onlocationfound', () => {
        this.leaveFocus();
      });

      // const tileUrl = '/api/osm-tiles/{z}/{x}/{y}.png';
      // const tileUrl = 'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png';
      const tileUrl = 'https://api.mapbox.com/styles/v1/mapbox/dark-v10/tiles/{z}/{x}/{y}?access_token={accessToken}';
      L.tileLayer(tileUrl, {
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
        accessToken: 'pk.eyJ1IjoiYW5icmF0ZW4iLCJhIjoiY2s1ZTg5bXJwMDI4eTNscnVkNmFldzM5biJ9.hrN_sy18PEbgu8QYYIYXiA',
      }).addTo(this.osmap);

      // zoom (+ & -) buttons
      L.control.zoom({
        position: 'bottomright',
      }).addTo(this.osmap);

      // gps locator button
      L.control.locate({
        position: 'bottomright',
        flyTo: true,
      }).addTo(this.osmap);

      // go to last visited location or center kiel
      const savedView = this.$store.state.map.view || null;
      if (!this.focusVehicle && !this.focusStop) {
        if (savedView) {
          this.setView(savedView.center, savedView.zoom); // center last location
        } else {
          this.setView([54.321, 10.131], 13); // center kiel city
        }
      }

      // add layer for vehicle markers
      this.vehicleLayer = L.layerGroup();
      this.vehicleLayer.addTo(this.osmap);

      // add layer for stop markers
      this.stopLayer = L.layerGroup();
      this.stopLayer.addTo(this.osmap);
    },
    setView(latlng, zoom) {
      if (!this.osmap) { return; }
      this.isProgramaticViewUpdate = true;
      this.osmap.setView(latlng, zoom);
      this.isProgramaticViewUpdate = false;
    },
    leaveFocus() {
      if (this.focusVehicle || this.focusStop) {
        if (this.focusStop) {
          this.stops.forEach((st) => {
            st.options.focused = false;
          });
        }
        this.focusData = null;
        if (this.$route.name !== 'map') {
          this.$router.replace({ name: 'map' });
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
          this.vehicles[v.id].options.heading = v.heading;
          this.vehicles[v.id].options.label = v.name.split(' ').shift();
          this.vehicles[v.id].options.focused = this.focusVehicle === v.id;
        } else {
          const marker = L.vehicleMarker([v.latitude / 3600000, v.longitude / 3600000], {
            id: v.id,
            heading: v.heading,
            label: v.name.split(' ')[0],
          }).addTo(this.vehicleLayer);

          // focus vehicle
          marker.on('click', (e) => {
            // prevent re-focus
            if (marker.options.focused) {
              return;
            }

            marker.options.focused = true;
            this.focusData = v;
            this.setView(e.latlng, 17);

            if (this.focusVehicle !== v.id) {
              this.$router.replace({ name: 'mapTrip', params: { vehicle: v.id, trip: v.tripId } });
            }
          });

          this.vehicles[v.id] = marker;
        }

        // re-center vehicle
        if (this.focusVehicle === v.id) {
          this.focusData = v;
          this.setView([v.latitude / 3600000, v.longitude / 3600000], 17);
        }
      });

      // remove not updated vehicles
      Object.keys(this.vehicles).forEach((vid) => {
        if (!vehicleUpdates.includes(vid)) {
          this.vehicles[vid].remove();
          if (this.focusData && this.focusData.id === vid) {
            this.leaveFocus();
          }
          delete this.vehicles[vid];
        }
      });
    },
    updateStops({ stops }) {
      this.stopLayer.clearLayers();
      this.stops = [];

      stops.forEach((s) => {
        const marker = L.stopMarker([s.latitude / 3600000, s.longitude / 3600000], {}).addTo(this.stopLayer);

        // focus stop
        marker.on('click', (e) => {
          // prevent re-focus
          if (marker.options.focused) {
            return;
          }

          this.stops.forEach((st) => {
            st.options.focused = false;
          });
          marker.options.focused = true;
          this.focusData = s;
          this.setView(e.latlng, 17);

          if (this.focusStop !== s.shortName) {
            this.$router.replace({ name: 'mapStop', params: { stop: s.shortName } });
          }
        });

        this.stops.push(marker);

        if (this.focusStop === s.shortName) {
          this.focusData = s;
          this.setView([s.latitude / 3600000, s.longitude / 3600000], 17);
        }
      });

      // hack to make sure vehicles are in front of the stops
      if (this.osmap.hasLayer(this.vehicleLayer)) {
        this.osmap.removeLayer(this.vehicleLayer);
        this.osmap.addLayer(this.vehicleLayer);
      }
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
  },
  beforeDestroy() {
    this.unload();
    this.$store.commit('map/setView', {
      center: this.osmap.getCenter(),
      zoom: this.osmap.getZoom(),
    });
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
    border-bottom: 1px solid #b5b5b5;

    .back {
      position: absolute;
      top: 1rem;
      left: 1rem;
      z-index: 500;
    }

    .map-overlay {
      width: 100%;
      height: 100%;
    }

    #map {
      width: 100%;
      height: 100%;
    }

    .focus-popup {
      position: absolute;
      display: flex;
      flex-direction: row;
      bottom: -1px;
      left: 50%;
      height: 3rem;
      width: 100%;
      margin: 0 auto;
      padding: 1rem .5rem;
      background: #fff;
      z-index: 1000;
      align-items: center;
      justify-content: space-between;
      transform: translate(-50%, 0%);
      border-bottom: 1px solid #b5b5b5;

      .body {
        display: flex;
        align-items: center;
        justify-content: space-between;
        cursor: pointer;
        width: 100%;

        .route {
          display: flex;
          font-size: 1rem;
          line-height: 1.2rem;
          height: 1rem;

          i {
            margin-right: .5rem;
          }
        }

        > * {
          margin: .5rem;
        }
      }

      span {
        line-height: 1rem;
      }

      .close {
        margin-left: 1rem;
        margin-right: 0;
        height: auto;
        padding: calc(.375em - 1px) 0.50em;
      }

      @media only screen and (min-width: 768px) {
        width: auto;
        border: 1px solid #b5b5b5;
        border-bottom: 0;
        border-top-left-radius: .5rem;
        border-top-right-radius: .5rem;
        max-width: 40rem;
        padding: 1.5rem;

        .close {
          margin-left: 3rem;
        }
      }
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

  .leaflet-bottom {
    bottom: 3rem;

    @media only screen and (min-width: 768px) {
      bottom: 0;
    }
  }
</style>
