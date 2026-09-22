import type { VehicleType } from '~/api/types';

// Same icons as VehiclePopup.vue, imported directly (rather than as Vue
// components) so only the used icons are bundled and we get the raw markup
// to recolor for the canvas-drawn map markers.
import bicycleSvg from '~icons/carbon/bicycle?raw';
import trainProfileSvg from '~icons/carbon/train-profile?raw';
import busSvg from '~icons/mdi/bus?raw';
import ferrySvg from '~icons/mdi/ferry?raw';
import mopedElectricSvg from '~icons/mdi/moped-electric?raw';
import mopedSvg from '~icons/mdi/moped?raw';
import carSvg from '~icons/ph/car?raw';
import scooterSvg from '~icons/ph/scooter?raw';
import subwaySvg from '~icons/ph/subway?raw';
import tramSvg from '~icons/ph/tram?raw';

const glyphs: Record<VehicleType, string> = {
  bus: busSvg,
  bike: bicycleSvg,
  car: carSvg,
  'e-scooter': scooterSvg,
  ferry: ferrySvg,
  train: trainProfileSvg,
  subway: subwaySvg,
  tram: tramSvg,
  moped: mopedSvg,
  'e-moped': mopedElectricSvg,
};

export function vehicleGlyphDataUrl(type: VehicleType, color: string): string {
  // the raw icon import has no xmlns, which <img>-loaded SVGs require to render
  const svg = glyphs[type].replace('<svg', `<svg xmlns="http://www.w3.org/2000/svg" style="color:${color}"`);
  return `data:image/svg+xml;charset=utf-8,${encodeURIComponent(svg)}`;
}
