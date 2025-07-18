import 'package:flutter/material.dart';
import 'dart:io';
import 'dart:math';

import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart'; // ignore: unnecessary_import
import 'package:kiel_live/widgets/map/map.dart';
import 'package:kiel_live/widgets/popups/BusStopPopup.dart';
import 'package:kiel_live/widgets/popups/Popup.dart';
import 'package:location/location.dart';
import 'package:sliding_up_panel/sliding_up_panel.dart';

class MapScreen extends StatefulWidget {
  static const String routeName = '/map';

  const MapScreen({super.key});

  @override
  MapScreenState createState() => MapScreenState();
}

class MapScreenState extends State<MapScreen> {
  final PanelController panelController = PanelController();

  @override
  Widget build(BuildContext context) {
    var height =
        MediaQuery.of(context).size.height - MediaQuery.of(context).padding.top;

    return Scaffold(
      floatingActionButton: FloatingActionButton(
        onPressed: () async {
          var location = Location();
          var hasPermissions = await location.hasPermission();
          if (hasPermissions != PermissionStatus.granted) {
            await location.requestPermission();
          } else {
            print('done');
          }
        },
        child: const Icon(Icons.search),
      ),
      body: SlidingUpPanel(
        controller: panelController,
        minHeight: height * 0.2,
        maxHeight: height * 1,
        borderRadius: const BorderRadius.vertical(top: Radius.circular(15)),
        color: Colors.grey[900]!,
        panelBuilder: (controller) =>
            Popup(controller: controller, panelController: panelController),
        body: const Map(),
      ),
    );
  }
}
