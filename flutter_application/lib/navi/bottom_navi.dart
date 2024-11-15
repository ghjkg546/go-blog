import 'package:flutter/material.dart';
import 'package:flutter_application_2/navi/navi_controller.dart';
import 'package:get/state_manager.dart';

Widget bottomNavi(NaviController naviController)  {

return Obx(()=>  BottomNavigationBar(
        items: const [
          BottomNavigationBarItem(icon: Icon(Icons.home), label: "首页"),
          // BottomNavigationBarItem(icon: Icon(Icons.share), label: "推荐"),
          BottomNavigationBarItem(
              icon: Icon(Icons.account_box), label: "我的"),
        ],
        currentIndex: naviController.selectedIndex.value,
        onTap: naviController.onBottomNavItemTapped,
      ));

}