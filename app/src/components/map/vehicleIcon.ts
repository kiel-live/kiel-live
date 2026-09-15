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
const OUTLINE_COLOR = '#ffffff';
const RADIUS_RATIO = 0.6;

// unselected badge circle radius, in canvas pixels — shared with the nose
const BADGE_RADIUS = (BADGE_SIZE / 2) * RADIUS_RATIO;

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
  canvas.width = NOSE_CANVAS_SIZE;
  canvas.height = NOSE_CANVAS_SIZE;

  const context = canvas.getContext('2d');
  if (!context) {
    throw new Error('Could not create 2d canvas context for vehicle nose icon');
  }

  const center = NOSE_CANVAS_SIZE / 2;
  const noseHeight = BADGE_RADIUS * 0.75;
  const noseWidth = BADGE_RADIUS * 0.7;
  const gap = BADGE_RADIUS * 0.25;
  const baseY = -BADGE_RADIUS - gap;
  const tipY = baseY - noseHeight;

  context.translate(center, center);

  context.beginPath();
  context.moveTo(0, tipY);
  context.lineTo(-noseWidth / 2, baseY);
  context.lineTo(noseWidth / 2, baseY);
  context.closePath();
  context.fillStyle = color;
  context.fill();
  context.lineWidth = 2;
  context.lineJoin = 'round';
  context.strokeStyle = OUTLINE_COLOR;
  context.stroke();

  return context.getImageData(0, 0, NOSE_CANVAS_SIZE, NOSE_CANVAS_SIZE);
}
