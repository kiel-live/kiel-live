import type { VehicleType } from '~/api/types';

export const defaultMarkerColor = '#495057';

// a stop can serve multiple vehicle types at once (e.g. Hauptbahnhof is a bus,
// train and bike stop), so stops all share one neutral color rather than being
// colored per type.
export const stopColor = '#33414d';

export const vehicleColors: Record<VehicleType, string> = {
  bus: '#aa0000',
  tram: '#862e9c',
  train: '#1864ab',
  subway: '#0b7285',
  ferry: '#0c8599',
  bike: '#2b8a3e',
  'e-scooter': '#f08c00',
  car: defaultMarkerColor,
  moped: defaultMarkerColor,
  'e-moped': defaultMarkerColor,
};

// vehicle types for which showing a route/line number label next to the marker is useful
export const labeledVehicleTypes = new Set<VehicleType>(['bus', 'tram', 'train', 'subway', 'ferry']);
