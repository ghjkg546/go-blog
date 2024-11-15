import 'package:bruno/bruno.dart';
import 'package:flutter/material.dart';
import 'package:flutter_application_2/apis/app.dart';
import 'package:flutter_application_2/entity/user.dart';

import 'package:flutter_application_2/pages/login.dart';
import 'package:flutter_application_2/utils/user_preference.dart';

class SettingPage extends StatefulWidget {
  const SettingPage({super.key});

  @override
  _SettingPageState createState() => _SettingPageState();
}

class _SettingPageState extends State<SettingPage> {
  String? username = "";
  @override
  void initState() {
    super.initState();
  }

  void _doLogout() {
    UserPreferences.clearUserInfo(); // Assuming you have a function like this

    // Navigate to the login page and remove all previous routes
    Navigator.push(
      context,
      MaterialPageRoute(builder: (context) => LoginPage()),
    );
  }



  @override
  Widget build(BuildContext context) {
    
    return DefaultTabController(
      
      length: 3, // 标签页数量
      child: Scaffold(
        appBar: AppBar(
        title: const Text('设置'),
        actions: [
          IconButton(
            icon: const Icon(Icons.favorite_border),
            onPressed: () {
               Navigator.pop(context, 'refresh');
            },
          ),
        ],
      ),
        body: SingleChildScrollView(
          // Center is a layout widget. It takes a single child and positions it
          // in the middle of the parent.
          child: Column(
            mainAxisAlignment: MainAxisAlignment.start,
            children: <Widget>[
              const SizedBox(
                width: 55,
              ),
              
              BrnBigMainButton(
                title: '退出登录',
                onTap: () {
                  _doLogout();
                },
              ),
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
