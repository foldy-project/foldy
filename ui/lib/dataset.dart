import 'package:flutter/material.dart';
import 'dart:html' as html;

class Dataset {
  String name;
  String namespace;
  String status;

  Dataset({
    @required this.name,
    @required this.namespace,
  });
}

class DatasetListItemProperty extends StatelessWidget {
  final String name;

  const DatasetListItemProperty(this.name);

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

class DatasetListItem extends StatelessWidget {
  final Dataset dataset;

  const DatasetListItem({
    Key key,
    @required this.dataset,
  }) : super(key: key);

  void navigateTo(BuildContext context) {
    Navigator.of(context).pushNamed("/models", arguments: {
      "name": dataset.name,
      "namespace": dataset.namespace,
    });
  }

  @override
  Widget build(BuildContext context) {
    Color indicatorColor;
    switch (dataset.status) {
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
                                DatasetListItemProperty(dataset.name),
                                DatasetListItemProperty(
                                    dataset.namespace),
                              ],
                            ),
                          ),
                          Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: <Widget>[
                              DatasetListItemProperty("proteinnet-casp11"),
                              DatasetListItemProperty("gromacs"),
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

class DatasetList extends StatelessWidget {
  final List<Dataset> items;

  const DatasetList({
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
              items.map((e) => DatasetListItem(dataset: e)).toList(),
        ),
      ),
    );
  }
}

class DatasetListHeader extends StatelessWidget {
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
                    "Datasets",
                    style: Theme.of(context).textTheme.subtitle2,
                  ),
                ),
                Padding(
                  padding: const EdgeInsets.all(8.0),
                  child: Text(
                    "DatasetS",
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
                            "New Dataset",
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

class DatasetsPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Column(
        // Column is also a layout widget. It takes a list of children and
        // arranges them vertically. By default, it sizes itself to fit its
        // children horizontally, and tries to be as tall as its parent.
        //
        // Invoke "debug painting" (press "p" in the console, choose the
        // "Toggle Debug Paint" action from the Flutter Inspector in Android
        // Studio, or the "Toggle Debug Paint" command in Visual Studio Code)
        // to see the wireframe for each widget.
        //
        // Column has various properties to control how it sizes itself and
        // how it positions its children. Here we use mainAxisAlignment to
        // center the children vertically; the main axis here is the vertical
        // axis because Columns are vertical (the cross axis would be
        // horizontal).
        mainAxisAlignment: MainAxisAlignment.start,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: <Widget>[
          DatasetListHeader(),
          Expanded(
            child: DatasetList(items: [
              Dataset(name: "my-dataset-0", namespace: "default"),
              Dataset(name: "my-dataset-1", namespace: "default"),
              Dataset(name: "my-dataset-2", namespace: "default"),
            ]),
          ),
        ],
      ),
    );
  }
}
