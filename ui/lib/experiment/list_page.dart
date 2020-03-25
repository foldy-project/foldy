
import 'package:flutter/material.dart';

import 'list.dart';

class ExperimentsPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Column(
        mainAxisAlignment: MainAxisAlignment.start,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: <Widget>[
          ExperimentListHeader(),
          Expanded(
            child: ExperimentList(items: [
              Experiment(name: "my-experiment-0", namespace: "default"),
              Experiment(name: "my-experiment-1", namespace: "default"),
              Experiment(name: "my-experiment-2", namespace: "default"),
            ]),
          ),
        ],
      ),
    );
  }
}
