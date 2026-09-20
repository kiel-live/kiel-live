import type { VehicleType } from '~/api/types';
import { vehicleGlyphDataUrl } from '~/components/map/vehicleGlyphs';

// The badge (circle + glyph) and nose (heading pointer) are separate images so
// the nose alone can be rotated via icon-rotate. Both are anchored at the
// vehicle's point and drawn at the same pixelRatio/icon-size, so a canvas
// pixel is the same screen size in both — they line up without an offset, and
// the nose canvas can simply be made bigger to give its pointer more room.

const BADGE_SIZE = 64;
const BADGE_SELECTED_SIZE = 84;
const NOSE_CANVAS_SIZE = 128;
const RADIUS_RATIO = 0.6;

// unselected badge circle radius, in canvas pixels — shared with the nose
const BADGE_RADIUS = (BADGE_SIZE / 2) * RADIUS_RATIO;

export interface VehicleBadgeIconOptions {
  type: VehicleType;
  color: string;
  // dimmed in dark mode so the outline doesn't glare against the dark basemap
  outlineColor: string;
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

export async function createVehicleBadgeIcon({ type, color, outlineColor, selected = false }: VehicleBadgeIconOptions) {
  const size = selected ? BADGE_SELECTED_SIZE : BADGE_SIZE;
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

  const glyph = await loadImage(vehicleGlyphDataUrl(type, outlineColor));

  context.translate(center, center);

  context.beginPath();
  context.arc(0, 0, radius, 0, Math.PI * 2);
  context.fillStyle = color;
  context.fill();
  context.lineWidth = strokeWidth;
  context.strokeStyle = outlineColor;
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

// scales up with the badge so the selected vehicle's nose reads as
// highlighted too, and clears the badge's selected halo ring
const NOSE_SELECTED_SCALE = BADGE_SELECTED_SIZE / BADGE_SIZE;

export function createVehicleNoseIcon(color: string, outlineColor: string, selected = false) {
  const canvas = document.createElement('canvas');
  canvas.width = NOSE_CANVAS_SIZE;
  canvas.height = NOSE_CANVAS_SIZE;

  const context = canvas.getContext('2d');
  if (!context) {
    throw new Error('Could not create 2d canvas context for vehicle nose icon');
  }

  const radius = selected ? BADGE_RADIUS * NOSE_SELECTED_SCALE : BADGE_RADIUS;
  const center = NOSE_CANVAS_SIZE / 2;
  const noseHeight = radius * 0.75;
  const noseWidth = radius * 0.7;
  const gap = radius * 0.25;
  const baseY = -radius - gap;
  const tipY = baseY - noseHeight;

  context.translate(center, center);

  context.beginPath();
  context.moveTo(0, tipY);
  context.lineTo(-noseWidth / 2, baseY);
  context.lineTo(noseWidth / 2, baseY);
  context.closePath();
  context.fillStyle = color;
  context.fill();
  context.lineWidth = selected ? 3 : 2;
  context.lineJoin = 'round';
  context.strokeStyle = outlineColor;
  context.stroke();

  return context.getImageData(0, 0, NOSE_CANVAS_SIZE, NOSE_CANVAS_SIZE);
}
