import 'package:bruno/bruno.dart';
import 'package:flutter/material.dart';
import 'package:flutter_application_2/apis/app.dart';
import 'package:flutter_application_2/entity/category.dart';
import 'package:flutter_application_2/entity/data.dart';
import 'package:flutter_application_2/components/item_card.dart';
import 'package:flutter_application_2/navi/navi_controller.dart';
import 'package:get/get.dart';

import 'package:number_paginator/number_paginator.dart';

// void main() => runApp(const MyApp());

class ListWidget extends StatefulWidget {
  const ListWidget({super.key});

  // @override
  @override
  State<ListWidget> createState() => _MyAppState();
}

class _MyAppState extends State<ListWidget> with TickerProviderStateMixin {
  late TabController _tabController;
  final NaviController naviController = Get.find<NaviController>();
  List cates = [];
  List resource_items = [];
  int _page = 1;
  int _categoryId = 0;
  int _totalPage = 10;
  final int _pageSize = 50;

  String _keyword = "";
  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 0, vsync: this);
    print(naviController.keyword.value);
    if (naviController.keyword.value != "") {
      _controller.text = naviController.keyword.value;
      _onSearch();
    }
    getReasouceItem();
  }

  // 获取列表
  getReasouceItem() async {
    var c1 = await userApi.getCategories();

    var cates1 = CategoryRes.fromJson(c1).data.list;

    setState(() {
      cates = cates1;
      _page = 0;
      _totalPage = (CategoryRes.fromJson(c1).data.total / _pageSize).ceil();
      final validInitialPage = _page < _totalPage ? _page : _totalPage - 1;
      _page = validInitialPage;

      _tabController = TabController(length: cates.length, vsync: this);
      if (cates.isEmpty) {
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
          });
        }
      });
    });
  }

  getListData(int categoryId) async {
    var c1 = await userApi.getListData(categoryId, _page, _keyword);

    var cates1 = DataRes.fromJson(c1).data.list;
    var tmpTotalpage = (CategoryRes.fromJson(c1).data.total / _pageSize).ceil();
    if (tmpTotalpage < 1) {
      tmpTotalpage = 1;
    }
    final validInitialPage = _page < _totalPage ? _page : _totalPage - 1;
    _page = validInitialPage;

    setState(() {
      _page = validInitialPage;

      _totalPage = tmpTotalpage;
      resource_items = cates1;
    });
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  final TextEditingController _controller = TextEditingController();

  void _onSearch() {
    _keyword = _controller.text.trim();
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
        home: Column(
          children: [
            Container(
              color: Colors.blue,
              child: TabBar(
                controller: _tabController,
                labelColor: Colors.white,
                unselectedLabelColor: Colors.white70,
                indicatorColor: Colors.yellow,
                tabs:
                    cates.map((category) => Tab(text: category.name)).toList(),
              ),
            ),
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
                    _onSearch();
                  },
                ),
              ],
            ),
            const SizedBox(height: 20),
            resource_items.length == 0
                ? Container(
                    width: double.infinity,
                    height: 300,
                    color: Colors.white, // 设置背景色（可选）
                    child: Center(
                      child: Text(
                        '空空如也',
                        style: TextStyle(
                          fontSize: 24,
                          fontWeight: FontWeight.bold,
                          color: Colors.grey,
                        ),
                      ),
                    ),
                  )
                : Expanded(
                    child: TabBarView(
                      controller: _tabController,
                      children: resource_items.map((item) {
                        // 使用 GridView.builder 替换内容
                        return GridView.builder(
                          padding: const EdgeInsets.all(8.0),
                          gridDelegate:
                              const SliverGridDelegateWithFixedCrossAxisCount(
                            crossAxisCount: 2, // 每行显示 2 个项
                            childAspectRatio: 1, // 设置子项的宽高比
                            crossAxisSpacing: 8.0,
                            mainAxisSpacing: 8.0,
                          ),
                          itemCount: resource_items.length, // 这里你可以根据具体数据调整数量
                          itemBuilder: (context, index) {
                            return ItemCard(
                                item: resource_items[index]); // 传入每个项
                          },
                        );
                      }).toList(),
                    ),
                  ),
            NumberPaginator(
              initialPage: 0,
              numberPages: _totalPage,
              onPageChange: (int index) {
                _page = index + 1;

                getListData(_categoryId);
              },
            ),
          ],
        ));
  }
}
