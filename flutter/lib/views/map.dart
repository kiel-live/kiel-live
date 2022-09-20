import 'package:flutter/material.dart';
import 'dart:io';
import 'dart:math';

import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart'; // ignore: unnecessary_import
import 'package:maplibre_gl/mapbox_gl.dart';

class Map extends StatefulWidget {
  const Map();

  @override
  State createState() => MapState();
}

class MapState extends State<Map> {
  MaplibreMapController? mapController;
  var isLight = true;

  _onMapCreated(MaplibreMapController controller) {
    mapController = controller;
  }

  _onStyleLoadedCallback() {
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      content: Text("Style loaded :)"),
      backgroundColor: Theme.of(context).primaryColor,
      duration: Duration(seconds: 1),
    ));
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Stack(children: [
        MaplibreMap(
          styleString: "https://tiles.ju60.de/styles/gray-matter/style.json",
          onMapCreated: _onMapCreated,
          initialCameraPosition:
              const CameraPosition(target: LatLng(10.1283, 54.3166)),
          onStyleLoadedCallback: _onStyleLoadedCallback,
        ),
      ]),
    );
  }
}
