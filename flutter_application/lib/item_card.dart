import 'package:flutter/material.dart';
import 'package:flutter_application_2/pages/detail.dart';
import 'package:flutter_application_2/entity/data.dart';

class ItemCard extends StatelessWidget {
  final Item item; // 假设 Item 是你的数据模型

  const ItemCard({Key? key, required this.item}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Card(
      elevation: 4.0,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          GestureDetector(
            onTap: () {
              Navigator.push(
                context,
                MaterialPageRoute(
                  builder: (context) => SecondRoute(item: item), // 替换为你的目标页面
                ),
              );
            },
            child: Padding(
            padding: const EdgeInsets.all(8.0),
            child: Text(
              item.name,
              style: const TextStyle(
                fontWeight: FontWeight.bold,
                fontSize: 16,
              ),
            ),
          ),
            // child: SizedBox(
            //   height: 150, // Specify a height for the image
            //   width: double.infinity,
            //   child: item.coverImg.isNotEmpty
            //       ? Image.network(
            //           item.coverImg,
            //           fit: BoxFit.cover,
            //           loadingBuilder: (context, child, loadingProgress) {
            //             if (loadingProgress == null) {
            //               return child;
            //             } else {
            //               return const Center(
            //                 child: CircularProgressIndicator(),
            //               );
            //             }
            //           },
            //           errorBuilder: (context, error, stackTrace) {
            //             return Image.asset('images/noimage.png');
            //           },
            //         )
            //       : Image.asset('images/noimage.png', fit: BoxFit.cover),
            // ),
          ),
          
        ],
      ),
    );
  }
}
