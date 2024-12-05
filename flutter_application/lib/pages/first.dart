import 'package:bruno/bruno.dart';
import 'package:flutter/material.dart';
import 'package:flutter_application_2/apis/app.dart';
import 'package:flutter_application_2/components/index_list_row.dart';

import 'package:flutter_application_2/entity/index_item.dart';
import 'package:flutter_application_2/navi/navi_controller.dart';
import 'package:get/get.dart';


// void main() => runApp(const MyApp());

class FirstPageWidget extends StatefulWidget {
  const FirstPageWidget({super.key});

  // @override
  @override
  State<FirstPageWidget> createState() => _MyAppState();
}

class _MyAppState extends State<FirstPageWidget> with TickerProviderStateMixin {


  List cates = [];
  List resource_items = [];
  int _page = 1;
  int _categoryId = 0;


  String _keyword = "";
  @override
  void initState() {
    super.initState();
   
    getListData(0);
    // getReasouceItem();
  }

  // 获取列表
 

  getListData(int categoryId) async {
    var c1 = await userApi.getIndexData( _page, _keyword);

    var cates1 = IndexRes.fromJson(c1).data.list;

   
   
    print("xxxxxxxx");
print(cates1);
    setState(() {
      _page = 1;

     cates1.map((item) {
                  print("xxxx22");
                  print(item);
                 
                 });
      resource_items = cates1;
      
    });
  }

  @override
  void dispose() {
    // _tabController.dispose();
    super.dispose();
  }

  final TextEditingController _controller = TextEditingController();

  void _onSearch() {
    _keyword = _controller.text.trim();

    getListData(_categoryId);
  }

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
        title: '资源列表',
        theme: ThemeData(
          colorScheme: ColorScheme.fromSeed(seedColor: Colors.deepPurple),
        ),
        home: SingleChildScrollView(
          child: Column(
            children: [
               Row(
                children: [
                  // 使用 Expanded 使 TextField 占据剩余空间
                  Expanded(
                    child: TextField(
                      controller: _controller,
                      decoration: const InputDecoration(
                        hintText: '请输入搜索内容',
                        border: OutlineInputBorder(),
                      ),
                      onSubmitted: (value) => _onSearch(),
                    ),
                  ),
                  const SizedBox(width: 10),
          
                  BrnSmallMainButton(
                    title: '搜索',
                    onTap: () {
                       NaviController naviController = Get.find<NaviController>();
        naviController.keyword.value = _controller.text.trim();
                  naviController.selectedIndex.value = 1;
                      // _onSearch();
                    },
                  ),
                ],
              ),
            //           SingleChildScrollView(
            // child: 
            
            Column(
              children: resource_items.map((item) {
                return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            
            SizedBox(
              height: 350, // 限制每个 IndexListRow 的高度
              child: IndexListRow(resource_items: item.data),
            ),
          ],
                );
              }).toList(),
            ),
          // )
          
             
             
            ],
          ),
        )
        );
  }
}
