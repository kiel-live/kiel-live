import { ArrivalType, VehicleType } from '~/api/types';

export type Marker = {
  type: ArrivalType | VehicleType;
  id: string;
};
