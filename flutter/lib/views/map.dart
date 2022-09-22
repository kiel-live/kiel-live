import 'package:flutter/material.dart';
import 'dart:io';
import 'dart:math';

import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart'; // ignore: unnecessary_import
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
      body: SlidingUpPanel(
        minHeight: MediaQuery.of(context).size.height * 0.2,
        maxHeight: MediaQuery.of(context).size.height * 1,
        borderRadius: const BorderRadius.vertical(top: Radius.circular(15)),
        panel: const Center(
          child: Text("This is the sliding Widget"),
        ),
        body: MaplibreMap(
          styleString: "https://tiles.ju60.de/styles/gray-matter/style.json",
          onMapCreated: _onMapCreated,
          initialCameraPosition: const CameraPosition(
              target: LatLng(54.3166, 10.1283), zoom: 11.0),
          onStyleLoadedCallback: _onStyleLoadedCallback,
        ),
      ),
    );
  }
}
