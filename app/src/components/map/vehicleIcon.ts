// Draws small colored, outlined dot markers for vehicles, optionally with a "nose"
// triangle pointing north. The nose is rotated on the map via the layer's
// `icon-rotate` property (bound to the vehicle's heading), so the icon itself is
// only ever drawn pointing "up".

const SIZE = 64;
const SELECTED_SIZE = 84;
const OUTLINE_COLOR = '#ffffff';

export interface VehicleIconOptions {
  color: string;
  selected?: boolean;
  nose?: boolean;
}

export function createVehicleIcon({ color, selected = false, nose = false }: VehicleIconOptions): ImageData {
  const size = selected ? SELECTED_SIZE : SIZE;
  const canvas = document.createElement('canvas');
  canvas.width = size;
  canvas.height = size;

  const context = canvas.getContext('2d');
  if (!context) {
    throw new Error('Could not create 2d canvas context for vehicle icon');
  }

  const center = size / 2;
  const radius = center * 0.55;
  const strokeWidth = selected ? 5 : 4;

  context.translate(center, center);

  if (nose) {
    const noseHeight = radius * 1.15;
    const noseWidth = radius * 1.1;

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
  }

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

  return context.getImageData(0, 0, size, size);
}
