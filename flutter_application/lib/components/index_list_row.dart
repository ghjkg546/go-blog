import 'package:bruno/bruno.dart';
import 'package:flutter/material.dart';
import 'package:flutter_application_2/components/simple_item_card.dart';
import 'package:flutter_application_2/entity/index_item.dart';
import 'package:flutter_application_2/entity/resource_item.dart';
import 'package:flutter_application_2/navi/navi_controller.dart';
import 'package:flutter_application_2/pages/detail.dart';
import 'package:flutter_application_2/entity/data.dart';
import 'package:flutter_application_2/pages/list.dart';
import 'package:flutter_html/flutter_html.dart';
import 'package:get/get.dart';

class IndexListRow extends StatelessWidget {
  final List<ResourceItem> resource_items;
  Map<int, String> diskMap = {
    1: '百度',
    2: '夸克',
    3: '阿里',
    4: '移动彩云',
  };
  IndexListRow({super.key, required this.resource_items});

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [

        Padding(
          padding: const EdgeInsets.all(8.0),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                "每日推荐",
                style: const TextStyle(
                  fontWeight: FontWeight.bold,
                  fontSize: 14,
                ),
              ),
              GestureDetector(
                onTap: () {
                  NaviController naviController = Get.find<NaviController>();
                  // Use naviController for your logic
                  naviController.selectedIndex.value = 1;
                  
                },
                child: const Text(
                  "更多",
                  style: TextStyle(
                    fontWeight: FontWeight.bold,
                    fontSize: 12,
                  ),
                ),
              )
            ],
          ),
        ),
        const SizedBox(height: 20),
        Expanded(
          // height: 50,
          child: GridView.builder(
            padding: const EdgeInsets.all(8.0),
            gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
              crossAxisCount: 4, // 每行显示 4 个项
              childAspectRatio: 1, // 设置子项的宽高比
              crossAxisSpacing: 8.0, // 设置子项之间的横向间距
              mainAxisSpacing: 8.0, // 设置子项之间的纵向间距
            ),
            itemCount: resource_items.length, // 数据项的数量
            itemBuilder: (context, index) {

              return SimpleItemCard(item: resource_items[index]); // 渲染每个子项
            },
          ),
        ),
      ],
    );
  }
}
