import { StopType, VehicleType } from '~/api/types';

export type Marker = {
  type: StopType | VehicleType;
  id: string;
};
