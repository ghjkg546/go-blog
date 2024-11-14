import 'package:flutter/material.dart';
import 'package:flutter_application_2/apis/app.dart';
import 'package:flutter_application_2/entity/user.dart';
import 'package:flutter_application_2/mycomponent.dart';
import 'package:flutter_application_2/pages/login.dart';
import 'package:flutter_application_2/pages/setting.dart';
import 'package:flutter_application_2/utils/user_preference.dart';
import 'package:get/get.dart';

class MyPage extends StatefulWidget {
  @override
  _UserInfoPageState createState() => _UserInfoPageState();
}

class _UserInfoPageState extends State<MyPage> {
  String? username = "";
  @override
  void initState() {
    super.initState();

  
   
    getToken();
    // getInfo();
  }

 
 getToken() async {
    String? accessToken = await UserPreferences.getAccessToken();
    if (accessToken != null) {
      print("tk:" + accessToken);
    //  Get.toNamed("login");
    getInfo();
    } else {
      print("null");
      Get.toNamed("login");
    }
  }

  getInfo() async {
    var c1 = await userApi.getInfo();
    var res = UserRes.fromJson(c1);
    if (res.code != 0) {
      Get.toNamed("login");
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
                        // const CircleAvatar(
                        //   radius: 25, // Adjust the size of the avatar
                        //   backgroundImage: AssetImage(
                        //       'avatar.jpg'), // Replace with your asset path
                          
                        // ),
                        const SizedBox(
                          width: 10,
                        ),
                        Column(
                          mainAxisAlignment: MainAxisAlignment.start,
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(username ?? '点击登录',
                                style: TextStyle(fontSize: 24)),
                      
                          ],
                        ),
                      ],
                    ),
                    IconButton(
                      icon: Icon(Icons.settings),
                      onPressed: () {
                                                      Navigator.push(
                                context,
                                MaterialPageRoute(
                                    builder: (context) => SettingPage()),
                              );
                      
                      },
                    )
                  ],
                ),
              ),
              const SizedBox(
                width: 55,
              ),
              Row(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: [
            // 收藏
            InkWell(
              onTap: () {
     
                Get.to(MycomponentPage());
              },
              child: const Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Icon(Icons.favorite, color: Colors.red),
                  SizedBox(height: 4),
                  Text('我的收藏', style: TextStyle(fontSize: 14)),
                ],
              ),
            ),

            // 评论
            // InkWell(
            //   onTap: () {
            //     print('评论 clicked');
            //   },
            //   child: const Column(
            //     mainAxisSize: MainAxisSize.min,
            //     children: [
            //       Icon(Icons.comment, color: Colors.blue),
            //       SizedBox(height: 4),
            //       Text('评论', style: TextStyle(fontSize: 14)),
            //     ],
            //   ),
            // ),
          ],
        ),
        SizedBox(height: 200,),
              Container(
                width: 980,
                height: 400,
            
                child: Padding(
                  padding: EdgeInsets.all(13),
                  child: Column(children: [
                  
                    SizedBox(
                      height: 10,
                    ),
                    
                    
                    SizedBox(
                      height: 20,
                    ),
                    // Row(
                    //   children: [
                    //     Expanded(
                    //       child: ElevatedButton(
                    //         onPressed: () {
                    //           print("jum");
                    //           Navigator.push(
                    //             context,
                    //             MaterialPageRoute(
                    //                 builder: (context) => MycomponentPage()),
                    //           );
                    //         },
                    //         child: Text(
                    //           "我的组件",
                    //           style: TextStyle(fontSize: 16),
                    //         ),
                    //       ),
                    //     ),
                    //     SizedBox(width: 10), // 设置按钮之间的间距
                    //     Expanded(
                    //       child: ElevatedButton(
                    //         onPressed: () {
                    //           // 第二个按钮点击事件
                    //         },
                    //         child: Text(
                    //           "我的收藏",
                    //           style: TextStyle(fontSize: 16),
                    //         ),
                    //       ),
                    //     ),
                    //   ],
                    // ),
                    SizedBox(
                      height: 15,
                    ),
                    
                    SizedBox(
                      height: 26,
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
