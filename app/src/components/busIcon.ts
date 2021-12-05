import maplibregl from 'maplibre-gl';

const colorPrimary = 'rgb(170, 0, 0)';
const colorSecondary = '#aaa';

export default class PulsingDot {
  width: number;
  height: number;
  data: Uint8ClampedArray;
  map: maplibregl.Map;
  focused: boolean;
  route: string;
  heading: number;
  rendered: boolean = false;
  context: CanvasRenderingContext2D | undefined;

  constructor(map: maplibregl.Map, focused: boolean, route: string, heading: number) {
    this.map = map;
    this.focused = focused;
    this.route = route;
    this.heading = heading;
    if (focused) {
      this.width = 100;
      this.height = 100;
    } else {
      this.width = 80;
      this.height = 80;
    }
    this.data = new Uint8ClampedArray(this.width * this.height * 4);
  }

  // get rendering context for the map canvas when layer is added to the map
  onAdd() {
    const canvas = document.createElement('canvas');
    canvas.width = this.width;
    canvas.height = this.height;
    this.context = canvas.getContext('2d') || undefined;
  }

  // called once before every frame where the icon will be used
  render() {
    if (this.rendered || !this.context) {
      return false;
    }
    const radius = (this.width / 2) * 0.6;
    const { context } = this;

    // clear canvas
    context.save();
    context.fillStyle = '#fff';
    context.clearRect(0, 0, this.width, this.height);
    // context.fillRect(0, 0, this.width, this.height);

    context.translate(this.width / 2, this.height / 2);

    if (this.focused) {
      // draw arrow
      context.rotate((this.heading * Math.PI) / 180);
      const lineWidth = 6;

      context.beginPath();
      context.moveTo(0, -this.height / 2 + lineWidth);
      context.lineTo(35 - lineWidth, 35 - lineWidth);
      context.lineTo(0, 25 - lineWidth);
      context.lineTo(-35 + lineWidth, 35 - lineWidth);
      context.closePath();

      context.lineWidth = lineWidth;
      context.strokeStyle = colorSecondary;
      context.stroke();

      context.fillStyle = colorPrimary;
      context.fill();

      context.rotate((-this.heading * Math.PI) / 180);
    } else {
      // draw heading nose
      if (typeof this.heading !== 'undefined' && this.heading !== null) {
        context.rotate((this.heading * Math.PI) / 180);
        context.beginPath();
        context.fillStyle = colorSecondary;
        const height = 15;
        const width = 18;
        context.moveTo(0, 0 - radius - height);
        context.lineTo(0 - width / 2, 0 - radius);
        context.lineTo(0 + width / 2, 0 - radius);
        context.closePath();
        context.fill('evenodd');
        context.rotate((-this.heading * Math.PI) / 180);
      }

      // draw base (circle)
      context.beginPath();
      context.arc(0, 0, radius, 0, 2 * Math.PI);
      context.lineWidth = 4;
      context.strokeStyle = colorSecondary;
      context.fillStyle = colorPrimary;
      context.fill('evenodd');
      context.stroke();
    }
    // draw text (route)
    context.fillStyle = '#eee';
    context.font = '20px Arial';
    context.textAlign = 'center';
    context.textBaseline = 'middle';
    context.fillText(this.route, 0, 0);

    context.restore();

    // update this image's data with data from the canvas
    this.data = context.getImageData(0, 0, this.width, this.height).data;
    this.rendered = true;

    // return `true` to let the map know that the image was updated
    return true;
  }
}
