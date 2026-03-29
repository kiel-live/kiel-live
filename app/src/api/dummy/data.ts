import type { Stop, Trip, Vehicle } from '../types';

function inMinutes(minutes: number): string {
  const future = new Date(Date.now() + minutes * 60 * 1000);
  return future.toISOString();
}

export const DUMMY_VEHICLES: Vehicle[] = [
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
  {
    id: 'kvg--7638104967632549039',
    provider: 'kvg',
    name: '12 Strande via Hohenleuchte',
    type: 'bus',
    state: 'onfire',
    battery: '',
    location: {
      longitude: 36472522,
      latitude: 195534934,
      heading: 0,
    },
    tripId: 'kvg-1610077840790697737',
    description: '',
  },
  {
    id: 'kvg--7638104967632549161',
    provider: 'kvg',
    name: '11 Wik Kanal',
    type: 'bus',
    state: 'onfire',
    battery: '',
    location: {
      longitude: 36513870,
      latitude: 195519676,
      heading: 270,
    },
    tripId: 'kvg-1610077840790681351',
    description: '',
  },
];

export const DUMMY_STOPS: Stop[] = [
  {
    id: 'stop-1',
    provider: 'dummy',
    name: 'Dummy Stop 1',
    type: 'bus-stop',
    routes: ['1', '2'],
    alerts: ['Alert 1'],
    location: { latitude: 54.3233 * 3600000, longitude: 10.1228 * 3600000, heading: 0 },
    actions: [],
    departures: [
      {
        name: 'Dummy Vehicle 1',
        type: 'bus',
        vehicleId: 'bus-1',
        tripId: 'trip-1',
        routeId: '1',
        routeName: 'Route 1',
        direction: 'North',
        state: 'predicted',
        planned: inMinutes(5),
        actual: inMinutes(5 + 3), // 3 minutes delay
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
    departures: [
      {
        name: 'Dummy Vehicle 2',
        type: 'bus',
        vehicleId: 'bus-2',
        tripId: 'trip-2',
        routeId: '2',
        routeName: 'Route 2',
        direction: 'South',
        state: 'predicted',
        planned: inMinutes(7),
        actual: inMinutes(7 + 1), // 1 minute delay
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
    departures: [
      {
        name: 'Hassee',
        type: 'bus',
        vehicleId: 'bus-2',
        tripId: 'trip-2',
        routeId: '1610073983892324369',
        routeName: '51',
        direction: 'Hassee',
        state: 'predicted',
        planned: inMinutes(0),
        actual: inMinutes(0 + 1), // 1 minute delay
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
        planned: inMinutes(6),
        actual: inMinutes(6 - 2), // 2 minutes early
        platform: '',
      },
    ],
    location: {
      longitude: 36460935,
      latitude: 195591727,
      heading: 0,
    },
  },
  {
    id: 'kvg-2387',
    provider: 'kvg',
    name: 'Hauptbahnhof',
    type: 'bus-stop',
    routes: ['80', '11', '12', '51', '91'],
    alerts: ['Alert 1', 'Alert 2'],
    departures: [
      {
        name: 'Strande',
        type: 'bus',
        vehicleId: 'kvg--7638104967632549039',
        tripId: 'kvg-1610077840790697737',
        routeId: '1610073983892324353',
        routeName: '12',
        direction: 'Strande',
        state: 'predicted',
        planned: inMinutes(8),
        actual: inMinutes(8 + 1), // 1 minute delay
        platform: '',
      },
      {
        name: 'Wik Kanal',
        type: 'bus',
        vehicleId: 'kvg--7638104967632549161',
        tripId: 'kvg-1610077840790681351',
        routeId: '1610073983892324359',
        routeName: '11',
        direction: 'Wik Kanal',
        state: 'predicted',
        planned: inMinutes(15),
        actual: inMinutes(15 + 23), // 23 minutes delay
        platform: '',
      },
      {
        name: 'Hassee',
        type: 'bus',
        vehicleId: 'bus-3',
        tripId: 'trip-3',
        routeId: '1610073983892324369',
        routeName: '51',
        direction: 'Hassee',
        state: 'predicted',
        planned: inMinutes(10),
        actual: inMinutes(10 - 2), // 2 minutes early
        platform: '',
      },
      {
        name: 'Wik',
        type: 'bus',
        vehicleId: 'bus-2',
        tripId: 'trip-2',
        routeId: '1610073983892324385',
        routeName: '80',
        direction: 'Wik',
        state: 'predicted',
        planned: inMinutes(78),
        actual: inMinutes(78 + 3), // 3 minutes delay
        platform: '',
      },
    ],
    location: {
      longitude: 36472006,
      latitude: 195536026,
      heading: 0,
    },
  },
  {
    id: 'kvg-1256',
    provider: 'kvg',
    name: 'Andreas-Gayk-Straße',
    type: 'bus-stop',
    routes: null,
    alerts: [],
    departures: [
      {
        name: 'Strande',
        type: 'bus',
        vehicleId: 'kvg--7638104967632549039',
        tripId: 'kvg-1610077840790697737',
        routeId: '1610073983892324353',
        routeName: '12',
        direction: 'Strande',
        state: 'predicted',
        planned: inMinutes(1),
        actual: inMinutes(1),
        platform: '',
      },
      {
        name: 'Wik Kanal',
        type: 'bus',
        vehicleId: 'kvg--7638104967632549161',
        tripId: 'kvg-1610077840790681351',
        routeId: '1610073983892324359',
        routeName: '11',
        direction: 'Wik Kanal',
        state: 'predicted',
        planned: inMinutes(90), // in 1.5 hours
        actual: inMinutes(90 + 5), // 5 minutes delay
        platform: '',
      },
    ],
    location: {
      longitude: 36482329,
      latitude: 195548642,
      heading: 0,
    },
  },
];

export const DUMMY_TRIPS: Trip[] = [
  {
    id: 'trip-1',
    provider: 'dummy',
    direction: 'North',
    path: [{ latitude: 54.3233 * 3600000, longitude: 10.1228 * 3600000, heading: 0 }],
    departures: [
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
    departures: [
      {
        id: 'stop-2',
        name: 'Dummy Stop 2',
        state: 'predicted',
        planned: inMinutes(19),
      },
    ],
  },
  {
    id: 'trip-3',
    provider: 'kvg',
    direction: 'Uni/Botan. Garten',
    departures: [
      {
        id: 'kvg-1256',
        name: 'Andreas-Gayk-Straße',
        state: 'departed',
        planned: inMinutes(-4),
      },
      {
        id: 'kvg-2387',
        name: 'Hauptbahnhof',
        state: 'predicted',
        planned: inMinutes(1),
      },
      {
        id: 'stop-3',
        name: 'Kirchhofallee',
        state: 'predicted',
        planned: inMinutes(5),
      },
    ],
    path: [],
  },
  {
    id: 'kvg-1610077840790697737',
    provider: 'kvg',
    direction: 'Strande',
    departures: [
      {
        id: 'kvg-2387',
        name: 'Hauptbahnhof',
        state: 'departed',
        planned: inMinutes(-6),
      },
      {
        id: 'kvg-1256',
        name: 'Andreas-Gayk-Straße',
        state: 'stopping',
        planned: inMinutes(0),
      },
    ],
    path: [],
  },
  {
    id: 'kvg-1610077840790681351',
    provider: 'kvg',
    direction: 'Wik Kanal',
    departures: [
      {
        id: 'kvg-2387',
        name: 'Hauptbahnhof',
        state: 'predicted',
        planned: inMinutes(12),
      },
      {
        id: 'kvg-1256',
        name: 'Andreas-Gayk-Straße',
        state: 'predicted',
        planned: inMinutes(18),
      },
    ],
    path: [],
  },
];
