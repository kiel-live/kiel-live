import L from 'leaflet';

// general vektor marker
L.VectorMarker = L.CircleMarker.extend({
  options: {
    customDraw: (layer, ctx) => {
      const p = layer._point;
      ctx.arc(p.x, p.y, 6, 0, 2 * Math.PI);
    },
  },

  _updatePath() {
    this._renderer._updateCustomPath(this);
  },

  _customDraw(ctx) {
    ctx.save();
    this.options.customDraw(this, ctx);
    ctx.restore();
  },
});

L.vectorMarker = (latlng, options) => new L.VectorMarker(latlng, options);

// stop vector marker
L.StopMarker = L.VectorMarker.extend({
  options: {
    customDraw: (layer, ctx) => {
      const { options } = layer;
      const p = layer._point;
      const colorPrimary = 'rgb(77, 151, 255)';
      const colorFocused = 'rgb(0, 106, 255)';
      const radius = 7;

      ctx.translate(p.x, p.y);

      ctx.fillStyle = colorPrimary;
      if (options.focused) {
        ctx.scale(2, 2);
        ctx.fillStyle = colorFocused;
      }

      ctx.beginPath();
      ctx.arc(0, 0, radius, 0, 2 * Math.PI);
      ctx.fill('evenodd');
    },
  },
});

L.stopMarker = (latlng, options) => new L.StopMarker(latlng, options);

// vehicle vector marker
L.VehicleMarker = L.VectorMarker.extend({
  options: {
    customDraw: (layer, ctx) => {
      const { options } = layer;
      const p = layer._point;
      const radius = 12;
      const colorPrimary = 'rgb(170, 0, 0)';
      const colorSecondary = 'rgb(170, 100, 100)';
      const colorSelected = 'rgb(170, 50, 50)';

      ctx.translate(p.x, p.y);

      if (options.focused) {
        ctx.scale(1.5, 1.5);
      }

      if (typeof options.heading !== 'undefined' && options.heading !== null) {
        ctx.rotate((options.heading * Math.PI) / 180);
        ctx.beginPath();
        ctx.fillStyle = colorPrimary;
        const height = 10;
        const width = 12;
        ctx.moveTo(0, 0 - radius - height);
        ctx.lineTo(0 - width / 2, 0 - radius);
        ctx.lineTo(0 + width / 2, 0 - radius);
        ctx.closePath();
        ctx.fill('evenodd');
        ctx.rotate((-options.heading * Math.PI) / 180); // rotate reversed
      }

      ctx.beginPath();
      ctx.arc(0, 0, radius, 0, 2 * Math.PI);
      ctx.lineWidth = 3;
      ctx.strokeStyle = colorPrimary;
      ctx.fillStyle = colorSecondary;
      if (options.focused) {
        ctx.fillStyle = colorSelected;
      }
      ctx.fill('evenodd');
      ctx.stroke();

      if (options.label) {
        ctx.font = '10px Arial';
        ctx.fillStyle = '#fff';
        const textSize = ctx.measureText(options.label);
        ctx.fillText(options.label, 0 - textSize.width / 2, 0 + 4);
      }
    },
  },
});

L.vehicleMarker = (latlng, options) => new L.VehicleMarker(latlng, options);
