<template>
  <div class="relative h-full w-full items-center justify-center overflow-hidden">
    <Map :geojson="geojson" :selected-marker="selectedMarker" @marker-click="selectedMarker = $event" />
    <DetailsPopup :is-open="!!selectedMarker" @close="selectedMarker = undefined">
      <MarkerPopup v-if="selectedMarker" :marker="selectedMarker" />
    </DetailsPopup>
    <Appbar />
  </div>
</template>

<script lang="ts">
import { GeoJSONSourceRaw } from 'maplibre-gl';
import { computed, defineComponent, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { stops, subscribe, vehicles } from '~/api';
import Appbar from '~/components/AppBar.vue';
import DetailsPopup from '~/components/DetailsPopup.vue';
import Map from '~/components/Map.vue';
import MarkerPopup from '~/components/popups/MarkerPopup.vue';
import { Marker } from '~/types';

export default defineComponent({
  // eslint-disable-next-line vue/multi-word-component-names
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
      async set(marker) {
        if (!marker) {
          await router.replace({ name: 'home' });
          return;
        }
        await router.replace({ name: 'map-marker', params: { markerType: marker.type, markerId: marker.id } });
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
          coordinates: [v.location.longitude / 3600000, v.location.latitude / 3600000],
        },
      })),
    );

    const stopsGeoJson = computed(() =>
      Object.values(stops.value).map((s) => ({
        type: 'Feature',
        properties: { type: 'stop', name: s.name, id: s.id },
        geometry: {
          type: 'Point',
          coordinates: [s.location.longitude / 3600000, s.location.latitude / 3600000],
        },
      })),
    );

    const geojson = computed<GeoJSONSourceRaw['data']>(() => ({
      type: 'FeatureCollection',
      features: [...vehiclesGeoJson.value, ...stopsGeoJson.value],
    }));

    onMounted(async () => {
      await subscribe('data.map.vehicle.>', vehicles);
      await subscribe('data.map.stop.>', stops);
    });

    return { geojson, selectedMarker };
  },
});
</script>
