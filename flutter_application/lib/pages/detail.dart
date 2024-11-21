import 'package:bruno/bruno.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_application_2/apis/app.dart';
import 'package:flutter_application_2/components/comment_section.dart';
import 'package:flutter_application_2/entity/resource_item.dart';
import 'package:flutter_application_2/pages/login.dart';
import 'package:flutter_application_2/pages/write_comment.dart';

import 'package:flutter_html/flutter_html.dart';
import 'package:get/get.dart';

class DetailPage extends StatefulWidget {
  // final Item item; // Define the item variable
  final int itemId;
  const DetailPage({super.key, required this.itemId});

  @override
  State<DetailPage> createState() => _MyAppState();
}

// class _MyAppState extends State<ListWidget> with SingleTickerProviderStateMixin {
class _MyAppState extends State<DetailPage>
    with SingleTickerProviderStateMixin {
  late Info? item;
  final List<String> _categories = ['详情', '评论'];
  late List comments;
  String url = "";
  Map<int, String> diskMap = {
    1: '百度',
    2: '夸克',
    3: '阿里',
    4: '移动彩云',
  };
  bool isFavorite = false; // 本地收藏状态
  late TabController _tabController;
  void getDetailInfo() async {
    var result;
    item = null;
    result = await userApi.getDetail(widget.itemId);

    final res = ApiResponse.fromJson(result);
    setState(() {
      item = res.data.info;
      comments = res.data.comments;
      isFavorite = res.data.info.isFavorite;
    });
  }

  @override
  void initState() {
    super.initState();
    // Initialize isFavorite based on widget.item

    _tabController = TabController(length: _categories.length, vsync: this);
    _tabController.addListener(() {
      setState(() {}); // 监听 Tab 切换并刷新界面
    });
    getDetailInfo();
  }

  @override
  Widget build(BuildContext context) {
    if (item == null) {
      return const Center(
        child: Column(
          mainAxisSize:
              MainAxisSize.min, // Ensure the content is vertically centered
          children: [
            SizedBox(
              width: 50,
              height: 50,
              child: CircularProgressIndicator(
                strokeWidth: 5,
                valueColor: AlwaysStoppedAnimation<Color>(Colors.green),
              ),
            ),
            SizedBox(height: 10), // Spacing between spinner and text
            Text(
              "加载中...",
              style: TextStyle(fontSize: 16, color: Colors.grey),
            ),
          ],
        ),
      ); // Show a loading indicator
    }
    return Scaffold(
      appBar: AppBar(
        title: const Text("资源详情"),
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
              tabs: _categories.map((category) => Tab(text: category)).toList(),
            ),
          ),
          Expanded(
            child: TabBarView(
              controller: _tabController,
              children: [
                SingleChildScrollView(
                  child: Padding(
                    padding: const EdgeInsets.all(16.0),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                          children: [
                            // 收藏
                            InkWell(
                              onTap: () async {
                                var res = await userApi.fav(item!.id);
                                if (res['code'] != 0) {
                                  BrnToast.show(res['msg'], context);
                                  Get.to(LoginPage());
                                  return;
                                }
                                setState(() {
                                  isFavorite = res['data'];
                                });
                              },
                              child: Column(
                                mainAxisSize: MainAxisSize.min,
                                children: [
                                  Icon(
                                      isFavorite
                                          ? Icons.favorite
                                          : Icons.favorite_border,
                                      color: Colors.red),
                                  const SizedBox(height: 4),
                                  const Text('收藏',
                                      style: TextStyle(fontSize: 14)),
                                ],
                              ),
                            ),
                          ],
                        ),
                        // item.coverImg.isNotEmpty
                        //     ? Image.network(item.coverImg)
                        //     : Image.asset('images/noimage.png'),
                        const SizedBox(height: 16),

                        Text(
                          item!.name,
                          style: const TextStyle(
                              fontSize: 24, fontWeight: FontWeight.bold),
                        ),
                        const SizedBox(height: 16),
                        Html(
                          data: item!.description,
                        ),

                        const SizedBox(height: 16),
                        Column(
                          children: item!.diskItemsArray.map((item) {
                            return Column(
                              children: [
                                const SizedBox(
                                  height: 10,
                                ),
                                ElevatedButton(
                                  onPressed: () {
                                    url = item.url;
                                    Clipboard.setData(ClipboardData(text: url));
                                    ScaffoldMessenger.of(context).showSnackBar(
                                      SnackBar(
                                          content: Text(
                                              "复制${diskMap[item.type]}链接成功!")),
                                    );
                                  },
                                  style: ElevatedButton.styleFrom(
                                    // fixedSize: Size(246, 80), // 设置宽度和高度
                                    shape: RoundedRectangleBorder(
                                      borderRadius:
                                          BorderRadius.circular(5), // 设置圆角半径
                                    ),
                                    backgroundColor: Colors.blue, // 按钮背景色
                                  ),
                                  child: Text(
                                    '复制${diskMap[item.type]}链接',
                                    style: const TextStyle(
                                        color: Colors.white,
                                        fontSize: 24), // 文字颜色
                                  ),
                                ),
                              ],
                            );
                          }).toList(),
                        ),
                      ],
                    ),
                  ),
                ),
                SingleChildScrollView(child: CommentSection(comments: comments))
              ],
            ),
          ),
        ],
      ),
      // 右下角的漂浮按钮
      floatingActionButton: Visibility(
        visible: _tabController.index == 1,
        child: FloatingActionButton(
          onPressed: () async {
            // 点击按钮时跳转到新页面
            final result = await Navigator.push(
              context,
              MaterialPageRoute(
                builder: (context) => WriteCommentPage(id: item!.id),
              ),
            );
            if (result == true) {
              // 如果返回时传递的参数是 true，重新获取评论
              getDetailInfo();
            }
          },
          backgroundColor: Colors.blue,
          child: const Icon(Icons.edit), // 使用写作图标
        ),
      ),

      // ),
    );
  }
}
