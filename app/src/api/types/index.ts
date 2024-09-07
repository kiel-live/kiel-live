import type { Stop } from './stop';
import type { Trip } from './trip';
import type { Vehicle } from './vehicle';

export * from './arrival';
export * from './location';
export * from './marker';
export * from './stop';
export * from './trip';
export * from './vehicle';

export type Models = Vehicle | Stop | Trip;
