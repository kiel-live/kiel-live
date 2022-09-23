import 'package:flutter/material.dart';
import 'dart:io';
import 'dart:math';

import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart'; // ignore: unnecessary_import
import 'package:kiel_live/widgets/popups/BusStopPopup.dart';
import 'package:kiel_live/widgets/popups/Popup.dart';
import 'package:maplibre_gl/mapbox_gl.dart';
import 'package:sliding_up_panel/sliding_up_panel.dart';

class Map extends StatefulWidget {
  const Map();

  @override
  State createState() => MapState();
}

class MapState extends State<Map> {
  MaplibreMapController? mapController;
  var isLight = true;
  final String mapStyle = "https://tiles.ju60.de/styles/gray-matter/style.json";

  _onMapCreated(MaplibreMapController controller) {
    mapController = controller;
    controller.addCircle(
      const CircleOptions(
        geometry: LatLng(54.3166, 10.1283),
        circleRadius: 10,
        circleColor: "#ff0000",
        circleOpacity: 0.5,
      ),
    );
    controller.addSymbol(
      const SymbolOptions(
        geometry: LatLng(54.3166, 10.1283),
        iconImage: "bus-stop",
        iconSize: 1.5,
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return MaplibreMap(
      styleString: mapStyle,
      onMapCreated: _onMapCreated,
      initialCameraPosition:
          const CameraPosition(target: LatLng(54.3166, 10.1283), zoom: 11.0),
      myLocationEnabled: true,
      myLocationRenderMode: MyLocationRenderMode.GPS,
      myLocationTrackingMode: MyLocationTrackingMode.Tracking,
      doubleClickZoomEnabled: true,
    );
  }
}
