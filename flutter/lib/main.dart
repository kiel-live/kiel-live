import 'package:flutter/material.dart';
import 'package:kiel_live/views/home.dart';
import 'package:kiel_live/views/map.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter Demo',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      home: Map(),
      // home: const MyHomePage(title: 'Flutter Demo Home Page'),
    );
  }
}
