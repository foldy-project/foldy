import 'package:fluro/fluro.dart';
import 'package:flutter/material.dart';
import 'dart:html' as html;

import 'package:get_it/get_it.dart';

class Experiment {
  String name;
  String namespace;
  String status;

  Experiment({
    @required this.name,
    @required this.namespace,
  });
}

class ExperimentListItemProperty extends StatelessWidget {
  final String name;

  const ExperimentListItemProperty(this.name);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(6.0),
      child: Text(name),
    );
  }
}

class Clickable extends StatelessWidget {
  static final appContainer =
      html.window.document.getElementById('app-container');

  final void Function() onTap;

  final Widget child;

  const Clickable({
    Key key,
    @required this.child,
    @required this.onTap,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return MouseRegion(
      onHover: (event) {
        appContainer.style.cursor = 'pointer';
      },
      // When it exits set it back to default
      onExit: (event) {
        appContainer.style.cursor = 'default';
      },
      child: GestureDetector(
        onTap: onTap,
        child: child,
      ),
    );
  }
}

class ExperimentListItem extends StatelessWidget {
  final Experiment experiment;

  const ExperimentListItem({
    Key key,
    @required this.experiment,
  }) : super(key: key);

  void navigateTo(BuildContext context) =>GetIt.I<Router>().navigateTo(context, "/experiments/${experiment.namespace}/${experiment.name}");

  @override
  Widget build(BuildContext context) {
    Color indicatorColor;
    switch (experiment.status) {
      case "Pending":
        indicatorColor = Colors.yellow;
        break;
      case "Healthy":
        indicatorColor = Colors.green;
        break;
      case "Error":
        indicatorColor = Theme.of(context).errorColor;
        break;
      default:
        indicatorColor = Colors.blue;
        break;
    }
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 6),
      child: Clickable(
        onTap: () => navigateTo(context),
        child: Card(
          child: Container(
            decoration: BoxDecoration(
              border: Border(
                left: BorderSide(
                  color: indicatorColor,
                  width: 4.0,
                  style: BorderStyle.solid,
                ),
              ),
            ),
            child: Padding(
              padding: const EdgeInsets.all(8.0),
              child: Container(
                child: Row(
                  children: <Widget>[
                    Expanded(
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: <Widget>[
                          Container(
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: <Widget>[
                                ExperimentListItemProperty(experiment.name),
                                ExperimentListItemProperty(
                                    experiment.namespace),
                              ],
                            ),
                          ),
                          Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: <Widget>[
                              ExperimentListItemProperty("proteinnet-casp11"),
                              ExperimentListItemProperty("gromacs"),
                            ],
                          ),
                        ],
                      ),
                    ),
                    Container(
                      decoration: BoxDecoration(
                        borderRadius: BorderRadius.all(Radius.circular(8.0)),
                      ),
                      child: ClipRRect(
                        borderRadius: BorderRadius.circular(8.0),
                        child: Material(
                          color: Theme.of(context).cardColor,
                          child: InkWell(
                            onTap: () {},
                            child: Padding(
                              padding: const EdgeInsets.symmetric(
                                  vertical: 6.0, horizontal: 2.0),
                              child: Icon(Icons.more_vert),
                            ),
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class ExperimentList extends StatelessWidget {
  final List<Experiment> items;

  const ExperimentList({
    Key key,
    @required this.items,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      color: Theme.of(context).secondaryHeaderColor,
      child: Padding(
        padding: const EdgeInsets.symmetric(
          vertical: 8.0,
          horizontal: 16.0,
        ),
        child: Column(
          children:
              items.map((e) => ExperimentListItem(experiment: e)).toList(),
        ),
      ),
    );
  }
}

class ExperimentListHeader extends StatelessWidget {
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: Theme.of(context).scaffoldBackgroundColor,
        border: Border(
          bottom: BorderSide(
            color: Theme.of(context).dividerColor,
          ),
        ),
      ),
      child: Padding(
        padding: const EdgeInsets.all(8.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: <Widget>[
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: <Widget>[
                Padding(
                  padding: const EdgeInsets.all(8.0),
                  child: Text(
                    "Experiments",
                    style: Theme.of(context).textTheme.subtitle2,
                  ),
                ),
                Padding(
                  padding: const EdgeInsets.all(8.0),
                  child: Text(
                    "EXPERIMENTS",
                    style: Theme.of(context).textTheme.subtitle2,
                  ),
                ),
              ],
            ),
            Divider(),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: <Widget>[
                Padding(
                  padding: const EdgeInsets.all(4.0),
                  child: ClipRRect(
                    borderRadius: BorderRadius.circular(16.0),
                    child: Material(
                      color: Theme.of(context).buttonColor,
                      child: InkWell(
                        onTap: () {},
                        child: Padding(
                          padding: const EdgeInsets.symmetric(
                            vertical: 10.0,
                            horizontal: 12.0,
                          ),
                          child: Text(
                            "New Experiment",
                            style:
                                Theme.of(context).textTheme.bodyText2.copyWith(
                                      fontWeight: FontWeight.bold,
                                    ),
                          ),
                        ),
                      ),
                    ),
                  ),
                ),
                ClipRRect(
                  borderRadius: BorderRadius.circular(16.0),
                  child: Material(
                    color: Theme.of(context).scaffoldBackgroundColor,
                    child: InkWell(
                      onTap: () {},
                      child: Padding(
                        padding: const EdgeInsets.all(10.0),
                        child: Text("Logout"),
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
