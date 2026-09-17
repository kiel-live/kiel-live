import type { VehicleType } from '~/api/types';

export const defaultMarkerColor = '#495057';

// one shared color: a stop can serve multiple vehicle types at once (e.g.
// Hauptbahnhof is a bus, train and bike stop), so coloring by type doesn't work
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

// types with a meaningful route/line number to label the marker with
export const labeledVehicleTypes = new Set<VehicleType>(['bus', 'tram', 'train', 'subway', 'ferry']);
