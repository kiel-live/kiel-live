import type { VehicleType } from '~/api/types';
import { vehicleGlyphDataUrl } from '~/components/map/vehicleGlyphs';

// Vehicles are drawn as two separate map images/layers sharing the same canvas
// size and center point:
//  - the "badge" (colored, outlined circle with the vehicle-type glyph) which
//    is always drawn upright and never rotated;
//  - the "nose" (a small triangle pointing outward from the badge), rotated on
//    the map via the layer's `icon-rotate` (bound to the vehicle's heading).
// Because both share the same canvas size and are rendered at the same
// `icon-size`, they line up perfectly at the feature's point without either
// layer needing a manual offset.

const SIZE = 64;
const SELECTED_SIZE = 84;
const OUTLINE_COLOR = '#ffffff';
const RADIUS_RATIO = 0.6;

export interface VehicleBadgeIconOptions {
  type: VehicleType;
  color: string;
  selected?: boolean;
}

function loadImage(url: string): Promise<HTMLImageElement> {
  return new Promise((resolve, reject) => {
    const image = new Image();
    image.onload = () => resolve(image);
    image.onerror = () => reject(new Error(`Could not load vehicle icon image: ${url}`));
    image.src = url;
  });
}

export async function createVehicleBadgeIcon({ type, color, selected = false }: VehicleBadgeIconOptions) {
  const size = selected ? SELECTED_SIZE : SIZE;
  const canvas = document.createElement('canvas');
  canvas.width = size;
  canvas.height = size;

  const context = canvas.getContext('2d');
  if (!context) {
    throw new Error('Could not create 2d canvas context for vehicle icon');
  }

  const center = size / 2;
  const radius = center * RADIUS_RATIO;
  const strokeWidth = selected ? 5 : 4;

  const glyph = await loadImage(vehicleGlyphDataUrl(type, OUTLINE_COLOR));

  context.translate(center, center);

  context.beginPath();
  context.arc(0, 0, radius, 0, Math.PI * 2);
  context.fillStyle = color;
  context.fill();
  context.lineWidth = strokeWidth;
  context.strokeStyle = OUTLINE_COLOR;
  context.stroke();

  if (selected) {
    context.beginPath();
    context.arc(0, 0, radius + strokeWidth / 2 + 3, 0, Math.PI * 2);
    context.lineWidth = 2;
    context.strokeStyle = color;
    context.stroke();
  }

  const glyphSize = radius * 1.15;
  context.drawImage(glyph, -glyphSize / 2, -glyphSize / 2, glyphSize, glyphSize);

  return context.getImageData(0, 0, size, size);
}

export function createVehicleNoseIcon(color: string) {
  const canvas = document.createElement('canvas');
  canvas.width = SIZE;
  canvas.height = SIZE;

  const context = canvas.getContext('2d');
  if (!context) {
    throw new Error('Could not create 2d canvas context for vehicle nose icon');
  }

  const center = SIZE / 2;
  const radius = center * RADIUS_RATIO;
  const noseHeight = radius * 1.15;
  const noseWidth = radius * 1.1;

  context.translate(center, center);

  context.beginPath();
  context.moveTo(0, -radius - noseHeight);
  context.lineTo(-noseWidth / 2, -radius + 3);
  context.lineTo(noseWidth / 2, -radius + 3);
  context.closePath();
  context.fillStyle = color;
  context.fill();
  context.lineWidth = 2;
  context.strokeStyle = OUTLINE_COLOR;
  context.stroke();

  return context.getImageData(0, 0, SIZE, SIZE);
}
