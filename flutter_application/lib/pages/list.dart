import 'dart:async';
import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter_application_2/apis/app.dart';
import 'package:flutter_application_2/utils/user_preference.dart';
import 'package:flutter_application_2/entity/category.dart';
import 'package:flutter_application_2/entity/data.dart';
import 'package:flutter_application_2/pages/detail.dart';
import 'package:flutter_application_2/item_card.dart';

import 'package:number_paginator/number_paginator.dart';

// void main() => runApp(const MyApp());

class ListPage extends StatefulWidget {
  const ListPage({super.key});

  // @override
  State<ListPage> createState() => _MyAppState();
}

class _MyAppState extends State<ListPage> with SingleTickerProviderStateMixin {
  late TabController _tabController;
  // late Future<Resource> futureReasorce;
  // late Future<CategoryRes> futureCategory;
  // final List<String> _categories = ['pc游戏', '安卓游戏'];
  List cates = [];
  List resource_items = [];
  int _page = 1;
  int _categoryId = 0;
  int _totalPage = 10;
  String _keyword = "";
  @override
  void initState() {
    super.initState();
    // getInfo();
    // futureReasorce = fetchAlbum(_categoryId, _page);
    getReasouceItem();
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
    print(_keyword);
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
    // print(cates1);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  final TextEditingController _controller = TextEditingController();

  void _onSearch() {
    _keyword = _controller.text.trim();
   
    // if (query.isNotEmpty) {
    //   // 执行搜索逻辑
    //   print("搜索内容: $query");
    //   // 这里可以触发你的搜索逻辑，如 API 请求
    // } else {
    //   print("输入框为空，请输入搜索内容");
    // }
    getListData(_categoryId);
  }

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
              bottom: TabBar(
                controller: _tabController,
                tabs:
                    cates.map((category) => Tab(text: category.name)).toList(),
              ),
            ),
            body: Column(
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
            ElevatedButton(
              onPressed: _onSearch,
              child: const Text('搜索'),
              style: ElevatedButton.styleFrom(
                padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 16),
              ),
            ),
          ],
                ),
                const SizedBox(height: 20),
                Expanded(
                  child: TabBarView(
                    controller: _tabController,
                    children: resource_items.map((_item) {
                      // 使用 GridView.builder 替换内容
                      return GridView.builder(
                        gridDelegate:
                            const SliverGridDelegateWithFixedCrossAxisCount(
                          crossAxisCount: 2, // 每行显示 2 个项
                          childAspectRatio: 1, // 设置子项的宽高比
                        ),
                        itemCount: resource_items.length, // 这里你可以根据具体数据调整数量
                        itemBuilder: (context, index) {
                          return ItemCard(item: resource_items[index]); // 传入每个项
                        },
                      );
                    }).toList(),
                    // 每个 Tab 对应的内容
                    //       return GridView.builder(
                    //         padding: const EdgeInsets.all(8.0),
                    //         gridDelegate:
                    //             const SliverGridDelegateWithFixedCrossAxisCount(
                    //           crossAxisCount: 2,
                    //           crossAxisSpacing: 8.0,
                    //           mainAxisSpacing: 8.0,
                    //         ),
                    //         itemCount: resource.data.list.length,
                    //         itemBuilder: (context, index) {
                    //           final item = resource.data.list[index];

                    // }).toList(),
                    // return Center(child: Text('内容 for ${item.name}'));
                    // }).toList(),

                    //       final resource = snapshot.data!;

                    //       return GridView.builder(
                    //         padding: const EdgeInsets.all(8.0),
                    //         gridDelegate:
                    //             const SliverGridDelegateWithFixedCrossAxisCount(
                    //           crossAxisCount: 2,
                    //           crossAxisSpacing: 8.0,
                    //           mainAxisSpacing: 8.0,
                    //         ),
                    //         itemCount: resource.data.list.length,
                    //         itemBuilder: (context, index) {
                    //           final item = resource.data.list[index];
                    //           return Card(
                    //             elevation: 4.0,
                    //             child: Column(
                    //               crossAxisAlignment: CrossAxisAlignment.start,
                    //               children: [
                    //                 GestureDetector(
                    //                   onTap: () {
                    //                     Navigator.push(
                    //                       context,
                    //                       MaterialPageRoute(
                    //                         builder: (context) =>
                    //                             SecondRoute(item: item),
                    //                       ),
                    //                     );
                    //                   },
                    //                   child: SizedBox(
                    //                     height:
                    //                         150, // Specify a height for the image
                    //                     width: double.infinity,
                    //                     child: item.coverImg.isNotEmpty
                    //                         ? Image.network(
                    //                             item.coverImg,
                    //                             fit: BoxFit.cover,
                    //                             loadingBuilder: (context, child,
                    //                                 loadingProgress) {
                    //                               if (loadingProgress == null) {
                    //                                 return child;
                    //                               } else {
                    //                                 return const Center(
                    //                                     child:
                    //                                         CircularProgressIndicator());
                    //                               }
                    //                             },
                    //                             errorBuilder: (context, error,
                    //                                 stackTrace) {
                    //                               return Image.asset(
                    //                                   'images/noimage.png');
                    //                             },
                    //                           )
                    //                         : Image.asset('images/noimage.png',
                    //                             fit: BoxFit.cover),
                    //                   ),
                    //                 ),
                    //                 Padding(
                    //                   padding: const EdgeInsets.all(8.0),
                    //                   child: Text(
                    //                     item.name,
                    //                     style: const TextStyle(
                    //                       fontWeight: FontWeight.bold,
                    //                       fontSize: 16,
                    //                     ),
                    //                   ),
                    //                 ),
                    //               ],
                    //             ),
                    //           );
                    //         },
                    //       );
                    //     },
                    //   );
                    // }).toList(),
                  ),
                ),
                NumberPaginator(
                  initialPage: 0,
                  numberPages: _totalPage,
                  onPageChange: (int index) {
                    _page = index + 1;

                    getListData(_categoryId);
                    // setState(() {
                    //   futureReasorce = fetchAlbum(_categoryId, index + 1);
                    // });
                  },
                ),
              ],
            )));
  }
}
