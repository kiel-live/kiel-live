import type { VehicleType } from '~/api/types';
import { vehicleGlyphDataUrl } from '~/components/map/vehicleGlyphs';

// Vehicles are drawn from two separate map images/layers, both anchored at the
// vehicle's point:
//  - the "badge" (colored, outlined circle with the vehicle-type glyph) which
//    is always drawn upright and never rotated;
//  - the "nose" (a small pointer extending outward from the badge), rotated on
//    the map via the layer's `icon-rotate` (bound to the vehicle's heading).
// A canvas pixel is the same number of screen pixels in both images (same
// pixelRatio, same layer icon-size), so as long as each canvas is centered on
// the vehicle's point, the two line up without any manual offset — the nose
// canvas is simply drawn larger to give the pointer room without touching the
// badge's own geometry.

const BADGE_SIZE = 64;
const BADGE_SELECTED_SIZE = 84;
const NOSE_CANVAS_SIZE = 128;
const OUTLINE_COLOR = '#ffffff';
const RADIUS_RATIO = 0.6;

// radius of the (unselected) badge circle, in canvas pixels — shared with the
// nose so its pointer starts flush against the badge's edge
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
  const noseHeight = BADGE_RADIUS * 1.4;
  const noseWidth = BADGE_RADIUS;
  // start a couple of pixels inside the badge's edge so the badge (drawn on
  // top) cleanly covers the seam instead of leaving an antialiasing gap
  const baseY = -BADGE_RADIUS + 3;
  const tipY = -BADGE_RADIUS - noseHeight;

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
