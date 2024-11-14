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
  late TabController _tabController;

  List cates = [];
  List resource_items = [];

  
  bool isLogin = false;
  int _page = 1;
  // int _categoryId = 0;
  @override
  void initState() {
    super.initState();

    // futureReasorce = fetchAlbum(_categoryId, _page);
    // getGoods();
  }

  // 获取列表
  getGoods() async {
    var c1 = await userApi.getReasources();

    var cates1 = CategoryRes.fromJson(c1).data.list;

    setState(() {
      cates = cates1;
      _tabController = TabController(length: cates.length, vsync: this);
      if (cates.length <= 0) {
        return;
      }
      getListData(cates[0].id);
      // print("cates")
      // print(fu)
      _tabController.addListener(() {
        if (_tabController.indexIsChanging) {
          // Fetch new data when tab changes
          setState(() {
            getListData(cates[_tabController.index].id);
            // print(object)
            // futureReasorce = fetchAlbum(_categoryId, _page);
          });
        }
      });
    });
  }

  getListData(int categoryId) async {
    var res = await userApi.getListData(categoryId,1,"");
    var cates1 = DataRes.fromJson(res).data.list;
    setState(() {
      resource_items = cates1;
    });
    // print(cates1);
  }

  void _onItemTapped(int index) {

    //  getInfo();
   
   
  }

  getInfo() async {
    String? accessToken = await UserPreferences.getAccessToken();
    if (accessToken != null) {
      print("tk:" + accessToken);
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

    // return MaterialApp(
    //     title: '资源列表',
    //     theme: ThemeData(
    //       colorScheme: ColorScheme.fromSeed(seedColor: Colors.deepPurple),
    //     ),
    //     home: Scaffold(
    //       resizeToAvoidBottomInset: false,
    //       // appBar: AppBar(
    //       //   title: const Text('资源列表'),
    //       //   // bottom: _showTabBar?TabBar(
    //       //   //   controller: _tabController,
    //       //   //   tabs:
    //       //   //       cates.map((category) => Tab(text: category.name)).toList(),
    //       //   // ):null,
    //       // ),
    //       body: IndexedStack(
    //         index: _selectedIndex, // 根据选中的索引显示不同的组件
    //         children: <Widget>[
    //           //  Center(child: RegisterPage()),
    //           Center(child: ListPage()),
    //           // Center(child: Tougao()),
    //           // Center(child: ComponentDetail()),
            
    //           // Center(child: WebViewPage()),
    //           MyPage()
    //           // isLogin? Center(child: MyPage()):Center(child: LoginPage()),

           
    //           //  Center(child: UserPage()),
    //         ],
    //       ),
    //       bottomNavigationBar: BottomNavigationBar(
    //         items: const <BottomNavigationBarItem>[
    //           BottomNavigationBarItem(
    //             icon: Icon(Icons.home),
    //             label: '列表',
    //           ),
    //           // BottomNavigationBarItem(
    //           //   icon: Icon(Icons.search),
    //           //   label: '搜索',
    //           // ),
    //           // BottomNavigationBarItem(
    //           //   icon: Icon(Icons.notifications),
    //           //   label: '通知',
    //           // ),
    //           BottomNavigationBarItem(
    //             icon: Icon(Icons.person),
    //             label: '我的',
    //           ),
    //         ],
    //         currentIndex: _selectedIndex,
    //         selectedItemColor: Colors.blue,
    //         unselectedItemColor: Colors.grey,
    //         onTap: _onItemTapped,
    //       ),
    //       // body: Column(
    //       //   children: [
    //       //     Expanded(
    //       //       child: TabBarView(
    //       //         controller: _tabController,
    //       //         children: resource_items.map((_item) {
    //       //     // 使用 GridView.builder 替换内容
    //       //     return GridView.builder(
    //       //       gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
    //       //         crossAxisCount: 2, // 每行显示 2 个项
    //       //         childAspectRatio: 1, // 设置子项的宽高比
    //       //       ),
    //       //       itemCount: resource_items.length, // 这里你可以根据具体数据调整数量
    //       //       itemBuilder: (context, index) {
    //       //         return ItemCard(item: resource_items[index]); // 传入每个项
    //       //       },
    //       //     );
    //       //   }).toList(),
    //       //     // 每个 Tab 对应的内容
    //       //      //       return GridView.builder(
    //       //         //         padding: const EdgeInsets.all(8.0),
    //       //         //         gridDelegate:
    //       //         //             const SliverGridDelegateWithFixedCrossAxisCount(
    //       //         //           crossAxisCount: 2,
    //       //         //           crossAxisSpacing: 8.0,
    //       //         //           mainAxisSpacing: 8.0,
    //       //         //         ),
    //       //         //         itemCount: resource.data.list.length,
    //       //         //         itemBuilder: (context, index) {
    //       //         //           final item = resource.data.list[index];

    //       //         // }).toList(),
    //       //     // return Center(child: Text('内容 for ${item.name}'));
    //       //   // }).toList(),

    //       //         // _categories.map((category) {
    //       //         //   return FutureBuilder<Resource>(
    //       //         //     future: futureReasorce,
    //       //         //     builder: (context, snapshot) {
    //       //         //       if (snapshot.connectionState ==
    //       //         //           ConnectionState.waiting) {
    //       //         //         return const Center(
    //       //         //             child: CircularProgressIndicator());
    //       //         //       } else if (snapshot.hasError) {
    //       //         //         return Center(
    //       //         //             child: Text('Error: ${snapshot.error}'));
    //       //         //       } else if (!snapshot.hasData ||
    //       //         //           snapshot.data!.data.list.isEmpty) {
    //       //         //         return const Center(
    //       //         //             child: Text('No data available'));
    //       //         //       }

    //       //         //       final resource = snapshot.data!;

    //       //         //       return GridView.builder(
    //       //         //         padding: const EdgeInsets.all(8.0),
    //       //         //         gridDelegate:
    //       //         //             const SliverGridDelegateWithFixedCrossAxisCount(
    //       //         //           crossAxisCount: 2,
    //       //         //           crossAxisSpacing: 8.0,
    //       //         //           mainAxisSpacing: 8.0,
    //       //         //         ),
    //       //         //         itemCount: resource.data.list.length,
    //       //         //         itemBuilder: (context, index) {
    //       //         //           final item = resource.data.list[index];
    //       //         //           return Card(
    //       //         //             elevation: 4.0,
    //       //         //             child: Column(
    //       //         //               crossAxisAlignment: CrossAxisAlignment.start,
    //       //         //               children: [
    //       //         //                 GestureDetector(
    //       //         //                   onTap: () {
    //       //         //                     Navigator.push(
    //       //         //                       context,
    //       //         //                       MaterialPageRoute(
    //       //         //                         builder: (context) =>
    //       //         //                             SecondRoute(item: item),
    //       //         //                       ),
    //       //         //                     );
    //       //         //                   },
    //       //         //                   child: SizedBox(
    //       //         //                     height:
    //       //         //                         150, // Specify a height for the image
    //       //         //                     width: double.infinity,
    //       //         //                     child: item.coverImg.isNotEmpty
    //       //         //                         ? Image.network(
    //       //         //                             item.coverImg,
    //       //         //                             fit: BoxFit.cover,
    //       //         //                             loadingBuilder: (context, child,
    //       //         //                                 loadingProgress) {
    //       //         //                               if (loadingProgress == null) {
    //       //         //                                 return child;
    //       //         //                               } else {
    //       //         //                                 return const Center(
    //       //         //                                     child:
    //       //         //                                         CircularProgressIndicator());
    //       //         //                               }
    //       //         //                             },
    //       //         //                             errorBuilder: (context, error,
    //       //         //                                 stackTrace) {
    //       //         //                               return Image.asset(
    //       //         //                                   'images/noimage.png');
    //       //         //                             },
    //       //         //                           )
    //       //         //                         : Image.asset('images/noimage.png',
    //       //         //                             fit: BoxFit.cover),
    //       //         //                   ),
    //       //         //                 ),
    //       //         //                 Padding(
    //       //         //                   padding: const EdgeInsets.all(8.0),
    //       //         //                   child: Text(
    //       //         //                     item.name,
    //       //         //                     style: const TextStyle(
    //       //         //                       fontWeight: FontWeight.bold,
    //       //         //                       fontSize: 16,
    //       //         //                     ),
    //       //         //                   ),
    //       //         //                 ),
    //       //         //               ],
    //       //         //             ),
    //       //         //           );
    //       //         //         },
    //       //         //       );
    //       //         //     },
    //       //         //   );
    //       //         // }).toList(),
    //       //       ),
    //       //     ),
    //       //     // NumberPaginator(
    //       //     //   numberPages: 10,
    //       //     //   onPageChange: (int index) {
    //       //     //     setState(() {
    //       //     //       futureReasorce = fetchAlbum(_categoryId, index + 1);
    //       //     //     });
    //       //     //   },
    //       //     // ),
    //       //   ],
    //       // )
    //     ));
  }
}
