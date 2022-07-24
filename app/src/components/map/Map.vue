<template>
  <div id="map" class="w-full h-full" />
</template>

<script lang="ts">
// eslint-disable-next-line no-restricted-imports
import 'maplibre-gl/dist/maplibre-gl.css';

import type { Feature, FeatureCollection, Point } from 'geojson';
import {
  AttributionControl,
  CircleLayerSpecification,
  GeoJSONSource,
  GeolocateControl,
  LineLayerSpecification,
  Map,
  Source,
  SymbolLayerSpecification,
} from 'maplibre-gl';
import { computed, defineComponent, onMounted, PropType, Ref, toRef, watch } from 'vue';

import { stops, subscribe, trips, vehicles } from '~/api';
import { Marker } from '~/api/types';
import BusIcon from '~/components/map/busIcon';
import { usePrefersColorSchemeDark } from '~/compositions/usePrefersColorScheme';
import { brightMapStyle, darkMapStyle } from '~/config';

export default defineComponent({
  name: 'Map',

  props: {
    selectedMarker: {
      type: Object as PropType<Marker>,
      default: () => ({}),
    },
  },

  emits: {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    markerClick: (_marker?: Marker) => true,
  },

  setup(props, { emit }) {
    let map: Map;
    let initial = true;

    const prefersColorSchemeDark = usePrefersColorSchemeDark();

    const vehiclesGeoJson = computed<Feature[]>(() =>
      Object.values(vehicles.value).map((v) => ({
        type: 'Feature',
        properties: {
          type: v.type,
          name: v.name,
          id: v.id,
          number: v.name.split(' ')[0],
          to: v.name.split(' ').slice(1).join(' '),
          iconName: `${v.type}-unfocused-${v.name.split(' ')[0]}-${v.location.heading}`,
          iconNameFocused: `${v.type}-focused-${v.name.split(' ')[0]}-${v.location.heading}`,
        },

        geometry: {
          type: 'Point',
          coordinates: [v.location.longitude / 3600000, v.location.latitude / 3600000],
        },
      })),
    );

    const stopsGeoJson = computed<Feature[]>(() =>
      Object.values(stops.value).map((s) => ({
        type: 'Feature',
        properties: { type: s.type, name: s.name, id: s.id },
        geometry: {
          type: 'Point',
          coordinates: [s.location.longitude / 3600000, s.location.latitude / 3600000],
        },
      })),
    );

    const selectedMarker = toRef(props, 'selectedMarker');

    const selectedVehicle = computed(() => vehicles.value[selectedMarker.value.id]);

    const trip = computed(() => {
      if (!trips.value || !selectedVehicle.value) {
        return null;
      }
      return trips.value[selectedVehicle.value.tripId];
    });

    const tripsGeoJson = computed<Feature[]>(() => {
      if (selectedMarker.value.type === 'bus' && trip.value?.arrivals) {
        return [
          {
            type: 'Feature',
            properties: {
              type: 'trip',
            },
            geometry: {
              type: 'LineString',
              coordinates: trip.value?.arrivals
                ?.map((arrival) => {
                  if (stops.value[arrival.id]) {
                    return [
                      stops.value[arrival.id].location.longitude / 3600000,
                      stops.value[arrival.id].location.latitude / 3600000,
                    ];
                  }
                  return [0, 0];
                })
                .filter((c) => c[0] !== 0 && c[1] !== 0),
            },
          },
        ];
      }
      return [];
    });

    const geojson = computed<FeatureCollection>(() => ({
      type: 'FeatureCollection',
      features: [...vehiclesGeoJson.value, ...stopsGeoJson.value, ...tripsGeoJson.value],
    }));

    const stopsLayer: Ref<CircleLayerSpecification> = computed(() => ({
      id: 'stops',
      type: 'circle',
      source: 'geojson',
      filter: ['==', 'type', 'bus-stop'],
      paint: {
        'circle-color': ['match', ['get', 'id'], selectedMarker.value.id || '', '#1673fc', '#4f96fc'],
        'circle-radius': ['match', ['get', 'id'], selectedMarker.value.id || '', 8, 5],
        'circle-opacity': selectedMarker.value.type === 'bus' ? 0.5 : 1,
        'circle-stroke-opacity': 0,
        'circle-stroke-width': 5,
      },
    }));

    const vehiclesLayer: Ref<SymbolLayerSpecification> = computed(() => ({
      id: 'vehicles',
      type: 'symbol',
      source: 'geojson',
      paint: {
        'icon-opacity': [
          'match',
          ['get', 'number'],
          vehicles.value[selectedMarker.value.id]?.name.split(' ')[0] ?? '',
          1,
          selectedMarker.value.type === 'bus' ? 0.3 : 1,
        ],
      },
      filter: ['==', 'type', 'bus'],
      layout: {
        'icon-image': [
          'match',
          ['get', 'id'],
          selectedMarker.value.id || '',
          ['get', 'iconNameFocused'],
          ['get', 'iconName'],
        ],
        'icon-rotation-alignment': 'map',
        'icon-allow-overlap': true,
        'symbol-sort-key': [
          'match',
          ['get', 'number'],
          vehicles.value[selectedMarker.value.id]?.name.split(' ')[0] ?? '',
          2,
          1,
        ],
      },
    }));

    const tripsLayer: Ref<LineLayerSpecification> = computed(() => ({
      id: 'trips',
      type: 'line',
      source: 'geojson',
      filter: ['==', 'type', 'trip'],
      paint: {
        'line-width': 3,
        'line-color': 'rgb(170, 0, 0)',
      },
    }));

    function flyTo(center: [number, number]) {
      if (!map) {
        return;
      }

      map.flyTo({
        center,
        padding: {
          bottom: 500, // TODO use 3/4 of screen height
        },
      });
    }

    onMounted(async () => {
      await subscribe('data.map.vehicle.>', vehicles);
      await subscribe('data.map.stop.>', stops);

      map = new Map({
        container: 'map',
        // style: 'https://demotiles.maplibre.org/style.json',
        style: prefersColorSchemeDark.value ? darkMapStyle : brightMapStyle,
        minZoom: 11,
        maxZoom: 18,
        center: [10.1283, 54.3166],
        zoom: 14,
        // [west, south, east, north]
        maxBounds: [9.8, 54.21, 10.44, 54.52],
        attributionControl: false,
      });

      const attributionControl = new AttributionControl({ compact: true });
      map.addControl(attributionControl, 'bottom-left');

      map.addControl(
        new GeolocateControl({
          positionOptions: {
            enableHighAccuracy: true,
          },
          trackUserLocation: true,
        }),
        'bottom-right',
      );

      // map.addControl(new NavigationControl({}), 'bottom-right');

      map.on('styleimagemissing', (e) => {
        const [type, focus, route, heading] = e.id.split('-');
        if (type === 'bus') {
          map.addImage(e.id, new BusIcon(map, focus === 'focused', route, Number.parseInt(heading, 10)), {
            pixelRatio: 2,
          });
        }
      });

      map.on('load', () => {
        map.addSource('geojson', {
          type: 'geojson',
          data: Object.freeze(geojson.value),
        });

        map.addLayer(stopsLayer.value);

        map.addLayer(tripsLayer.value);

        map.addLayer(vehiclesLayer.value);

        initial = false;
      });

      // Change the cursor to a pointer when the it enters a feature in the 'symbols' layer.
      map.on('mouseenter', 'vehicles', () => {
        map.getCanvas().style.cursor = 'pointer';
      });

      // Change it back to a pointer when it leaves.
      map.on('mouseleave', 'vehicles', () => {
        map.getCanvas().style.cursor = '';
      });

      // Change the cursor to a pointer when the it enters a feature in the 'symbols' layer.
      map.on('mouseenter', 'stops', () => {
        map.getCanvas().style.cursor = 'pointer';
      });

      // Change it back to a pointer when it leaves.
      map.on('mouseleave', 'stops', () => {
        map.getCanvas().style.cursor = '';
      });

      map.on('click', (e) => {
        const features = map.queryRenderedFeatures(e.point, {
          layers: ['stops', 'vehicles'],
        });

        // Deselect marker when the map is clicked
        if (features.length === 0) {
          emit('markerClick');
          return;
        }

        const feature = features[0] as unknown as {
          geometry: Point;
          properties: Marker;
        };

        // Prevent reloading the same marker
        if (feature.properties.id === selectedMarker.value.id) {
          return;
        }

        flyTo(feature.geometry.coordinates as [number, number]);
        emit('markerClick', { type: feature.properties.type, id: feature.properties.id });
      });

      map.on('drag', () => {
        if (selectedMarker.value !== null) {
          emit('markerClick');
        }
      });
    });

    watch(prefersColorSchemeDark, () => {
      if (prefersColorSchemeDark.value) {
        map.setStyle(darkMapStyle);
      } else {
        map.setStyle(brightMapStyle);
      }
      // TODO: properly re-render custom layers
      // eslint-disable-next-line no-restricted-globals
      location.reload();
    });

    watch(geojson, () => {
      if (!map) {
        return;
      }

      const geoJSONSource = map.getSource('geojson');
      const isGeoJsonSource = (source?: Source): source is GeoJSONSource => source?.type === 'geojson';
      if (isGeoJsonSource(geoJSONSource)) {
        geoJSONSource.setData(Object.freeze(geojson.value));
      }
    });

    watch(stopsLayer, () => {
      if (!map || initial) {
        return;
      }

      if (stopsLayer.value.layout) {
        Object.keys(stopsLayer.value.layout).forEach((key) => {
          if (stopsLayer.value.layout) {
            map.setLayoutProperty('stops', key, stopsLayer.value.layout[key as keyof typeof stopsLayer.value.layout]);
          }
        });
      }

      if (stopsLayer.value.paint) {
        Object.keys(stopsLayer.value.paint).forEach((key) => {
          if (stopsLayer.value.paint) {
            map.setPaintProperty('stops', key, stopsLayer.value.paint[key as keyof typeof stopsLayer.value.paint]);
          }
        });
      }
    });

    watch(vehiclesLayer, () => {
      if (!map || initial) {
        return;
      }

      if (vehiclesLayer.value.layout) {
        Object.keys(vehiclesLayer.value.layout).forEach((key) => {
          if (vehiclesLayer.value.layout) {
            map.setLayoutProperty(
              'vehicles',
              key,
              vehiclesLayer.value.layout[key as keyof typeof vehiclesLayer.value.layout],
            );
          }
        });
      }

      if (vehiclesLayer.value.paint) {
        Object.keys(vehiclesLayer.value.paint).forEach((key) => {
          if (vehiclesLayer.value.paint) {
            map.setPaintProperty(
              'vehicles',
              key,
              vehiclesLayer.value.paint[key as keyof typeof vehiclesLayer.value.paint],
            );
          }
        });
      }
    });

    const selectedMarkerItem = computed(() => {
      const marker = selectedMarker.value;
      if (!marker) {
        return undefined;
      }
      return geojson.value.features.find((f) => (f.properties as Marker).id === marker.id);
    });
    watch(selectedMarkerItem, (_selectedMarkerItem) => {
      if (!map || !_selectedMarkerItem) {
        return;
      }

      flyTo((_selectedMarkerItem.geometry as Point)?.coordinates as [number, number]);
    });
  },
});
</script>

<style scoped>
#map :deep(.maplibregl-ctrl-attrib) {
  box-sizing: content-box;
}
</style>
