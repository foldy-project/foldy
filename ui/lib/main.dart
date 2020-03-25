import 'package:foldy_ui/experiment/details.dart';

import 'dataset.dart';
import 'package:fluro/fluro.dart';
import 'package:get_it/get_it.dart';

import 'experiment/list_page.dart';
import 'toolbar.dart';
import 'package:flutter/material.dart';

void main() {
  final router = Router();
  GetIt.I.registerSingleton<Router>(router);
  GetIt.I.registerSingleton<NavigationService>(NavigationService());

  var experimentsHandler =
      Handler(handlerFunc: (BuildContext context, Map<String, dynamic> params) {
    final namespace = params["namespace"];
    final name = params["name"];
    if (namespace != null && name != null) {
      //print('namespace=${namespace[0]} name=${name[0]}');
      return ExperimentsDetailsPage(namespace: namespace[0], name: name[0]);
    }
    return ExperimentsPage();
  });
  var datasetsHandler =
      Handler(handlerFunc: (BuildContext context, Map<String, dynamic> params) {
    return DatasetsPage();
  });
  var modelsHandler =
      Handler(handlerFunc: (BuildContext context, Map<String, dynamic> params) {
    return ExperimentsPage();
  });
  router.define("/models", handler: modelsHandler);
  router.define("/datasets", handler: datasetsHandler);
  router.define("/", handler: experimentsHandler);
  router.define("/experiments", handler: experimentsHandler);
  router.define("/experiments/:namespace/:name", handler: experimentsHandler);

  runApp(MyApp());
}

class MyApp extends StatefulWidget {
  @override
  _MyAppState createState() => _MyAppState();
}

class _MyAppState extends State<MyApp> {
  bool dark = true;

  final GlobalKey<NavigatorState> navigatorKey =
      GlobalKey<NavigatorState>(debugLabel: 'navKey');

  MyNavigatorObserver navObs;
  @override
  void initState() {
    super.initState();
    navObs = MyNavigatorObserver();
  }

  void setTheme(bool dark) => setState(() => this.dark = dark);

  @override
  Widget build(BuildContext context) {
    final textColor = dark
        ? Color.fromARGB(255, 240, 245, 245)
        : Color.fromARGB(255, 54, 60, 74);
    final theme = ThemeData(
      brightness: dark ? Brightness.dark : Brightness.light,
      primarySwatch: Colors.blue,
      primaryTextTheme: Theme.of(context).primaryTextTheme.apply(
            bodyColor: textColor,
            displayColor: textColor,
          ),
      textTheme: Theme.of(context).textTheme.apply(
            bodyColor: textColor,
            displayColor: textColor,
          ),
      appBarTheme: Theme.of(context).appBarTheme.copyWith(
            color: Color.fromARGB(255, 73, 87, 99),
          ),
    );
    return MaterialApp(
      theme: theme,
      onGenerateRoute: (settings) => MaterialPageRoute(
        builder: (context) => Row(
          children: [
            FoldyToolbar(
              activeTab: "",
              setTheme: setTheme,
              navObs: navObs,
              navigatorKey: navigatorKey,
            ),
            Expanded(
              child: Navigator(
                key: navigatorKey,
                initialRoute: '/experiments',
                //observers: [navObs],
                onGenerateRoute: GetIt.I<Router>().generator,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class MyCustomRoute<T> extends MaterialPageRoute<T> {
  MyCustomRoute({WidgetBuilder builder, RouteSettings settings})
      : super(builder: builder, settings: settings);

  @override
  Widget buildTransitions(BuildContext context, Animation<double> animation,
      Animation<double> secondaryAnimation, Widget child) {
    // Fades between routes. (If you don't want any animation,
    // just return child.)
    return new FadeTransition(opacity: animation, child: child);
  }
}
