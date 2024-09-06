export type Bounds = {
  north: number;
  south: number;
  east: number;
  west: number;
};

export type GpsLocation = {
  longitude: number; // exp: 54.306 * 3600000 = longitude
  latitude: number; // exp: 10.149 * 3600000 = latitude
  heading: number; // in degree
};
