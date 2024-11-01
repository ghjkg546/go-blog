import 'dart:typed_data';

import 'package:flutter/material.dart';
import 'package:flutter_application_2/entity/data.dart';
import 'package:flutter_application_2/util.dart';
import 'package:screenshot/screenshot.dart';

import 'package:flutter/rendering.dart';
import 'package:image_gallery_saver/image_gallery_saver.dart';

class SecondRoute extends StatelessWidget {
  final Item item;
  SecondRoute({super.key, required this.item});

  final screenshotController = ScreenshotController();
  // final GlobalKey _globalKey = GlobalKey();
  
 _saveLocalImage(uint8List) async {
    try {
          final result = await ImageGallerySaver.saveImage(uint8List);
           showToast(
                message: 'File saved to gallery: $result',
                backgroundColor: Colors.blue,
                textColor: Colors.white,
              );
        // } else {
        //   print('No valid context for capturing screenshot.');
        // }
      } catch (e) {
       showToast(
                message: e.toString(),
                backgroundColor: Colors.blue,
                textColor: Colors.white,
              );
      }
  }



  @override
  Widget build(BuildContext context) {
    return Scaffold(
     appBar: AppBar(
        title: Text("资源详情"),
        actions: [
          IconButton(
            icon: Icon(Icons.favorite_border),
            onPressed: () {
              // 收藏操作
            },
          ),
        ],
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
         
          
            item.coverImg.isNotEmpty
                ? Image.network(item.coverImg)
                : Image.asset('images/noimage.png'),
            const SizedBox(height: 16),
            Text(
              item.name,
              style: const TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 16),
            Text(item.description),
            // Add more details here
          ],
        ),
      ),
    );
  }
}
