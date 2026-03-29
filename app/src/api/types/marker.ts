import type { StopType, VehicleType } from '~/api/types';

export interface Marker {
  type: StopType | VehicleType | 'stop' | 'vehicle';
  id: string;
}
