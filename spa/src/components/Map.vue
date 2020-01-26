<template>
  <div class="map-container">
    <a @click="$router.go(-1)" class="back button"><i class="fas fa-angle-double-left"/></a>
    <div id="map"></div>
    <div v-if="focusData" class="focus-popup">
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
      <a class="close button" @click="$router.replace({ name: 'map' })"><i class="fas fa-times" /></a>
    </div>
  </div>
</template>

<script>
import L from 'leaflet';
import 'leaflet.locatecontrol';
import '@/libs/LVectorMarker';
import { mapState } from 'vuex';

export default {
  name: 'osmap',
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
  data() {
    return {
      osmap: null,
      vehicleLayer: null,
      stopLayer: null,
      markers: {},
    };
  },
  computed: {
    ...mapState({
      vehicles: (state) => state.map.vehicles,
      stops: (state) => state.map.stops,
      savedView: (state) => state.map.savedView,
    }),
    focusMarker() {
      return this.markers[this.focusStop || this.focusVehicle] || null;
    },
    focusData() {
      return (this.focusMarker && this.focusMarker.options) || null;
    },
    focused() {
      return !!this.focusMarker;
    },
  },
  watch: {
    focusVehicle(vehicle) {
      if (vehicle) {
        // join vehicle
      }
    },
    focusStop(stop) {
      if (stop) {
        // join stop
      }
    },
    focusMarker(marker, old) {
      console.log('marker updated');

      // unselect old marker
      if (old) {
        old.options.focused = false;
      }

      // focus new marker
      if (marker) {
        marker.options.focused = true;

        if (this.osmap) {
          this.osmap.setView(marker.getLatLng(), 17);
        }
      }
    },
    focusData() {
      // TODO: fix location updates
      console.log('marker data updated');
    },
    focused(focused) {
      if (focused) {
        console.log('disable controls');
      } else {
        console.log('enable controls');
      }
      // TODO: re-enable this.setMapControlsEnabled(!focused);
    },
    vehicles(vehicles) {
      const updated = {}; // updated vehicle-ids list

      vehicles.forEach((v) => {
        if (!v.id || !v.latitude || !v.longitude) { return; } // skip incomplete records
        const marker = this.markers[v.id] || null;
        if (marker) {
          marker.setLatLng([v.latitude / 3600000, v.longitude / 3600000]);
          const options = {
            ...marker.options,
            label: v.name.split(' ').shift(),
          };
          this.$set(this.markers[v.id], 'options', options);
          this.$set(this.markers[v.id], 't', Date.now());
        } else {
          this.addVehicle(v);
        }
        updated[v.id] = true;
      });

      // remove none updated stop-markers
      this.cleanUpMarkers('vehicle', updated);
    },
    stops(stops) {
      const updated = {}; // updated stop-ids list

      stops.forEach((s) => {
        const marker = this.markers[s.id] || null;
        if (marker) {
          marker.setLatLng([s.latitude / 3600000, s.longitude / 3600000]);
          const options = {
            ...marker.options,
            ...s,
          };
          this.$set(this.markers[s.id], 'options', options);
        } else {
          this.addStop(s);
        }
        updated[s.id] = true;
      });

      this.cleanUpMarkers('stop', updated);
    },
  },
  methods: {
    initMap() {
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
        maxZoom: 18,
        zoomControl: false,
        maxBounds: [
          [54.52, 9.8],
          [54.21, 10.44],
        ], // kiel city: left=9.9 bottom=54.21 right=10.44 top=54.52
      });

      // leave possibile vehicle / stop focus to center gps location
      this.osmap.on('onlocationfound', () => {
        this.leaveFocus();
        // TODO: necessary?
      });

      // const tileUrl = '/api/osm-tiles/{z}/{x}/{y}.png';
      const tileUrl = 'https://maps.targomo.com/styles/gray-gl-style/rendered/{z}/{x}/{y}.png';
      L.tileLayer(tileUrl, {
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
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
      if (this.$route.name === 'map') {
        if (this.savedView) {
          this.osmap.setView(this.savedView.center, this.savedView.zoom); // center last location
        } else {
          this.osmap.setView([54.321, 10.131], 13); // center kiel city
        }
      }

      // add layer for vehicle markers
      this.vehicleLayer = L.layerGroup();
      this.vehicleLayer.addTo(this.osmap);

      // add layer for stop markers
      this.stopLayer = L.layerGroup();
      this.stopLayer.addTo(this.osmap);
    },
    addVehicle(v) {
      const marker = L.vehicleMarker([v.latitude / 3600000, v.longitude / 3600000], {
        ...v, // unpack vehicle data
        type: 'vehicle',
        label: v.name && v.name.split(' ')[0],
        focused: v.id === this.focusVehicle,
      });

      // focus vehicle
      marker.on('click', () => {
        if (this.focusVehicle === v.id) { return; } // prevent reloading of same vehicle
        this.$router.replace({ name: 'mapVehicle', params: { vehicle: v.id } });
      });

      marker.addTo(this.vehicleLayer);
      this.$set(this.markers, v.id, marker);
    },
    addStop(s) {
      const marker = L.stopMarker([s.latitude / 3600000, s.longitude / 3600000], {
        ...s, // unpack stop data
        type: 'stop',
        focused: s.id === this.focusStop,
      });

      // focus stop
      marker.on('click', () => {
        if (this.focusStop === s.id) { return; } // prevent reloading of same stop
        this.$router.replace({ name: 'mapStop', params: { stop: s.id } });
      });

      marker.addTo(this.stopLayer);
      this.$set(this.markers, s.id, marker);
    },
    cleanUpMarkers(type, ids) {
      // remove none updated markers
      Object.keys(this.markers).forEach((id) => {
        const marker = this.markers[id];
        if (marker.options.type === type && !ids[marker.options.id]) {
          console.log('deleted', marker.options.type);
          marker.remove(); // remove marker from map
          this.$delete(this.markers, id); // notify and delete marker from list via vue
        }
      });
    },
    setMapControlsEnabled(enabled) {
      const fnc = enabled ? 'enable' : 'disable';
      this.osmap.dragging[fnc]();
      this.osmap.touchZoom[fnc]();
      this.osmap.doubleClickZoom[fnc]();
      this.osmap.scrollWheelZoom[fnc]();
      this.osmap.boxZoom[fnc]();
      this.osmap.keyboard[fnc]();
      if (this.osmap.tap) this.osmap.tap[fnc]();

      if (enabled) {
        document.getElementById('map').style.cursor = 'grab';
      } else {
        document.getElementById('map').style.cursor = 'default';
      }
    },
    load() {
      this.initMap();
      this.$store.dispatch('map/load');
    },
    unload() {
      let view = null;
      if (this.osmap) {
        try {
          view = {
            center: this.osmap.getCenter(),
            zoom: this.osmap.getZoom(),
          };
        } catch (e) {
          // IGNORE
        }
      }
      this.$store.dispatch('map/unload', view);
    },
  },
  mounted() {
    this.load();
  },
  beforeDestroy() {
    this.unload();
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

    #map {
      position: absolute;
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
