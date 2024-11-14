import 'package:flutter/material.dart';
import 'package:flutter_application_2/apis/app.dart';
import 'package:flutter_application_2/entity/fav.dart';
import 'package:flutter_application_2/pages/detail.dart';
import 'package:number_paginator/number_paginator.dart';



class MycomponentPage extends StatefulWidget {
  const MycomponentPage({super.key});

  // @override
  State<MycomponentPage> createState() => _MyAppState();
}


class _MyAppState  extends State<MycomponentPage> with SingleTickerProviderStateMixin {
  List favs = [];
  int _page=1;
  int _totalPage=1;
  getListData(int categoryId) async {

    var c1 = await userApi.getFavList(categoryId,_page,"");

     var res = FavRes.fromJson(c1).data.list;
     print(res);
    var tmp_totalPage = (FavRes.fromJson(c1).data.total / 10).ceil();
    if (tmp_totalPage < 1) {
      tmp_totalPage = 1;
    }
    final validInitialPage = _page < _totalPage ? _page : _totalPage - 1;
    _page = validInitialPage;

    setState(() {
       _page = validInitialPage;
      favs = res;
       _totalPage = tmp_totalPage;
      
    });
  }

  @override
  void initState() {
    super.initState();
    getListData(1);
  }
  
  @override
  Widget build(BuildContext context) {
    return DefaultTabController(
      length: 6, // 标签页数量
      child: Scaffold(
        appBar: AppBar(
          title: Text("我的收藏", style: TextStyle(fontSize: 36)),
          
        ),
        body: Column(
          children: [
            Column(
               children: (favs.map((_item) {
                          // 使用 GridView.builder 替换内容
                          return ItemCard(item: _item);
                         
                        }).toList()
            )
            ),
            NumberPaginator(
                  initialPage: 0,
                  numberPages: _totalPage,
                  onPageChange: (int index) {
                    _page = index + 1;

                    getListData(1);
                    // setState(() {
                    //   futureReasorce = fetchAlbum(_categoryId, index + 1);
                    // });
                  },
                ),
          ],
        )

      ),
    );
  }
}


class ItemCard extends StatelessWidget {
  // final int index;
  FavItem item;

   ItemCard({Key? key, required this.item}) : super(key: key);

  @override
  Widget build(BuildContext context) {
  return GestureDetector(
      onTap: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute(
                      builder: (context) => DetailPage(itemId:item.id), // 替换为你的目标页面
                    ),
                  );
                },
      
    child: Container(
      height: 100,  // 设置高度，Row 会撑满这个高度
      margin: EdgeInsets.all(4),
      padding: EdgeInsets.all(12),
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
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.stretch, // 让Row的子组件高度撑满
        children: [
          // 左侧的子组件
         
          
          // 右侧的子组件
          Expanded(
            
            child: Container(
          
              child:  Column(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      item.name,
                      style: TextStyle(fontSize: 16),
                    ),
                    
                  ],
                ),
            ),
          ),
    
          //  Expanded(
          //       flex: 1,
          //       child: Column(
          //          mainAxisAlignment: MainAxisAlignment.end,
          //         children:[ IconButton(
          //           icon: Icon(Icons.delete),
          //           onPressed: () {
          //             // Delete action here
          //           },
                   
          //         )],
          //       ),
          //     ),
        ],
      ),
    ),
  );
}


}
