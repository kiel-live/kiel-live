import { Icon } from 'leaflet';
import 'leaflet/dist/leaflet.css';
import markerIcon from 'leaflet/dist/images/marker-icon.png';
import markerRetinaIcon from 'leaflet/dist/images/marker-icon-2x.png';
import markerShadowIcon from 'leaflet/dist/images/marker-shadow.png';

// this part resolve an issue where the markers would not appear
delete Icon.Default.prototype._getIconUrl;

Icon.Default.mergeOptions({
  iconUrl: markerIcon,
  iconRetinaUrl: markerRetinaIcon,
  shadowUrl: markerShadowIcon,
});
