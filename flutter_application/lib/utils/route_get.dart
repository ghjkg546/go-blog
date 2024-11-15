import 'package:flutter_application_2/pages/index.dart';

import 'package:flutter_application_2/pages/login.dart';
import 'package:flutter_application_2/pages/my.dart';
import 'package:get/get.dart';

class RouteGet {
  /// page name
  static const String index = "/index";
  static const String my = "/my";
  static const String login = "/login";

  ///pages map
  static final List<GetPage> getPages = [
    GetPage(
        name: index, 
        page: () => const IndexPage(), 
    ),
    GetPage(
        name: my, 
        page: () => MyPage(), 
    ),
     GetPage(
        name: login, 
        page: () => LoginPage(), 
    ),
  ];
}
