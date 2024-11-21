import 'package:bruno/bruno.dart';
import 'package:flutter/material.dart';
import 'package:flutter_application_2/pages/detail.dart';
import 'package:flutter_application_2/entity/data.dart';
import 'package:flutter_html/flutter_html.dart';

class ItemCard extends StatelessWidget {
  final Item item; // 假设 Item 是你的数据模型

  Map<int, String> diskMap = {
    1: '百度',
    2: '夸克',
    3: '阿里',
    4: '移动彩云',
  };
  ItemCard({super.key, required this.item});

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
                  style: const TextStyle(
                    fontWeight: FontWeight.bold,
                    fontSize: 16,
                  ),
                ),
              ),
            ),

            // ),
            const SizedBox(
              height: 15,
            ),
            Padding(
                padding: const EdgeInsets.all(8.0),
                child: Wrap(
                  spacing: 10.0, // 元素之间的水平间距
                  runSpacing: 10.0, // 元素之间的垂直间距
                  children: item.diskItemsArray.map((i) {
                    return BrnStateTag(
                      tagText: diskMap[i.type].toString(),
                      tagState: TagState.succeed,
                    );
                  }).toList(),
                )),
                Row(
                  mainAxisAlignment: MainAxisAlignment.end,
                  children: [
                    const Icon(Icons.visibility),
                    const SizedBox(width: 5,),
                    Text(item.views.toString(),style: const TextStyle(fontSize: 12),),
                    const SizedBox(width: 10,),
                  ],
                )
          ],
        ),
      ),
    );
  }
}
