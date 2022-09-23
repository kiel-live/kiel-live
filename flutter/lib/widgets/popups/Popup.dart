import 'package:flutter/material.dart';
import 'package:kiel_live/widgets/popups/BusStopPopup.dart';
import 'package:sliding_up_panel/sliding_up_panel.dart';

class Popup extends StatelessWidget {
  final ScrollController controller;
  final PanelController panelController;

  const Popup({key, required this.controller, required this.panelController})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return ListView(
      padding: EdgeInsets.zero,
      controller: controller,
      children: [
        GestureDetector(
          child: Column(children: [
            const SizedBox(height: 12),
            Container(
              width: 30,
              height: 5,
              decoration: BoxDecoration(
                color: Colors.grey[700],
                borderRadius: BorderRadius.circular(12),
              ),
            ),
          ]),
          onTap: () => panelController.isPanelOpen
              ? panelController.close()
              : panelController.open(),
        ),
        BusStop(),
      ],
    );
  }
}
