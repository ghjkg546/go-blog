import 'dart:async';
import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter_application_2/apis/app.dart';
import 'package:flutter_application_2/entity/category.dart';
import 'package:flutter_application_2/entity/data.dart';
import 'package:flutter_application_2/pages/my.dart';

import 'package:flutter_application_2/pages/list.dart';
import 'package:flutter_application_2/pages/login.dart';
import 'package:flutter_application_2/pages/register.dart';
import 'package:flutter_application_2/pages/web_view.dart';
import 'package:flutter_application_2/utils/route_get.dart';
import 'package:flutter_application_2/utils/user_preference.dart';
import 'package:get/get.dart';
import 'package:number_paginator/number_paginator.dart';




void main() {
  // setupWindow();
  runApp(const MyApp());
}



// void main() => runApp(const MyApp());

class MyApp extends StatefulWidget {
  const MyApp({super.key});

  @override
  State<MyApp> createState() => _MyAppState();
}

class _MyAppState extends State<MyApp> with SingleTickerProviderStateMixin {
  // late TabController _tabController;

  List cates = [];
  List resource_items = [];

  
  bool isLogin = false;

  // int _categoryId = 0;
  @override
  void initState() {
    super.initState();

    // futureReasorce = fetchAlbum(_categoryId, _page);
    // getGoods();
  }

  // 获取列表
  // getGoods() async {
  //   var c1 = await userApi.getReasources();

  //   var cates1 = CategoryRes.fromJson(c1).data.list;
  //   print(cates1);
  //   setState(() {
  //     cates = cates1;
  //     _tabController = TabController(length: cates.length, vsync: this);
  //     if (cates.length <= 0) {
  //       return;
  //     }
  //     getListData(cates[0].id);
      
  //     _tabController.addListener(() {
  //       if (_tabController.indexIsChanging) {
  //         // Fetch new data when tab changes
  //         setState(() {
  //           getListData(cates[_tabController.index].id);
  //           // print(object)
  //           // futureReasorce = fetchAlbum(_categoryId, _page);
  //         });
  //       }
  //     });
  //   });
  // }

  // getListData(int categoryId) async {
  //   var res = await userApi.getListData(categoryId,1,"");
  //   var cates1 = DataRes.fromJson(res).data.list;
  //   setState(() {
  //     resource_items = cates1;
  //   });
  //   // print(cates1);
  // }

  void _onItemTapped(int index) {

    //  getInfo();
   
   
  }

  getInfo() async {
    String? accessToken = await UserPreferences.getAccessToken();
    if (accessToken != null) {
      print("tk:$accessToken");
      isLogin = true;
    } else {
      print("null");
      isLogin = false;
    }
  }

  @override
  Widget build(BuildContext context) {
   
     return GetMaterialApp(
      title: 'Flutter Demo',
      initialRoute: RouteGet.index,
        
      getPages: RouteGet.getPages,
      theme: ThemeData(
        primarySwatch: Colors.blue,
      )
    );

    
  }
}
