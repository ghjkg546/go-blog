import 'package:flutter/material.dart';
import 'package:flutter_application_2/pages/detail.dart';
import 'package:flutter_application_2/entity/data.dart';

class ItemCard extends StatelessWidget {
  final Item item; // 假设 Item 是你的数据模型

  const ItemCard({Key? key, required this.item}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (context) => SecondRoute(item: item), // 替换为你的目标页面
                    ),
                  );
                },
      child: Card(
        elevation: 4.0,

          child: Column(
            crossAxisAlignment: CrossAxisAlignment.center,
            mainAxisAlignment:MainAxisAlignment.center ,
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
              
            ],
          ),
        
      ),
    );
  }
}
