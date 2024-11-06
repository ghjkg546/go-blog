import 'package:flutter/material.dart';
import 'package:flutter_application_2/apis/app.dart';
import 'package:flutter_application_2/entity/register.dart';
import 'package:flutter_application_2/entity/user.dart';
import 'package:flutter_application_2/mycomponent.dart';
import 'package:flutter_application_2/pages/login.dart';
import 'package:flutter_application_2/utils/user_preference.dart';

class MyPage extends StatefulWidget {
  @override
  _UserInfoPageState createState() => _UserInfoPageState();
}

class _UserInfoPageState extends State<MyPage> {
  String? username = "";
  @override
  void initState() {
    super.initState();
    // getInfo();
    // futureReasorce = fetchAlbum(_categoryId, _page);
    getGoods();
  }

  void _doLogout() {
    UserPreferences.clearUserInfo(); // Assuming you have a function like this

    // Navigate to the login page and remove all previous routes
    Navigator.pushAndRemoveUntil(
      context,
      MaterialPageRoute(builder: (context) => LoginPage()),
      (route) => false,
    );
  }

  getGoods() async {
    var c1 = await userApi.getInfo();
    var res = UserRes.fromJson(c1);
    if (res.code != 0) {
      Navigator.push(
        context,
        MaterialPageRoute(
          builder: (context) => LoginPage(),
        ),
      );
    } else {
      setState(() {
        username = res.data?.userName.toString() ?? "";
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return DefaultTabController(
      length: 3, // 标签页数量
      child: Scaffold(
        body: SingleChildScrollView(
          // Center is a layout widget. It takes a single child and positions it
          // in the middle of the parent.
          child: Column(
            mainAxisAlignment: MainAxisAlignment.start,
            children: <Widget>[
              Padding(
                padding: EdgeInsets.all(12),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Row(
                      children: [
                        const CircleAvatar(
                          radius: 25, // Adjust the size of the avatar
                          backgroundImage: AssetImage(
                              'avatar.jpg'), // Replace with your asset path
                          // If you want to use a placeholder:
                          // child: Text('A'), // Use initials or a placeholder text
                        ),
                        const SizedBox(
                          width: 10,
                        ),
                        Column(
                          mainAxisAlignment: MainAxisAlignment.start,
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(username ?? '点击登录',
                                style: TextStyle(fontSize: 24)),
                            const Text(
                              "登录体验更多功能",
                              style:
                                  TextStyle(fontSize: 16, color: Colors.grey),
                            ),
                          ],
                        ),
                      ],
                    ),
                    IconButton(
                      icon: Icon(Icons.settings),
                      onPressed: () {
                        _doLogout();
                      },
                    )
                  ],
                ),
              ),
              SizedBox(
                width: 55,
              ),
              Container(
                width: 980,
                height: 400,
                decoration: BoxDecoration(
                  color: const Color(0xFFFED88D).withOpacity(0.6), // 设置颜色
                  borderRadius: BorderRadius.circular(30), // 设置圆角
                ),
                child: Padding(
                  padding: EdgeInsets.all(13),
                  child: Column(children: [
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        const Row(
                          children: [
                            SizedBox(
                              width: 120.0, // 设置宽度
                              height: 80.0, // 设置高度
                              child: Icon(
                                Icons.fmd_good, // 替换为你需要的图标
                                size: 40.0, // 图标的大小
                                color: Colors.black, // 图标颜色
                              ),
                            ),
                            SizedBox(
                              width: 10,
                            ),
                          ],
                        ),
                        ElevatedButton(
                          onPressed: () {
                            // 按钮点击事件
                            print("xxxxxxxxxx");
                            Navigator.push(
                              context,
                              MaterialPageRoute(
                                  builder: (context) => MycomponentPage()),
                            );
                          },
                          style: ElevatedButton.styleFrom(
                            fixedSize: Size(246, 80), // 设置宽度和高度
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(10), // 设置圆角半径
                            ),
                            backgroundColor: Colors.blue, // 按钮背景色
                          ),
                          child: const Text(
                            "点击我",
                            style: TextStyle(
                                color: Colors.white, fontSize: 30), // 文字颜色
                          ),
                        )
                      ],
                    ),
                    SizedBox(
                      height: 10,
                    ),
                    SizedBox(
                      height: 10,
                    ),
                    const Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Row(
                          children: [
                            SizedBox(
                              height: 22,
                              width: 18,
                              child: Icon(
                                Icons.thumb_down,
                                size: 28,
                              ),
                            ),
                            SizedBox(
                              width: 15,
                            ),
                          ],
                        ),
                      ],
                    ),
                    SizedBox(
                      height: 20,
                    ),
                    Row(
                      children: [
                        Expanded(
                          child: ElevatedButton(
                            onPressed: () {
                              print("jum");
                              Navigator.push(
                                context,
                                MaterialPageRoute(
                                    builder: (context) => MycomponentPage()),
                              );
                            },
                            child: Text(
                              "我的组件",
                              style: TextStyle(fontSize: 16),
                            ),
                          ),
                        ),
                        SizedBox(width: 10), // 设置按钮之间的间距
                        Expanded(
                          child: ElevatedButton(
                            onPressed: () {
                              // 第二个按钮点击事件
                            },
                            child: Text(
                              "我的收藏",
                              style: TextStyle(fontSize: 16),
                            ),
                          ),
                        ),
                      ],
                    ),
                    SizedBox(
                      height: 15,
                    ),
                    Row(
                      children: [
                        Expanded(child: CustomList(count: 5)),
                      ],
                    ),
                    SizedBox(
                      height: 26,
                    ),
                    Row(
                      children: [
                        Expanded(
                          child: ListView(
                            shrinkWrap: true,
                            children: [
                              CustomListTile(
                                title: "主页",
                                icon: Icons.home,
                              ),

                              // 可以继续添加更多 ListTile
                            ],
                          ),
                        ),
                      ],
                    ),
                  ]),
                ),
              )
            ],
          ),
        ),
      ),
    );
  }
}

class ItemList extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return ListView.builder(
      itemCount: 10, // 示例数量
      itemBuilder: (context, index) {
        return Row(
          mainAxisAlignment: MainAxisAlignment.spaceAround,
          children: [
            Expanded(child: ItemCard(index: index)),
            SizedBox(width: 10), // 项目间的间距
            Expanded(child: ItemCard(index: index + 1)),
          ],
        );
      },
    );
  }
}

class ItemCard extends StatelessWidget {
  final int index;

  const ItemCard({Key? key, required this.index}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: EdgeInsets.all(8.0),
      padding: EdgeInsets.all(10.0),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(10),
        boxShadow: [
          BoxShadow(
            color: Colors.grey.withOpacity(0.5),
            blurRadius: 5,
            spreadRadius: 2,
            offset: Offset(0, 3),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Image.network(
            'https://via.placeholder.com/100', // 示例图片
            height: 100,
            width: 100,
            fit: BoxFit.cover,
          ),
          SizedBox(height: 8),
          Text(
            '标题 $index',
            style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
          ),
          SizedBox(height: 4),
          Text('描述内容 $index'),
          SizedBox(height: 4),
          Text('作者: 作者 $index'),
        ],
      ),
    );
  }
}

class CustomListTile extends StatelessWidget {
  final String title;
  final IconData icon;

  const CustomListTile({
    Key? key,
    required this.title,
    required this.icon,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      height: 120,
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(20),
      ),
      child: ListTile(
        leading: Icon(icon),
        title: Text(
          title,
          style: TextStyle(fontSize: 28),
        ),
      ),
    );
  }
}

class CustomList extends StatelessWidget {
  final int count;

  const CustomList({Key? key, required this.count}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return ListView.builder(
      shrinkWrap: true,
      itemCount: count,
      itemBuilder: (context, index) {
        return Container(
          height: 120,
          decoration: BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.only(
              topLeft: index == 0 ? Radius.circular(20) : Radius.zero,
              topRight: index == 0 ? Radius.circular(20) : Radius.zero,
              bottomLeft:
                  index == count - 1 ? Radius.circular(20) : Radius.zero,
              bottomRight:
                  index == count - 1 ? Radius.circular(20) : Radius.zero,
            ),
            boxShadow: [
              BoxShadow(
                color: Colors.grey.withOpacity(0.5),
                blurRadius: 5,
                spreadRadius: 2,
                offset: Offset(0, 3),
              ),
            ],
          ),
          child: CustomListTile(
            icon: Icons.add_a_photo,
            title: "组件吐槽",
          ),
        );
      },
    );
  }
}
