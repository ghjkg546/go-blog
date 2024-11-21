import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter_application_2/apis/app.dart';
import 'package:flutter_application_2/entity/user.dart';
import 'package:flutter_application_2/components/mycomponent.dart';
import 'package:flutter_application_2/pages/login.dart';
import 'package:flutter_application_2/pages/setting.dart';
import 'package:flutter_application_2/pages/signin.dart';
import 'package:flutter_application_2/utils/user_preference.dart';
import 'package:get/get.dart';
import 'package:image_picker/image_picker.dart';

class MyPage extends StatefulWidget {
  const MyPage({super.key});

  @override
  _UserInfoPageState createState() => _UserInfoPageState();
}

class _UserInfoPageState extends State<MyPage> {
  String? username = "";
  int score = 0;
  String avatar = "";
  @override
  void initState() {
    super.initState();
    getToken();
  }

  final ImagePicker _picker = ImagePicker();
  File? _selectedImage;
  PickedFile? _pickedFile;

  getToken() async {
    String? accessToken = await UserPreferences.getAccessToken();
    if (accessToken != null) {
      getInfo();
    } else {
      final result = await Navigator.push(
        context,
        MaterialPageRoute(builder: (context) => LoginPage()),
      );

      // 检查返回结果，并更新导航栏
      if (result == 'refresh') {}
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
        score = res.data?.score ?? 0;
        avatar = res.data?.avatar ?? "";
      });
    }
  }

  // 从手机选择图片
  Future<void> pickImage() async {
    print("Attempting to pick an image...");

    final XFile? pickedFile;
    try {
      pickedFile = await _picker.pickImage(source: ImageSource.gallery);
      print("Image picker result received");
    } catch (e) {
      print("Error picking image: $e");
      return;
    }

    if (pickedFile != null) {
      setState(() {
        _selectedImage = File(pickedFile!.path);
      });

      // 选择图片后上传
      await userApi.uploadImage(_selectedImage!);
      getInfo();
      print("Image upload completed");
    } else {
      print("No image selected");
    }
  }

  @override
  Widget build(BuildContext context) {
    return DefaultTabController(
      length: 3, // 标签页数量
      child: Scaffold(
        body: SingleChildScrollView(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.start,
            children: <Widget>[
              Padding(
                padding: const EdgeInsets.all(12),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Row(
                      children: [
                        GestureDetector(
                          onTap: pickImage,
                          child: CircleAvatar(
                            radius: 25, // Adjust the size of the avatar
                            backgroundImage: NetworkImage(
                              avatar, // Replace with your image URL
                            ),
                          ),
                        ),
                        const SizedBox(
                          width: 10,
                        ),
                        Column(
                          mainAxisAlignment: MainAxisAlignment.start,
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(username ?? '点击登录',
                                style: const TextStyle(fontSize: 24)),
                          ],
                        ),
                      ],
                    ),
                    IconButton(
                      icon: const Icon(Icons.settings),
                      onPressed: () {
                        Get.to(SettingPage());
                      },
                    )
                  ],
                ),
              ),
              const SizedBox(
                width: 100,
              ),
              Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  children: [
                    const Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Icon(Icons.local_activity),
                        Text("积分"),
                      ],
                    ),
                    Text(score.toString())
                  ],
                ),
              ),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                children: [
                  // 收藏
                  InkWell(
                    onTap: () {
                      Get.to(const MycomponentPage());
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
                  InkWell(
                    onTap: () {
                      Get.to(const SignInPage());
                    },
                    child: const Column(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Icon(Icons.assignment, color: Colors.blue),
                        SizedBox(height: 4),
                        Text('签到', style: TextStyle(fontSize: 14)),
                      ],
                    ),
                  ),
                ],
              ),
              const SizedBox(
                height: 200,
              ),
              SizedBox(
                width: 980,
                height: 400,
                child: const Padding(
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
  const ItemList({super.key});

  @override
  Widget build(BuildContext context) {
    return ListView.builder(
      itemCount: 10, // 示例数量
      itemBuilder: (context, index) {
        return Row(
          mainAxisAlignment: MainAxisAlignment.spaceAround,
          children: [
            Expanded(child: ItemCard(index: index)),
            const SizedBox(width: 10), // 项目间的间距
            Expanded(child: ItemCard(index: index + 1)),
          ],
        );
      },
    );
  }
}

class ItemCard extends StatelessWidget {
  final int index;

  const ItemCard({super.key, required this.index});

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.all(8.0),
      padding: const EdgeInsets.all(10.0),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(10),
        boxShadow: [
          BoxShadow(
            color: Colors.grey.withOpacity(0.5),
            blurRadius: 5,
            spreadRadius: 2,
            offset: const Offset(0, 3),
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
          const SizedBox(height: 8),
          Text(
            '标题 $index',
            style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 4),
          Text('描述内容 $index'),
          const SizedBox(height: 4),
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
    super.key,
    required this.title,
    required this.icon,
  });

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
          style: const TextStyle(fontSize: 28),
        ),
      ),
    );
  }
}

class CustomList extends StatelessWidget {
  final int count;

  const CustomList({super.key, required this.count});

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
              topLeft: index == 0 ? const Radius.circular(20) : Radius.zero,
              topRight: index == 0 ? const Radius.circular(20) : Radius.zero,
              bottomLeft:
                  index == count - 1 ? const Radius.circular(20) : Radius.zero,
              bottomRight:
                  index == count - 1 ? const Radius.circular(20) : Radius.zero,
            ),
            boxShadow: [
              BoxShadow(
                color: Colors.grey.withOpacity(0.5),
                blurRadius: 5,
                spreadRadius: 2,
                offset: const Offset(0, 3),
              ),
            ],
          ),
          child: const CustomListTile(
            icon: Icons.add_a_photo,
            title: "组件吐槽",
          ),
        );
      },
    );
  }
}
