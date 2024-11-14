

import 'package:bruno/bruno.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_application_2/apis/app.dart';
import 'package:flutter_application_2/entity/data.dart';
import 'package:flutter_application_2/entity/resource_item.dart';
import 'package:flutter_application_2/pages/login.dart';

import 'package:flutter_html/flutter_html.dart';
import 'package:get/get.dart';



class DetailPage extends StatefulWidget {
  // final Item item; // Define the item variable
  final int itemId;
  const DetailPage({Key? key,  required this.itemId}) : super(key: key);

  @override
  State<DetailPage> createState() => _MyAppState();
}

// class _MyAppState extends State<ListWidget> with SingleTickerProviderStateMixin {
class _MyAppState extends State<DetailPage> {
 late Info item;
  String url = "";
  Map<int, String> diskMap = {
    1: '百度',
    2: '夸克',
    3: '阿里',
    4: '移动彩云',
  };
  bool isFavorite = false; // 本地收藏状态

  void getDetailInfo() async{
    var result ;
    // if(widget.itemId>0){
  result = await userApi.getDetail(widget.itemId);
  //   }else{
  // result = await userApi.getDetail(widget.item.id);
  //   }
   
    final res = ApiResponse.fromJson(result);
    setState(() {
      item = res.data.info;
      isFavorite = res.data.info.isFavorite;
    });
    
  }

 @override
  void initState() {
    super.initState();
    // Initialize isFavorite based on widget.item
    getDetailInfo();
    
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
               Navigator.pop(context, 'refresh');
            },
          ),
        ],
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
             Row(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: [
            // 收藏
            InkWell(
              onTap: () async{
               
                var res = await userApi.fav(item.id);
                if(res['code'] !=0){
                  BrnToast.show(res['msg'], context);
                  Get.to(LoginPage());
                  return;
                }
                setState(() {
                  isFavorite = res['data'];
                });
              },
              child:  Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Icon( isFavorite? Icons.favorite:Icons.favorite_border, color: Colors.red),
                  SizedBox(height: 4),
                  Text('收藏', style: TextStyle(fontSize: 14)),
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
            // item.coverImg.isNotEmpty
            //     ? Image.network(item.coverImg)
            //     : Image.asset('images/noimage.png'),
            const SizedBox(height: 16),

            Text(
              item.name,
              style: const TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 16),
        Html(
            data: item.description,
            
          ),
             
            
            const SizedBox(height: 16),
            Column(
              children: 
            item.diskItemsArray.map((item) {
              return Column(
                children: [
                  SizedBox(height: 10,),
                  ElevatedButton(
                  onPressed: () {
                    url = item.url;
                    Clipboard.setData(ClipboardData(text: url));
                    ScaffoldMessenger.of(context).showSnackBar(
                      SnackBar(content: Text("复制${diskMap[item.type]}链接成功!")),
                    );
                  },
                  style: ElevatedButton.styleFrom(
                    // fixedSize: Size(246, 80), // 设置宽度和高度
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(5), // 设置圆角半径
                    ),
                    backgroundColor: Colors.blue, // 按钮背景色
                  ),
                  child:   Text(
                  
                    '复制${diskMap[item.type]}链接',
                    style: TextStyle(color: Colors.white, fontSize: 24), // 文字颜色
                  ),
                              ),
                ],
              );
            }).toList(),
           


            )
            
            // Add more details here
          ],
        ),
      ),
    );
  }
}
