<template>
  <div class="relative h-full w-full items-center justify-center overflow-hidden">
    <Map :geojson="geojson" @marker-click="selectedMarker = $event" />
    <DetailsPopup :is-open="!!selectedMarker" @close="selectedMarker = undefined">
      <MarkerPopup v-if="selectedMarker" :marker="selectedMarker" />
    </DetailsPopup>
    <Appbar />
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, onMounted, ref } from 'vue';
import { GeoJSONSourceRaw } from 'maplibre-gl';
import 'maplibre-gl/dist/maplibre-gl.css';

import { vehicles, stops, loadApi } from '~/api';
import Map from '~/components/Map.vue';
import DetailsPopup from '~/components/DetailsPopup.vue';
import Appbar from '~/components/Appbar.vue';
import { useRoute, useRouter } from 'vue-router';
import { Marker } from '~/types';
import MarkerPopup from '~/components/popups/MarkerPopup.vue';

export default defineComponent({
  name: 'Home',

  components: { Map, DetailsPopup, Appbar, MarkerPopup },

  setup() {
    const route = useRoute();
    const router = useRouter();
    const selectedMarker = computed<Marker | undefined>({
      get() {
        if (route.name !== 'map-marker') {
          return undefined;
        }

        return {
          type: route.params.markerType,
          id: route.params.markerId,
        } as Marker;
      },
      set(marker) {
        if (!marker) {
          router.replace({ name: 'home' });
          return;
        }
        router.replace({ name: 'map-marker', params: { markerType: marker.type, markerId: marker.id } });
      },
    });

    const vehiclesGeoJson = computed(() =>
      Object.values(vehicles.value).map((v) => ({
        type: 'Feature',
        properties: {
          type: 'vehicle',
          name: v.name,
          id: v.id,
          number: v.name.split(' ')[0],
          to: v.name.split(' ').slice(1).join(' '),
          iconName: `busIcon-unfocused-${v.name.split(' ')[0]}-${v.location.heading}`,
          iconNameFocused: `busIcon-focused-${v.name.split(' ')[0]}-${v.location.heading}`,
        },

        geometry: {
          type: 'Point',
          coordinates: [v.location.longitude, v.location.latitude],
        },
      })),
    );

    const stopsGeoJson = computed(() =>
      Object.values(stops.value).map((s) => ({
        type: 'Feature',
        properties: { type: 'stop', name: s.name, id: s.id },
        geometry: {
          type: 'Point',
          coordinates: [s.location.longitude, s.location.latitude],
        },
      })),
    );

    const geojson = computed<GeoJSONSourceRaw['data']>(() => ({
      type: 'FeatureCollection',
      features: [...vehiclesGeoJson.value, ...stopsGeoJson.value],
    }));

    onMounted(async () => {
      await loadApi();
    });

    return { geojson, selectedMarker };
  },
});
</script>
