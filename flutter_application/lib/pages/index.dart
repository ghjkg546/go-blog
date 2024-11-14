import 'dart:async';
import 'dart:convert';

import 'package:bruno/bruno.dart';
import 'package:flutter/material.dart';
import 'package:flutter_application_2/apis/app.dart';
import 'package:flutter_application_2/entity/category.dart';
import 'package:flutter_application_2/entity/data.dart';

import 'package:flutter_application_2/item_card.dart';
import 'package:flutter_application_2/navi/bottom_navi.dart';
import 'package:flutter_application_2/navi/navi_controller.dart';
import 'package:flutter_application_2/pages/list.dart';
import 'package:flutter_application_2/pages/my.dart';
import 'package:get/get.dart';

import 'package:number_paginator/number_paginator.dart';

// void main() => runApp(const MyApp());

class IndexPage extends StatefulWidget {
  const IndexPage({super.key});

  // @override
  State<IndexPage> createState() => _MyAppState();
}

class _MyAppState extends State<IndexPage> with SingleTickerProviderStateMixin {
  late TabController _tabController;
  // late Future<Resource> futureReasorce;
  // late Future<CategoryRes> futureCategory;
  // final List<String> _categories = ['pc游戏', '安卓游戏'];
  List cates = [];
  List resource_items = [];
  int _page = 1;
  int _categoryId = 0;
  int _totalPage = 10;
  // var _selectedIndex=0.obs;
  String _keyword = "";
  @override
  void initState() {
    super.initState();
    // getReasouceItem();
  }

  // 获取列表
  getReasouceItem() async {
    var c1 = await userApi.getReasources();

    var cates1 = CategoryRes.fromJson(c1).data.list;

    setState(() {
      cates = cates1;
      _page = 0;
      _totalPage = (CategoryRes.fromJson(c1).data.total / 10).ceil();
      final validInitialPage = _page < _totalPage ? _page : _totalPage - 1;
      _page = validInitialPage;

      _tabController = TabController(length: cates.length, vsync: this);
      if (cates.length <= 0) {
        _page = validInitialPage;
        _totalPage = 1;
        return;
      }
      getListData(cates[0].id);

      _tabController.addListener(() {
        if (_tabController.indexIsChanging) {
          // Fetch new data when tab changes
          setState(() {
            _categoryId = cates[_tabController.index].id;
            _totalPage = _totalPage > 0 ? _totalPage : 1;
            getListData(cates[_tabController.index].id);
            // print(object)
            // futureReasorce = fetchAlbum(_categoryId, _page);
          });
        }
      });
    });
  }

  getListData(int categoryId) async {

    var c1 = await userApi.getListData(categoryId, _page,_keyword);

    var cates1 = DataRes.fromJson(c1).data.list;
    var tmp_totalPage = (CategoryRes.fromJson(c1).data.total / 10).ceil();
    if (tmp_totalPage < 1) {
      tmp_totalPage = 1;
    }
    final validInitialPage = _page < _totalPage ? _page : _totalPage - 1;
    _page = validInitialPage;

    setState(() {
      _page = validInitialPage;

      _totalPage = tmp_totalPage;
      resource_items = cates1;
    });
  }

  @override
  void dispose() {
   
    super.dispose();
  }

  final TextEditingController _controller = TextEditingController();
final NaviController naviController = Get.put(NaviController());
  void _onSearch() {
    _keyword = _controller.text.trim();
   
   
    getListData(_categoryId);
  }

  void _onItemTapped(int index) {
    print(index);
// setState(() {
//       _selectedIndex = index;
//     });
    // if(index==1){
    //   Get.toNamed("/my");
    // }
    
   
  }

  List<Widget> _widgetList = <Widget>[
      ListWidget(),
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
            appBar: AppBar(
              
            ),
            bottomNavigationBar: bottomNavi(naviController),
            
          
            
            // ),
            body:Obx(() => _widgetList.elementAt(naviController.selectedIndex.value)),));
            
            
//              Column(
//               children: [
                
//                 Row(
//                   children: [
//             // 使用 Expanded 使 TextField 占据剩余空间
//             Expanded(
//               child: TextField(
//                 controller: _controller,
//                 decoration: const InputDecoration(
//                   hintText: '请输入搜索内容',
//                   border: OutlineInputBorder(),
//                 ),
//                 onSubmitted: (value) => _onSearch(),
//               ),
//             ),
//             const SizedBox(width: 10),
            
// BrnSmallMainButton(
//   title: '搜索',
//   onTap: () {
//     _onSearch();
    // BrnToast.show('录需求信息', context);
//   },
// ),
            
//           ],
//                 ),
//                 const SizedBox(height: 20),
//                 Expanded(
//                   child: TabBarView(
//                     controller: _tabController,
//                     children: resource_items.map((_item) {
//                       // 使用 GridView.builder 替换内容
//                       return GridView.builder(
//                         padding:  const EdgeInsets.all(8.0),
//                         gridDelegate:
//                             const SliverGridDelegateWithFixedCrossAxisCount(
//                           crossAxisCount: 2, // 每行显示 2 个项
//                           childAspectRatio: 1, // 设置子项的宽高比
//                           crossAxisSpacing: 8.0,
//                                mainAxisSpacing: 8.0,
//                         ),
//                         itemCount: resource_items.length, // 这里你可以根据具体数据调整数量
//                         itemBuilder: (context, index) {
//                           return ItemCard(item: resource_items[index]); // 传入每个项
//                         },
//                       );
//                     }).toList(),
                    
//                   ),
//                 ),
//                 NumberPaginator(
//                   initialPage: 0,
//                   numberPages: _totalPage,
//                   onPageChange: (int index) {
//                     _page = index + 1;

//                     getListData(_categoryId);
//                     // setState(() {
//                     //   futureReasorce = fetchAlbum(_categoryId, index + 1);
//                     // });
//                   },
//                 ),
//               ],
//             )));
  }
}
