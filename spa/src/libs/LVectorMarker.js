import L from 'leaflet';

const CustomCanvas = L.Canvas.extend({
  _updateCustomPath(layer) {
    if (!this._drawing || layer._empty()) { return; }

    const ctx = this._ctx;
    layer._customDraw(ctx);
  },
});

// general vektor marker
L.VectorMarker = L.CircleMarker.extend({
  options: {
    renderer: new CustomCanvas(),
    draw: (layer, ctx) => {
      const p = layer._point;
      ctx.arc(p.x, p.y, 6, 0, 2 * Math.PI);
    },
  },

  _updatePath() {
    this._renderer._updateCustomPath(this);
  },

  _customDraw(ctx) {
    this.options.draw(this, ctx);
  },
});

L.vectorMarker = (latlng, options) => new L.VectorMarker(latlng, options);

// stop vector marker
L.StopMarker = L.VectorMarker.extend({
  options: {
    draw: (layer, ctx) => {
      const { options } = layer;
      const p = layer._point;

      ctx.beginPath();
      ctx.arc(p.x, p.y, layer._radius, 0, 2 * Math.PI);
      ctx.globalAlpha = options.fillOpacity;
      ctx.fillStyle = options.fillColor || options.color;
      ctx.fill('evenodd');
    },
  },
});

L.stopMarker = (latlng, options) => new L.StopMarker(latlng, options);

// vehicle vector marker
L.VehicleMarker = L.VectorMarker.extend({
  options: {
    label: 'test',
    draw: (layer, ctx) => {
      const { options } = layer;
      const p = layer._point;

      ctx.beginPath();
      ctx.arc(p.x, p.y, layer._radius, 0, 2 * Math.PI);
      ctx.globalAlpha = options.fillOpacity;
      ctx.fillStyle = options.fillColor || options.color;
      ctx.fill('evenodd');

      ctx.font = '10px Arial';
      ctx.fillStyle = '#fff';
      const textSize = ctx.measureText(options.label);
      ctx.fillText(options.label, p.x - textSize.width / 2, p.y + 4);
    },
  },
});

L.vehicleMarker = (latlng, options) => new L.VehicleMarker(latlng, options);
