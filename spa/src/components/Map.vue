<template>
  <div class="map">
    <l-map
      v-if="vehicles && stops"
      :zoom="zoom"
      :center="center"
      :maxBounds="maxBounds"
      :minZoom="minZoom"
      :maxZoom="maxZoom"
      @update:zoom="zoomUpdated"
      @update:center="centerUpdated"
      @update:bounds="boundsUpdated"
    >
      <l-tile-layer
        url="/api/osm-tiles/{z}/{x}/{y}.png"
        attribution="&copy; <a href='https://www.openstreetmap.org/copyright'>OpenStreetMap</a> contributors"
      />
      <l-rotated-marker
        v-for="v in vehicles"
        :key="v.id"
        :lat-lng="[v.latitude / 3600000, v.longitude / 3600000]"
        :rotationAngle="v.heading-90"
        :zIndexOffset="100"
        :name="v.name.split(' ')[0]"
        @click="clickVehicle({ vehicleId: v.id, tripId: v.tripId })"
      >
        <l-popup>{{ v.name }}</l-popup>
        <l-icon
          :className="v.heading > 180 ? 'vehiclemarker-rotated' : 'vehiclemarker'"
          :iconAnchor="[40 / 2, Math.round(((40 / 2) / 68) * 44)]"
          :iconSize="[40, Math.round((40 / 68) * 44)]"
          :popupAnchor="[-1, -14]"
          :shadowAnchor="[32, 32]"
          :shadowSize="[24, 24]"
        >
          <span>{{ v.name.split(' ')[0] }}</span>
        </l-icon>
      </l-rotated-marker>
      <template v-if="zoom > 14">
        <l-marker
          v-for="s in stops"
          :key="s.id"
          :lat-lng="[s.latitude / 3600000, s.longitude / 3600000]"
          :name="s.name"
          @click="clickStop(s.shortName)"
        >
          <l-popup>{{ s.name }}</l-popup>
        </l-marker>
      </template>
      <l-marker
        v-if="ownLocation"
        :lat-lng="ownLocation"
        :zIndexOffset="100"
        name="GPS"
      >
        <l-popup>GPS</l-popup>
      </l-marker>
    </l-map>
  </div>
</template>

<script>
import {
  LMap,
  LTileLayer,
  LIcon,
  LMarker,
  LPopup,
} from 'vue2-leaflet';
import LRotatedMarker from 'vue2-leaflet-rotatedmarker';

import Api from '@/api';

export default {
  name: 'osmap',
  components: {
    LMap,
    LTileLayer,
    LIcon,
    LPopup,
    LMarker,
    LRotatedMarker,
  },
  data() {
    return {
      gpsSupport: true,
      vehicles: null,
      stops: null,
      zoom: 14,
      center: [54.321, 10.131],
      maxBounds: [
        [54.52, 9.9],
        [54.21, 10.44],
      ], // kiel city: left=9.9 bottom=54.21 right=10.44 top=54.52
      ownLocation: null,
      minZoom: 11, // zoom: 11-16
      maxZoom: 16,
      bounds: null,
    };
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
    clickVehicle({ vehicleId, tripId }) {
      this.$router.push({ name: 'trip', params: { vehicle: vehicleId, trip: tripId } });
    },
    clickStop(stopId) {
      this.$router.push({ name: 'stop', params: { stop: stopId } });
    },
    zoomUpdated(zoom) {
      this.zoom = zoom;
    },
    centerUpdated(center) {
      this.center = center;
    },
    boundsUpdated(bounds) {
      this.bounds = bounds;
    },
    updateVehicles(data) {
      const vehicles = data.vehicles.filter((v) => v.id && v.name && v.longitude && v.latitude);

      // only list selected vehicle
      if (this.showVehicle) {
        const vehicle = vehicles.filter((v) => v.id === this.showVehicle)[0] || null;

        if (!vehicle) {
          return;
        }

        this.center = {
          lat: vehicle.latitude / 3600000,
          lng: vehicle.longitude / 3600000,
        };

        this.vehicles = [vehicle];
      } else {
        this.vehicles = vehicles;
      }
    },
    updateStops({ stops }) {
      if (this.showStop) {
        const stop = stops.filter((v) => v.shortName === this.showStop)[0] || null;

        if (!stop) {
          return;
        }

        this.center = {
          lat: stop.latitude / 3600000,
          lng: stop.longitude / 3600000,
        };

        this.stops = [stop];
      } else {
        this.stops = stops;
      }
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
  .map {
    position: relative;
    display: flex;
    width: 100%;
    flex-flow: column;
    flex-grow: 1;

    .map-container {
      position: absolute;
      width: 100%;
      height: 100%;
    }

    #osmap {
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
    transition: transform 1s linear;
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
