import 'package:flutter/material.dart';

class ExperimentsDetailsPage extends StatelessWidget {
  final String name;

  final String namespace;

  const ExperimentsDetailsPage({
    Key key,
    @required this.namespace,
    @required this.name,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Column(
        mainAxisAlignment: MainAxisAlignment.start,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: <Widget>[],
      ),
    );
  }
}
