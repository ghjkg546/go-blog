import 'package:flutter/material.dart';
import 'package:flutter_application_2/apis/app.dart';
import 'package:flutter_application_2/entity/category.dart';
import 'package:flutter_application_2/entity/data.dart';
import 'package:flutter_application_2/pages/detail.dart';
// import 'package:social_app_ui/util/data.dart';

class RankWidget extends StatefulWidget {
  @override
  _FriendsState createState() => _FriendsState();
}

class _FriendsState extends State<RankWidget> with TickerProviderStateMixin {
  late TabController _tabController;
  List item = [];
  List cates = [];
  int _categoryId = 0;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 0, vsync: this);
    getReasouceItem();
  }

  getListData() async {
    var c1 = await userApi.getRankList(_categoryId);

    var cates1 = DataRes.fromJson(c1).data.list;
    setState(() {
      // _page = validInitialPage;

      // _totalPage = tmpTotalpage;
      item = cates1;
    });
  }

  // 获取列表
  getReasouceItem() async {
    var c1 = await userApi.getCategories();
    var cates1 = CategoryRes.fromJson(c1).data.list;
    setState(() {
      cates = cates1;

      (CategoryRes.fromJson(c1).data.total / 10).ceil();
      _tabController = TabController(length: cates.length, vsync: this);
      if (cates.isEmpty) {
        return;
      }
      _categoryId = cates[0].id;
      getListData();

      _tabController.addListener(() {
        if (_tabController.indexIsChanging) {
          // Fetch new data when tab changes
          setState(() {
            _categoryId = cates[_tabController.index].id;

            getListData();
          });
        }
      });
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        actions: <Widget>[
          IconButton(
            icon: const Icon(
              Icons.filter_list,
            ),
            onPressed: () {},
          ),
        ],
      ),
      body: Column(
        children: [
          Container(
            color: Colors.white, // 设置背景为白色
            child: TabBar(
              controller: _tabController,
              labelColor: Colors.black, // 选中的标签文本颜色设置为黑色
              unselectedLabelColor: Colors.black54, // 未选中的标签文本颜色设置为浅黑色
              indicatorColor: Colors.blue, // 指示器颜色
              tabs: cates.map((category) => Tab(text: category.name)).toList(),
            ),
          ),
          Expanded(
            child: ListView.separated(
              padding: EdgeInsets.all(10),
              separatorBuilder: (BuildContext context, int index) {
                return Align(
                  alignment: Alignment.centerRight,
                  child: Container(
                    height: 0.5,
                    width: MediaQuery.of(context).size.width / 1.3,
                    child: Divider(),
                  ),
                );
              },
              itemCount: item.length,
              itemBuilder: (BuildContext context, int index) {
                Item friend = item[index];
                return Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 8.0),
                  child: ListTile(
                    leading: Text(
                      (index + 1).toString(),
                      style: TextStyle(
                          color: const Color.fromARGB(255, 0, 0, 0),
                          fontSize: 20),
                    ),
                    contentPadding: EdgeInsets.all(0),
                    title: Text(friend.name),
                    trailing: Row(
                      mainAxisSize: MainAxisSize.min, // Add this line
                      mainAxisAlignment: MainAxisAlignment.end,
                      children: [
                        const Icon(Icons.visibility),
                        const SizedBox(
                          width: 5,
                        ),
                        Text(
                          friend.views.toString(),
                          style: const TextStyle(fontSize: 12),
                        ),
                      ],
                    ),
                    onTap: () {
                      Navigator.push(
                        context,
                        MaterialPageRoute(
                          builder: (context) =>
                              DetailPage(itemId: item[index].id), // 替换为你的目标页面
                        ),
                      );
                    },
                  ),
                );
              },
            ),
          ),
        ],
      ),
    );
  }
}
