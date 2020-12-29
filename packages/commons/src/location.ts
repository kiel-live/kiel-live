export type Location = {
  longitude: number;
  latitude: number;
};

export type LocationHeading = Location & {
  heading: number;
};
