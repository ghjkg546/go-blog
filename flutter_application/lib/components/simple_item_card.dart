import 'package:bruno/bruno.dart';
import 'package:flutter/material.dart';
import 'package:flutter_application_2/entity/index_item.dart';
import 'package:flutter_application_2/entity/resource_item.dart';
import 'package:flutter_application_2/pages/detail.dart';
import 'package:flutter_application_2/entity/data.dart';
import 'package:flutter_html/flutter_html.dart';

class SimpleItemCard extends StatelessWidget {
  final ResourceItem item; // 假设 Item 是你的数据模型

  Map<int, String> diskMap = {
    1: '百度',
    2: '夸克',
    3: '阿里',
    4: '移动彩云',
  };
  SimpleItemCard({super.key, required this.item});

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: () {
        Navigator.push(
          context,
          MaterialPageRoute(
            builder: (context) => DetailPage(itemId: item.id), // 替换为你的目标页面
          ),
        );
      },
      child: Container(
        decoration: BoxDecoration(
          color: Colors.white,
          border: Border.all(
            color: Colors.grey.shade300, // 浅灰色边框
            width: 1.0,
          ),
          borderRadius: BorderRadius.circular(8), // 可选的圆角边框
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.center,
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            // GestureDetector(

            Padding(
              padding: const EdgeInsets.all(8.0),
              child: Center(
                child: Text(
                  item.name,
                  maxLines: 2, // 限制为单行
  overflow: TextOverflow.ellipsis, // 超出部分显示省略号
                  style: const TextStyle(
                    fontWeight: FontWeight.bold,
                    fontSize: 14,
                  ),
                ),
              ),
            ),

            // ),
            const SizedBox(
              height: 15,
            ),
           
               
          ],
        ),
      ),
    );
  }
}
