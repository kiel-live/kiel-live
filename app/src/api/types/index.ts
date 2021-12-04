import { Stop } from './stop';
import { Trip } from './trip';
import { Vehicle } from './vehicle';

export * from './location';
export * from './vehicle';
export * from './stop';
export * from './arrival';
export * from './trip';

export type Models = Vehicle | Stop | Trip;
