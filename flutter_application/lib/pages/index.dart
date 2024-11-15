
import 'package:flutter/material.dart';

import 'package:flutter_application_2/navi/bottom_navi.dart';
import 'package:flutter_application_2/navi/navi_controller.dart';
import 'package:flutter_application_2/pages/list.dart';
import 'package:flutter_application_2/pages/my.dart';
import 'package:get/get.dart';

// void main() => runApp(const MyApp());

class IndexPage extends StatefulWidget {
  final int? selectedIndex;
  const IndexPage({super.key, this.selectedIndex});

  // @override
  @override
  State<IndexPage> createState() => _MyAppState();
}

class _MyAppState extends State<IndexPage> with SingleTickerProviderStateMixin {

  List cates = [];
  List resource_items = [];


  @override
  void initState() {
    super.initState();

  }



  final NaviController naviController = Get.put(NaviController());


  final List<Widget> _widgetList = <Widget>[
    const ListWidget(),
    //  ListWidget(),
    MyPage()
    // testWidget(),
    // settingsWidget()
  ];

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
        title: '资源列表',
        theme: ThemeData(
          colorScheme: ColorScheme.fromSeed(seedColor: Colors.deepPurple),
        ),
        home: Scaffold(
          resizeToAvoidBottomInset: false,
          appBar: AppBar(),
          bottomNavigationBar: bottomNavi(naviController),

          // ),
          body: Obx(
              () => _widgetList.elementAt(naviController.selectedIndex.value)),
        ));
  }
}
