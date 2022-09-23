import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:kiel_live/screens/map.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    SystemChrome.setPreferredOrientations([
      DeviceOrientation.portraitUp,
      DeviceOrientation.portraitDown,
    ]);

    return MaterialApp(
      title: 'Flutter Demo',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      initialRoute: MapScreen.routeName,
      routes: {
        MapScreen.routeName: (context) => MapScreen(),
      },
    );
  }
}
