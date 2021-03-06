<template>
  <div class="map-container">
    <a class="back button" @click="$router.go(-1)"><i class="fas fa-angle-double-left" /></a>
    <template v-if="mapStyle">
      <MglMap
        id="map"
        :map-style="mapStyle"
        :center.sync="center"
        :min-zoom="minZoom"
        :max-zoom="maxZoom"
        :zoom.sync="zoom"
        :max-bounds="maxBounds"
        :attribution-control="false"
        @click="onClickMap"
        @load="onMapLoaded"
      >
        <MglAttributionControl position="top-right" />
        <MglNavigationControl position="bottom-right" />
        <MglGeolocateControl position="bottom-right" />
        <MglGeojsonLayer
          source-id="stops"
          :source="{ type: 'geojson', data: stopsGeoJson }"
          layer-id="stops"
          :layer="stopsLayer"
          @click="onClickStop"
          @mouseenter="onMouseEnter"
          @mouseleave="onMouseLeave"
        />
        <MglGeojsonLayer
          source-id="vehicles"
          :source="{ type: 'geojson', data: vehiclesGeoJson }"
          layer-id="vehicles"
          :layer="vehiclesLayer"
          @click="onClickVehicle"
          @mouseenter="onMouseEnter"
          @mouseleave="onMouseLeave"
        />
      </MglMap>
      <div v-if="focusData" class="focus-popup">
        <a v-if="focusStop" class="body" @click="$router.push({ name: 'stop', params: { stop: focusData && focusData.shortName } })">
          <i class="fas fa-sign" />
          <span>{{ focusData && focusData.name }}</span>
          <i class="fas fa-external-link-alt" />
        </a>
        <a v-if="focusVehicle" class="body" @click="$router.push({ name: 'trip', params: { vehicle: focusVehicle, trip: focusData && focusData.tripId } })">
          <span class="route"><i v-if="focusData && focusData.category === 'bus'" class="icon fas fa-bus" />{{ focusData && focusData.name.split(' ')[0] }}</span>
          <span class="direction">{{ focusData && focusData.name.split(' ').slice(1).join(' ') }}</span>
          <i class="fas fa-external-link-alt" />
        </a>
        <a class="close button" @click="$router.replace({ name: 'map' })"><i class="fas fa-times" /></a>
      </div>
    </template>
    <p v-else>
      Die Map ist nicht konfiguriert ;-D
    </p>
  </div>
</template>

<script>
import Mapbox from 'mapbox-gl';
import {
  MglMap,
  MglAttributionControl,
  MglNavigationControl,
  MglGeolocateControl,
  MglGeojsonLayer,
} from 'vue-mapbox';
import { mapState } from 'vuex';
import BusIcon from '@/libs/busIcon';
import config from '@/libs/config';

export default {
  components: {
    MglMap,
    MglAttributionControl,
    MglNavigationControl,
    MglGeolocateControl,
    MglGeojsonLayer,
  },
  props: {
    focusStop: {
      type: String,
      default: null,
    },
    focusVehicle: {
      type: String,
      default: null,
    },
  },
  data() {
    return {
      mapStyle: config('tile_server_url'),
      minZoom: 11,
      maxZoom: 18,
      // [west, south, east, north]
      maxBounds: [9.8, 54.21, 10.44, 54.52],
      center: null,
      zoom: null,
      needToFocus: false,
    };
  },
  computed: {
    ...mapState({
      vehicles: (state) => state.map.vehicles,
      stops: (state) => state.map.stops,
      savedView: (state) => state.map.savedView,
    }),
    stopsGeoJson() {
      if (!this.stops) { return null; }
      return {
        type: 'FeatureCollection',
        features: this.stops.map((stop) => ({
          type: 'Feature',
          geometry: {
            type: 'Point',
            coordinates: [this.convertLatLng(stop.longitude), this.convertLatLng(stop.latitude)],
          },
          properties: {
            id: stop.id,
          },
        })),
      };
    },
    vehiclesGeoJson() {
      if (!this.vehicles) { return null; }
      return {
        type: 'FeatureCollection',
        features: this.vehicles.filter((v) => v.id && v.latitude && v.longitude && v.name).map((vehicle) => ({
          type: 'Feature',
          geometry: {
            type: 'Point',
            coordinates: [this.convertLatLng(vehicle.longitude), this.convertLatLng(vehicle.latitude)],
          },
          properties: {
            id: vehicle.id,
            number: vehicle.name.split(' ')[0],
            to: vehicle.name.split(' ').slice(1).join(' '),
            iconName: `busIcon-unfocused-${vehicle.name.split(' ')[0]}-${vehicle.heading}`,
            iconNameFocused: `busIcon-focused-${vehicle.name.split(' ')[0]}-${vehicle.heading}`,
          },
        })),
      };
    },
    stopsLayer() {
      return {
        id: 'stops',
        type: 'circle',
        source: 'stops',
        paint: {
          'circle-color': [
            'match',
            ['get', 'id'],
            this.focusStop || '',
            '#1673fc',
            '#4f96fc',
          ],
          'circle-radius': [
            'match',
            ['get', 'id'],
            this.focusStop || '',
            8,
            5,
          ],
          'circle-stroke-opacity': 0,
          'circle-stroke-width': 5,
          'circle-opacity': this.focusVehicle ? 0.5 : 1,
        },
      };
    },
    vehiclesLayer() {
      return {
        id: 'vehicles',
        type: 'symbol',
        source: 'vehicles',
        paint: {
          'icon-opacity': [
            'match',
            ['get', 'number'],
            (this.focusData && this.focusData.name.split(' ')[0]) || '',
            1,
            this.focusVehicle ? 0.3 : 1,
          ],
        },
        layout: {
          'icon-image': [
            'match',
            ['get', 'id'],
            this.focusVehicle || '',
            ['get', 'iconNameFocused'],
            ['get', 'iconName'],
          ],
          'icon-rotation-alignment': 'map',
          'icon-allow-overlap': true,
          'symbol-sort-key': [
            'match',
            ['get', 'number'],
            (this.focusData && this.focusData.name.split(' ')[0]) || '',
            2,
            1,
          ],
        },
      };
    },
    focusData() {
      return (this.focusVehicle && this.vehicles && this.vehicles.find((v) => v.id === this.focusVehicle))
          || (this.focusStop && this.stops && this.stops.find((s) => s.id === this.focusStop));
    },
  },
  created() {
    // We need to set mapbox-gl library here in order to use it in template
    this.mapbox = Mapbox;
    this.map = null;
  },
  mounted() {
    this.$store.dispatch('map/load');
    if (this.focusVehicle || this.focusStop) {
      this.needToFocus = true;
    }
    if (this.savedView) {
      this.center = this.savedView.center;
      this.zoom = this.savedView.zoom;
    } else {
      this.center = [10.1283, 54.3166];
      this.zoom = 14;
    }
  },
  beforeDestroy() {
    const view = {
      center: this.center,
      zoom: this.zoom,
    };
    this.$store.dispatch('map/unload', view);
  },
  methods: {
    convertLatLng(value) {
      return value / 3600000;
    },
    onClickMap() {
      if (!this.focusData) { return; }
      this.$router.replace({ name: 'map' });
    },
    onClickStop(e) {
      if (this.focusStop === e.mapboxEvent.features[0].properties.id) { return; } // prevent reloading of same stop
      this.$router.replace({ name: 'mapStop', params: { stop: e.mapboxEvent.features[0].properties.id } });
    },
    onClickVehicle(e) {
      if (this.focusVehicle === e.mapboxEvent.features[0].properties.id) { return; } // prevent reloading of same stop
      this.$router.replace({ name: 'mapVehicle', params: { vehicle: e.mapboxEvent.features[0].properties.id } });
    },
    onMouseEnter(e) {
      e.map.getCanvas().style.cursor = 'pointer';
    },
    onMouseLeave(e) {
      e.map.getCanvas().style.cursor = '';
    },
    onMapLoaded(event) {
      this.map = event.map;
      this.map.on('styleimagemissing', (e) => {
        const [, focus, route, heading] = e.id.split('-');
        this.map.addImage(e.id, new BusIcon(this.map, focus === 'focused', route, heading), { pixelRatio: 2 });
      });
      if (this.needToFocus && this.focusData) {
        this.map.flyTo({
          center: [this.convertLatLng(this.focusData.longitude), this.convertLatLng(this.focusData.latitude)],
          zoom: 14,
        });
        this.needToFocus = false;
      }
    },
  },
};
</script>

<style lang="scss">
  @import '~mapbox-gl/dist/mapbox-gl';
</style>

<style lang="scss" scoped>
  .map-container {
    position: relative;
    display: flex;
    width: 100%;
    flex-flow: column;
    flex-grow: 1;
    border-bottom: 1px solid #b5b5b5;
    overflow: hidden;
    justify-content: center;

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

    ::v-deep .mapboxgl-ctrl-bottom-right {
      bottom: 3rem;
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
      transform: translate(-50%, 0);
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
        padding: calc(.375em - 1px) .5em;
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
