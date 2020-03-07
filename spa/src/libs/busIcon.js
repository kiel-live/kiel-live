const colorPrimary = 'rgb(170, 0, 0)';
const colorSecondary = 'rgb(170, 100, 100)';
const colorFocused = 'rgb(170, 50, 50)';

export default class PulsingDot {
  width = 30;

  height = 40;

  data = new Uint8Array(this.width * this.height * 4);

  map;

  focused;

  constructor(map, focused) {
    this.map = map;
    this.focused = focused;
  }

  // get rendering context for the map canvas when layer is added to the map
  onAdd() {
    const canvas = document.createElement('canvas');
    canvas.width = this.width;
    canvas.height = this.height;
    this.context = canvas.getContext('2d');
  }

  // called once before every frame where the icon will be used
  render() {
    const radius = (this.width / 2) * 0.7;
    const { context } = this;

    // clear canvas
    context.save();
    context.fillStyle = '#fff';
    context.clearRect(0, 0, this.width, this.height);
    // context.fillRect(0, 0, this.width, this.height);

    context.translate(this.width / 2, this.height / 2);

    // draw heading nose
    context.beginPath();
    context.fillStyle = colorPrimary;
    const height = 10;
    const width = 12;
    context.moveTo(0, 0 - radius - height);
    context.lineTo(0 - width / 2, 0 - radius);
    context.lineTo(0 + width / 2, 0 - radius);
    context.closePath();
    context.fill('evenodd');

    // draw base (circle)
    context.beginPath();
    context.arc(0, 0, radius, 0, 2 * Math.PI);
    context.lineWidth = 3;
    context.strokeStyle = colorPrimary;
    context.fillStyle = colorSecondary;
    if (this.focused) {
      context.scale(2, 2);
      context.fillStyle = colorFocused;
    }
    context.fill('evenodd');
    context.stroke();

    context.restore();

    // update this image's data with data from the canvas
    this.data = context.getImageData(
      0,
      0,
      this.width,
      this.height,
    ).data;

    // return `true` to let the map know that the image was updated
    return true;
  }
}
