import 'package:fluro/fluro.dart';
import 'package:flutter/material.dart';
import 'dart:html' as html;

import 'package:get_it/get_it.dart';

class MyNavigatorObserver extends NavigatorObserver {
  void Function(String) setRouteName;
  void Function(RouteSettings) onChange;

  MyNavigatorObserver();

  void didPush(Route<dynamic> route, Route<dynamic> previousRoute) {
    if (onChange != null) {
      onChange(route.settings);
    }
  }

  void didPop(Route<dynamic> route, Route<dynamic> previousRoute) {}
  void didRemove(Route<dynamic> route, Route<dynamic> previousRoute) {}
  void didReplace({Route<dynamic> newRoute, Route<dynamic> oldRoute}) {}
}

class NavigationService {
  final GlobalKey<NavigatorState> navigatorKey =
      new GlobalKey<NavigatorState>(debugLabel: '_navKey');

  Future<dynamic> navigateTo(String routeName) {
    return navigatorKey.currentState.pushNamed(routeName);
  }
}

class FoldyToolbarItem extends StatefulWidget {
  final String name;

  final IconData icon;

  final bool isActive;

  final bool indicator;

  final void Function() onPressed;

  final String tooltip;

  const FoldyToolbarItem({
    Key key,
    @required this.name,
    @required this.icon,
    @required this.onPressed,
    @required this.tooltip,
    this.isActive = false,
    this.indicator = true,
  }) : super(key: key);

  @override
  _FoldyToolbarItemState createState() => _FoldyToolbarItemState();
}

class _FoldyToolbarItemState extends State<FoldyToolbarItem> {
  bool hovering = false;

  @override
  Widget build(BuildContext context) {
    return Tooltip(
      message: widget.tooltip,
      preferBelow: true,
      child: MouseRegion(
        onHover: (event) => setState(() => hovering = true),
        onExit: (event) => setState(() => hovering = false),
        child: Padding(
          padding: const EdgeInsets.only(left: 8.0),
          child: Container(
            child: Row(
              children: [
                AnimatedOpacity(
                  opacity: widget.isActive && widget.indicator ? 1.0 : 0.0,
                  duration: Duration(milliseconds: 200),
                  child: Container(
                    width: 2,
                    height: 48,
                    decoration: BoxDecoration(
                      border: Border.all(
                        color: Theme.of(context).accentColor,
                        width: 1.0,
                        style: BorderStyle.solid,
                      ),
                    ),
                  ),
                ),
                Center(
                  child: IconButton(
                    color: Theme.of(context)
                        .primaryIconTheme
                        .color
                        .withAlpha(hovering || widget.isActive ? 255 : 128),
                    padding: EdgeInsets.symmetric(vertical: 24),
                    visualDensity: VisualDensity.comfortable,
                    icon: Icon(widget.icon),
                    onPressed: widget.onPressed,
                  ),
                )
              ],
            ),
          ),
        ),
      ),
    );
  }
}

class FoldyToolbar extends StatefulWidget {
  final void Function(bool dark) setTheme;

  //final void Function(String value) setActiveTab;
  final MyNavigatorObserver navObs;
  final String activeTab;
  final GlobalKey<NavigatorState> navigatorKey;

  FoldyToolbar({
    Key key,
    @required this.setTheme,
    @required this.activeTab,
    this.navObs,
    this.navigatorKey,
    //@required this.setActiveTab,
  }) : super(key: key);

  @override
  _FoldyToolbarState createState() => _FoldyToolbarState();
}

class _FoldyToolbarState extends State<FoldyToolbar> {
  static const _kFontFam = 'NightIcon';

  static const IconData moon = IconData(0xe800, fontFamily: _kFontFam);
  static const IconData moon_inv = IconData(0xe801, fontFamily: _kFontFam);
  String routeName = "";

  @override
  void initState() {
    super.initState();
    widget.navObs?.onChange = (settings) {
      if (!mounted) {
        return;
      }
      print(settings.name);
      () async {
        await Future.delayed(const Duration(milliseconds: 100));
        setState(() => routeName = settings.name);
      }();
    };
  }

  void navigateTo(BuildContext context, String path) {
    print('navigateTo $path, ${widget.navigatorKey}');
    widget.navigatorKey.currentState.pushNamed(path);
  }

  @override
  Widget build(BuildContext context) {
    final dark = Theme.of(context).brightness == Brightness.dark;
    var url = html.window.location.href;
    print(url);
    final start = '/#/';
    final path = url.substring(url.indexOf(start) + start.length);
    print('path=$path');
    return Material(
      color: Theme.of(context).appBarTheme.color,
      child: Container(
        width: 68,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          mainAxisAlignment: MainAxisAlignment.start,
          children: <Widget>[
            FoldyToolbarItem(
              name: "Experiments",
              icon: Icons.search,
              isActive: path.startsWith("experiments"),
              onPressed: () => navigateTo(context, "/experiments"),
              tooltip: "Experiments",
            ),
            FoldyToolbarItem(
              name: "Models",
              icon: Icons.gradient,
              isActive: path.startsWith("models"),
              onPressed: () => navigateTo(context, "/models"),
              tooltip: "Models",
            ),
            FoldyToolbarItem(
              name: "Datasets",
              icon: Icons.data_usage,
              isActive: path.startsWith("datasets"),
              onPressed: () => navigateTo(context, "/datasets"),
              tooltip: "Datasets",
            ),
            FoldyToolbarItem(
              name: "System",
              icon: Icons.settings,
              isActive: path.startsWith("system"),
              onPressed: () => navigateTo(context, "/system"),
              tooltip: "System",
            ),
            Expanded(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.end,
                children: <Widget>[
                  FoldyToolbarItem(
                    name: "Toggle Dark Mode",
                    indicator: false,
                    icon: dark ? moon_inv : moon,
                    onPressed: () => widget.setTheme(dark ? false : true),
                    tooltip: dark ? "Light Mode" : "Dark Mode",
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
