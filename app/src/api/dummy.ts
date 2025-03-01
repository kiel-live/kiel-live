import type { Ref } from 'vue';
import type { Api, Stop, Trip, Vehicle } from '~/api/types';
import { computed, ref } from 'vue';

const DUMMY_VEHICLES: Vehicle[] = [
  {
    id: 'bus-1',
    provider: 'dummy',
    name: 'Dummy Vehicle 1',
    type: 'bus',
    state: 'active',
    location: { latitude: 54.3239 * 3600000, longitude: 10.1228 * 3600000, heading: 0 },
    tripId: 'trip-1',
  },
  {
    id: 'bus-2',
    provider: 'dummy',
    name: 'Dummy Vehicle 2',
    type: 'bus',
    state: 'inactive',
    location: { latitude: 54.3237 * 3600000, longitude: 10.1229 * 3600000, heading: 0 },
    tripId: 'trip-2',
  },
  {
    id: 'bus-3',
    provider: 'kvg',
    name: '62 Russee, Schiefe Horn',
    type: 'bus',
    state: 'onfire',
    battery: '',
    location: {
      longitude: 36440127,
      latitude: 195621108,
      heading: 90,
    },
    tripId: 'trip-3',
    description: '',
  },
];

const DUMMY_STOPS: Stop[] = [
  {
    id: 'stop-1',
    provider: 'dummy',
    name: 'Dummy Stop 1',
    type: 'bus-stop',
    routes: ['1', '2'],
    alerts: ['Alert 1'],
    location: { latitude: 54.3233 * 3600000, longitude: 10.1228 * 3600000, heading: 0 },
    actions: [],
    arrivals: [
      {
        name: 'Dummy Vehicle 1',
        type: 'bus',
        vehicleId: 'bus-1',
        tripId: 'trip-1',
        routeId: '1',
        routeName: 'Route 1',
        direction: 'North',
        state: 'predicted',
        planned: '18:50',
        eta: 300,
        platform: '',
      },
    ],
  },
  {
    id: 'stop-2',
    provider: 'dummy',
    name: 'Dummy Stop 2',
    type: 'bus-stop',
    routes: ['3', '4'],
    alerts: ['Alert 2'],
    location: { latitude: 54.3234 * 3600000, longitude: 10.1229 * 3600000, heading: 0 },
    actions: [],
    arrivals: [
      {
        name: 'Dummy Vehicle 2',
        type: 'bus',
        vehicleId: 'bus-2',
        tripId: 'trip-2',
        routeId: '2',
        routeName: 'Route 2',
        direction: 'South',
        state: 'predicted',
        planned: '19:00',
        eta: 600,
        platform: '',
      },
    ],
  },
  {
    id: 'stop-3',
    provider: 'kvg',
    name: 'Lehmberg',
    type: 'bus-stop',
    routes: null,
    alerts: [],
    arrivals: [
      {
        name: 'Hassee',
        type: 'bus',
        vehicleId: 'bus-2',
        tripId: 'trip-2',
        routeId: '1610073983892324369',
        routeName: '51',
        direction: 'Hassee',
        state: 'predicted',
        planned: '19:23',
        eta: 145,
        platform: '',
      },
      {
        name: 'Bf. Melsdorf',
        type: 'bus',
        vehicleId: 'bus-3',
        tripId: 'trip-3',
        routeId: '1610073983892324377',
        routeName: '91',
        direction: 'Bf. Melsdorf',
        state: 'predicted',
        planned: '19:29',
        eta: 505,
        platform: '',
      },
    ],
    location: {
      longitude: 36460935,
      latitude: 195591727,
      heading: 0,
    },
  },
];

const DUMMY_TRIPS: Trip[] = [
  {
    id: 'trip-1',
    provider: 'dummy',
    direction: 'North',
    path: [{ latitude: 54.3233 * 3600000, longitude: 10.1228 * 3600000, heading: 0 }],
    arrivals: [
      {
        id: 'stop-1',
        name: 'Dummy Stop 1',
        state: 'predicted',
        planned: '18:50',
      },
    ],
  },
  {
    id: 'trip-2',
    provider: 'dummy',
    direction: 'South',
    path: [{ latitude: 54.3234 * 3600000, longitude: 10.1229 * 3600000, heading: 0 }],
    arrivals: [
      {
        id: 'stop-2',
        name: 'Dummy Stop 2',
        state: 'predicted',
        planned: '19:00',
      },
    ],
  },
  {
    id: 'trip-3',
    provider: 'kvg',
    direction: 'Uni/Botan. Garten',
    arrivals: [
      {
        id: 'stop-1',
        name: 'Andreas-Gayk-Straße',
        state: 'departed',
        planned: '19:23',
      },
      {
        id: 'stop-2',
        name: 'Hauptbahnhof',
        state: 'predicted',
        planned: '19:26',
      },
      {
        id: 'stop-3',
        name: 'Kirchhofallee',
        state: 'predicted',
        planned: '19:28',
      },
    ],
    path: [],
  },
];

export class DummyApi implements Api {
  isConnected = ref(true);

  useStops() {
    const stops = ref<Stop[]>(DUMMY_STOPS);

    return {
      stops,
      loading: ref(false),
      unsubscribe: async () => {},
    };
  }

  useVehicles() {
    const vehicles = ref<Vehicle[]>(DUMMY_VEHICLES);

    return {
      vehicles,
      loading: ref(false),
      unsubscribe: async () => {},
    };
  }

  useStop(stopId: Ref<string | undefined>) {
    const stop = computed(() => DUMMY_STOPS.find((s) => s.id === stopId.value) || null);

    return {
      stop,
      loading: ref(false),
      unsubscribe: async () => {},
    };
  }

  useVehicle(vehicleId: Ref<string | undefined>) {
    const vehicle = computed(() => DUMMY_VEHICLES.find((v) => v.id === vehicleId.value) || null);

    return {
      vehicle,
      loading: ref(false),
      unsubscribe: async () => {},
    };
  }

  useTrip(tripId: Ref<string | undefined>) {
    const trip = computed(() => {
      return DUMMY_TRIPS.find((t) => t.id === tripId.value) || null;
    });

    return {
      trip,
      loading: ref(false),
      unsubscribe: async () => {},
    };
  }
}
