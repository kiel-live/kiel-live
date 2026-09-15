import type { StopType, VehicleType } from '~/api/types';

export const defaultMarkerColor = '#495057';

export const vehicleColors: Record<VehicleType, string> = {
  bus: '#d9480f',
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

export const stopColors: Record<StopType, string> = {
  'bus-stop': vehicleColors.bus,
  'tram-stop': vehicleColors.tram,
  'train-stop': vehicleColors.train,
  'subway-stop': vehicleColors.subway,
  'ferry-stop': vehicleColors.ferry,
  'bike-stop': vehicleColors.bike,
  'parking-spot': defaultMarkerColor,
};

// vehicle types for which showing a route/line number label next to the marker is useful
export const labeledVehicleTypes = new Set<VehicleType>(['bus', 'tram', 'train', 'subway', 'ferry']);
